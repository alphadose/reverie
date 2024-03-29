package mongo

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/reverie/types"
	"github.com/reverie/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// userCollectionKey is the collection for all users
	userCollectionKey = "users"

	// userEmailKey is the key holding the email of a user
	userEmailKey = types.EmailKey

	// usernameKey is the key holding the username of a user
	usernameKey = types.UsernameKey

	// userPasswordKey is the key holding the password of a user/instance
	userPasswordKey = "password"

	// userRoleKey is the key denoting the role of a user
	userRoleKey = types.RoleKey

	// userInventoryKey is the key denoting the inventory of a user
	userInventoryKey = "inventory"

	// userVerifiedKey is the key denoting the whether the user is verified or not
	userVerifiedKey = "verified"
)

// The link to the user collection
var userCollection = db.Collection(userCollectionKey)

// upsertUser is an abstraction over UpdateOne which updates a user
// or inserts it if the corresponding document doesn't exist
func upsertUser(filter types.M, user *types.User) error {
	return updateOne(userCollection, filter, user, options.FindOneAndUpdate().SetUpsert(true))
}

// UpdateUser is an abstraction over UpdateOne which updates a user
func UpdateUser(filter types.M, data interface{}) error {
	return updateOne(userCollection, filter, data)
}

// InitVendorInventory initializes the vendor's inventory
// Should be called only once per vendor and this call should be authorized by us
func InitVendorInventory(vendorEmail string, inventory *types.Inventory) error {
	filter := types.M{
		userEmailKey: vendorEmail,
	}
	// Make a map for initializing the vendor's inventory
	initMap := make(map[string]int64)

	inventoryValues := reflect.ValueOf(*inventory)
	inventoryKeys := reflect.TypeOf(*inventory)

	for i := 0; i < inventoryValues.NumField(); i++ {
		key := concat(userInventoryKey, inventoryKeys.Field(i).Name)
		value := inventoryValues.Field(i).Int()
		if value < 0 {
			return errors.New("Vendor inventory values cannot be negative")
		}
		// No need to put fields with zero values into mongoDB, wastage of space
		if value == 0 {
			continue
		}
		initMap[key] = value
	}

	return updateOne(userCollection, filter, initMap)
}

