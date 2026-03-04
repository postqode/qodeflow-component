package mongodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/metadata"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	_ = activity.Register(&MongoDbActivity{}, New)
}

/*
Integration with MongoDb
inputs: {uri, dbName, collection, method, [keyName, keyValue, value]}
outputs: {output, count}
*/
type MongoDbActivity struct {
	settings *Settings
}

// New creates a new MongoDB activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	return &MongoDbActivity{settings: s}, nil
}

// Metadata returns the activity's metadata
func (a *MongoDbActivity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval implements activity.Activity.Eval - MongoDb integration
func (a *MongoDbActivity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	if a.settings.ConnectionURI == "" {
		return false, fmt.Errorf("connection URI is required")
	}

	clientOptions := options.Client().ApplyURI(a.settings.ConnectionURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		ctx.Logger().Errorf("Connection error: %v", err)
		return false, err
	}
	defer func() {
		if disconnectErr := client.Disconnect(context.Background()); disconnectErr != nil {
			ctx.Logger().Errorf("Error disconnecting from MongoDB: %v", disconnectErr)
		}
	}()

	db := client.Database(a.settings.DbName)
	coll := db.Collection(a.settings.Collection)

	filter := bson.M{input.KeyName: input.KeyValue}

	switch strings.ToUpper(input.Method) {
	case "GET":
		result := coll.FindOne(context.Background(), filter)
		val := make(map[string]interface{})
		err := result.Decode(&val)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				_ = ctx.SetOutputObject(&Output{Output: nil})
				return true, nil
			}
			return false, err
		}

		ctx.Logger().Debugf("Get Results %v", val)
		_ = ctx.SetOutputObject(&Output{Output: val})

	case "DELETE":
		result, err := coll.DeleteMany(context.Background(), filter)
		if err != nil {
			return false, err
		}

		ctx.Logger().Debugf("Delete Results %+v", result)
		_ = ctx.SetOutputObject(&Output{Count: result.DeletedCount})

	case "INSERT":
		result, err := coll.InsertOne(context.Background(), input.Data)
		if err != nil {
			return false, err
		}
		ctx.Logger().Debugf("Insert Results %+v", result)
		_ = ctx.SetOutputObject(&Output{Output: result.InsertedID})

	case "REPLACE":
		result, err := coll.ReplaceOne(context.Background(), filter, input.Data)
		if err != nil {
			return false, err
		}

		ctx.Logger().Debugf("Replace Results %+v", result)
		_ = ctx.SetOutputObject(&Output{
			Output: result.UpsertedID,
			Count:  result.ModifiedCount,
		})

	case "UPDATE":
		update := bson.M{"$set": input.Data}
		result, err := coll.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return false, err
		}

		ctx.Logger().Debugf("Update Results %+v", result)
		_ = ctx.SetOutputObject(&Output{
			Output: result.UpsertedID,
			Count:  result.ModifiedCount,
		})

	default:
		ctx.Logger().Errorf("unsupported method '%s'", input.Method)
		return false, fmt.Errorf("unsupported method '%s'", input.Method)
	}

	return true, nil
}
