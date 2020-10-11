package mongo

// db.getCollection('kekw').find({$or: [{ "looks.chu": { $gt: 0 } }, { "looks.meow": { $gt: 0 } }]})

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
		Username: adminInfo.Username,
		Email:    adminInfo.Email,
		Password: pwd,
		Role:     types.Admin,
	}
	filter := types.M{EmailKey: adminInfo.Email}
	if err := UpsertUser(filter, admin); err != nil && err != ErrNoDocuments {
		utils.LogError("Mongo-Connection-2", err)
	}
	utils.LogInfo("Mongo-Connection-3", "%s (%s) has been given admin privileges", adminInfo.Username, adminInfo.Email)
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
	}
}

func init() {
	go setup()
}