// UpdateVendorInventoryOnAcceptance updates a vendor's inventory after their offer has been accepted on a post
// This deducts the offer contents from the vendor's current inventory
func UpdateVendorInventoryOnAcceptance(vendorEmail string, offer types.Inventory) error {
	filter := types.M{
		userEmailKey: vendorEmail,
	}

	// Make a map for decrementing the vendor's inventory
	decrementMap := make(map[string]int64)

	offerValues := reflect.ValueOf(offer)
	offerKeys := reflect.TypeOf(offer)

	for i := 0; i < offerValues.NumField(); i++ {
		key := concat(userInventoryKey, offerKeys.Field(i).Name)
		value := offerValues.Field(i).Int() * -1
		if value == 0 {
			continue
		}
		decrementMap[key] = value
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	return userCollection.FindOneAndUpdate(ctx, filter, types.M{
		"$inc": decrementMap,
	}).Err()
}

// ReleaseSingleVendorInventory updates a vendor's inventory after their offer has been rejected from a post
// The remainder inventory items are added back to the vendor's inventory pool
func ReleaseSingleVendorInventory(vendorEmail string, offer types.Inventory) error {
	filter := types.M{
		userEmailKey: vendorEmail,
	}

	// Make a map for incrementing the vendor's inventory
	incrementMap := make(map[string]int64)

	offerValues := reflect.ValueOf(offer)
	offerKeys := reflect.TypeOf(offer)

	for i := 0; i < offerValues.NumField(); i++ {
		key := concat(userInventoryKey, offerKeys.Field(i).Name)
		value := offerValues.Field(i).Int()
		if value == 0 {
			continue
		}
		incrementMap[key] = value
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	return userCollection.FindOneAndUpdate(ctx, filter, types.M{
		"$inc": incrementMap,
	}).Err()
}

// FetchUsers returns all users given their email ids
func FetchUsers(emailList []string) ([]types.M, error) {
	return fetchDocs(userCollection, types.M{
		userEmailKey: types.M{
			"$in": emailList,
		},
	}, options.Find())
}

// ReleaseVendorInventories releases inventories of all vendors bound to a job after it is marked as COMPLETED by the client
// This set of inventory is then added back to their respective vendor's inventory pool
func ReleaseVendorInventories(acceptedOffers map[string]types.Offer) error {
	updates := make([]mongo.WriteModel, 0)

	for offerKey, offer := range acceptedOffers {
		vendorEmail, err := utils.Decrypt(offerKey)
		// TODO : Mail us in this case, because this means a vendor's inventory did not get released properly
		// Possible cause: Encryption key or nonce was changed in config.toml mid-production
		if err != nil {
			utils.LogError("kekw", err)
			continue
		}

		// Make a map for incrementing the vendor's inventory
		incrementMap := make(map[string]int64)

		offerValues := reflect.ValueOf(offer.Content)
		offerKeys := reflect.TypeOf(offer.Content)

		for i := 0; i < offerValues.NumField(); i++ {
			key := concat(userInventoryKey, offerKeys.Field(i).Name)
			value := offerValues.Field(i).Int()
			if value == 0 {
				continue
			}
			incrementMap[key] = value
		}

		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(types.M{
			userEmailKey: vendorEmail,
		})
		operation.SetUpdate(types.M{
			"$inc": incrementMap,
		})

		updates = append(updates, operation)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	result, err := userCollection.BulkWrite(ctx, updates)

	// TODO : proper logging on job completion
	utils.LogInfo("Released Inventories", "Released Inventories %v", *result)

	return err
}

// UpdatePassword is an abstraction over UpdateOne which updates a user's password
func UpdatePassword(email, newHashedPassword string) error {
	filter := types.M{
		userEmailKey: email,
	}
	updatePayload := types.M{
		userPasswordKey: newHashedPassword,
	}
	return updateOne(userCollection, filter, updatePayload, nil)
}

// RegisterUser is an abstraction over InsertOne which inserts user into the mongoDB
func RegisterUser(user *types.User) (interface{}, error) {
	return insertOne(userCollection, user)
}

// FetchSingleUser returns a user based on a email based filter
func FetchSingleUser(email string, opts ...*options.FindOneOptions) (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	user := &types.User{}
	err := userCollection.FindOne(ctx, types.M{userEmailKey: email}, opts...).Decode(user)
	return user, err
}

// FetchVendorInventory returns the inventory of a vendor
func FetchVendorInventory(vendorEmail string) (*types.Inventory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	user := &types.User{}
	err := userCollection.FindOne(ctx, types.M{
		userEmailKey: vendorEmail,
	}, options.FindOne().SetProjection(types.M{userInventoryKey: 1})).Decode(user)
	if err != nil {
		return nil, err
	}
	return user.Inventory, nil
}

// FetchSingleUserWithoutPassword returns a user based on a email based filter without his/her password
func FetchSingleUserWithoutPassword(email string) (*types.User, error) {
	return FetchSingleUser(
		email,
		&options.FindOneOptions{
			Projection: types.M{userPasswordKey: 0},
		})
}

// IsUniqueEmail checks if an email id is unique or not
func IsUniqueEmail(email string) (bool, error) {
	count, err := countDocs(userCollection, types.M{userEmailKey: email})
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// VerifyUserEmail sets the user's verified field to true
func VerifyUserEmail(userID string) error {
	docID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := types.M{
		primaryKey: docID,
	}
	updatePayload := types.M{
		userVerifiedKey: true,
	}
	return updateOne(userCollection, filter, updatePayload, nil)
}
