package database

import (
	"context"
	"curd-service/constants"
	"curd-service/logger"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//var CNX = Connection()

func Connection() *mongo.Client {
	// Set client options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	dbUrl := "mongodb+srv://" + constants.DB_USERNAME + ":" + constants.DB_PASSWORD + "@cluster0.gdbqwc7.mongodb.net/?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(dbUrl).SetServerAPIOptions(serverAPI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		msg := fmt.Sprintf("Error while connect to DB %s", err.Error())
		logger.Log.Printf(msg)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		msg := fmt.Sprintf("Error while connect to DB %s", err.Error())
		logger.Log.Printf(msg)
	}

	fmt.Println("Connected to MongoDB!")
	logger.Log.Printf("Connected to MongoDB!")
	return client
}
func CloseClientDB(client *mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		msg := fmt.Sprintf("Error while connect to DB %s", err.Error())
		logger.Log.Printf(msg)
	}

	// TODO optional you can log your closed MongoDB client
	fmt.Println("Connection to MongoDB closed.")
	logger.Log.Printf("Connection to MongoDB closed.")
}
