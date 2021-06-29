package Dboperation

import (
	"Giftcode-mongo-protobuf/config"
	"Giftcode-mongo-protobuf/model"
	"context"
	"fmt"
	_ "go.mongodb.org/mongo-driver/benchmark"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func Searchmongo(user, id string) (model.UerInfo, error) {
	user_filter := bson.D{{user, id}}
	var data model.UerInfo
	err := config.Collection.FindOne(context.TODO(), user_filter).Decode(&data)
	if data.UID == "" && err != nil {
		return data, err
	}
	return data, nil
}

func Insearchmongo(user model.UerInfo) bool {
	_, err := config.Collection.InsertOne(context.TODO(), user)
	if err != nil {
		return false
	}
	return true
}

func Updatamongo(userid, id, propertyid, property string) bool {
	user_filter := bson.D{{userid, id}}
	updata := bson.D{{"$set", bson.D{{propertyid, property}}}}
	result, err := config.Collection.UpdateOne(context.TODO(), user_filter, updata)
	if err != nil {
		log.Fatal(err)
	}
	if result.MatchedCount != 0 {
		fmt.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
		return false
	} else if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}
	return true
}

func ExistUserId(id string, col string) bool {

	result, _ := Searchmongo("Uid", id)
	if result.UID == "" {
		return false
	}
	return true
}
