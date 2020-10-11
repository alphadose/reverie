package mongo

import (
	"context"
	"time"

	"github.com/reverie/types"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpsertUser is an abstraction over UpdateOne which updates a user
// or inserts it if the corresponding document doesn't exist
func UpsertUser(filter types.M, user *types.User) error {
	return UpdateOne(UserCollection, filter, user, options.FindOneAndUpdate().SetUpsert(true))
}

// UpdateUser is an abstraction over UpdateOne which updates a user
func UpdateUser(filter types.M, data interface{}) error {
	return UpdateOne(UserCollection, filter, data, nil)
}

// UpdateVendorInventory is an abstraction over UpdateOne which updates the vendor's inventory
func UpdateVendorInventory(email string, inventory *types.Inventory) error {
	filter := types.M{
		EmailKey: email,
	}
	updatePayload := types.M{
		UserInventoryKey: inventory,
	}
	return UpdateOne(UserCollection, filter, updatePayload, nil)
}

// UpdatePassword is an abstraction over UpdateOne which updates a user's password
func UpdatePassword(email, newHashedPassword string) error {
	filter := types.M{
		EmailKey: email,
	}
	updatePayload := types.M{
		PasswordKey: newHashedPassword,
	}
	return UpdateOne(UserCollection, filter, updatePayload, nil)
}

// RegisterUser is an abstraction over InsertOne which inserts user into the mongoDB
func RegisterUser(user *types.User) (interface{}, error) {
	return InsertOne(UserCollection, user)
}

// FetchSingleUser returns a user based on a email based filter
func FetchSingleUser(email string, opts ...*options.FindOneOptions) (*types.User, error) {
	collection := link.Collection(UserCollection)
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	user := &types.User{}
	err := collection.FindOne(ctx, types.M{EmailKey: email}, opts...).Decode(user)
	return user, err
}

// FetchSingleUserWithoutPassword returns a user based on a email based filter without his/her password
func FetchSingleUserWithoutPassword(email string) (*types.User, error) {
	return FetchSingleUser(
		email,
		&options.FindOneOptions{
			Projection: types.M{PasswordKey: 0},
		})
}

// IsUniqueEmail checks if an email id is unique or not
func IsUniqueEmail(email string) (bool, error) {
	count, err := CountDocs(UserCollection, types.M{EmailKey: email})
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
