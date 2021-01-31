package mongo

import (
	"context"
	"os"
	"time"

	"github.com/reverie/types"
	"github.com/reverie/utils"

	"github.com/reverie/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var client, err = mongo.Connect(ctx, options.Client().ApplyURI(configs.MongoConfig.URL))
var db = client.Database(projectDatabase)

func setupAdmin() {
	adminInfo := configs.AdminConfig
	pwd, err := utils.HashPassword(adminInfo.Password)
	if err != nil {
		utils.LogError("Mongo-Connection-1", err)
		os.Exit(1)
	}
	admin := &types.User{
		Username:      adminInfo.Username,
		Email:         adminInfo.Email,
		Password:      pwd,
		Role:          types.Admin,
		Company:       "EzFlo",
		Designation:   "Co-Founder",
		Phone:         "911",
		OfficeAddress: "Rourkela",
		Verified:      true,
		Inventory: &types.Inventory{
			Truck:      100,
			Crane:      100,
			BoomLifter: 100,
		},
	}
	filter := types.M{userEmailKey: adminInfo.Email}
	if err := upsertUser(filter, admin); err != nil && err != ErrNoDocuments {
		utils.LogError("Mongo-Connection-2", err)
	}
	utils.LogInfo("Mongo-Connection-3", "%s (%s) has been given admin privileges", adminInfo.Username, adminInfo.Email)
}

func createGeoIndex() {
	indexes := []mongo.IndexModel{
		{
			Keys: types.M{
				postLocationKey: "2dsphere",
			},
		},
		{
			Keys: types.M{
				updatedKey: 1,
			},
		},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	if _, err := postCollection.Indexes().CreateMany(ctx, indexes, opts); err != nil {
		utils.LogError("Mongo-Connection-7", err)
	}
}

func setup() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		utils.Log("Mongo-Connection-4", "MongoDB connection was not established", utils.ErrorTAG)
		utils.LogError("Mongo-Connection-5", err)
		time.Sleep(5 * time.Second)
		setup()
	} else {
		utils.LogInfo("Mongo-Connection-6", "MongoDB Connection Established")
		setupAdmin()
		createGeoIndex()
	}
}

func init() {
	go setup()
}
