package order

import (
	"context"
	"log"
	"time"

	"github.com/Sonmezturk/telegram-bot-formatter/db"
	"github.com/Sonmezturk/telegram-bot-formatter/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)
  

func SaveOrder (orders []structs.Order, timestamp time.Time, userName string, originalFilename string) error  {
	clientInstance, clientInstanceError := db.GetMongoClient();
	if clientInstanceError != nil {
		return clientInstanceError
	}

	collection := clientInstance.Database("ketuna").Collection("orders")
	formatedOrder := structs.FormatedOrder{
		Orders: orders,
		Timestamp: timestamp,
		UserName:  userName,
		FileName: originalFilename,
	}
	
	_, err := collection.InsertOne(context.TODO(), formatedOrder )
	if err != nil {
		return err
	}

	return nil
}

func GetFilteredAggregatedData(from, to time.Time, fileName string) ([]bson.M, error) {
	client, _ := db.GetMongoClient();

	collection := client.Database("ketuna").Collection("orders")

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"timestamp": bson.M{
				"$gte": from,
				"$lte": to,
			},
			"filename": fileName,
		}}},
		{{"$unwind", "$orders"}},
		{{"$unwind", "$orders.items"}},
		{{"$group", bson.M{
			"_id": bson.M{
				"skuname": "$orders.items.skuname",
				"color":   "$orders.items.color",
				"size":    "$orders.items.customizations.size",
			},
			"totalQuantity": bson.M{"$sum": bson.M{"$toInt": "$orders.items.quantity"}},
			"orders":        bson.M{"$push": "$orders.orderid"},
			"users":         bson.M{"$addToSet": "$username"},
			"fileNames":     bson.M{"$addToSet": "$filename"},
		}}},
		{{"$sort", bson.M{"totalQuantity": -1}}},
		{{"$project", bson.M{
			"skuname":       "$_id.skuname",
			"color":         "$_id.color",
			"size":          "$_id.size",
			"totalQuantity": 1,
			"orders":        1,
			"users":         1,
			"fileNames":     1,
			"_id":           0,
		}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Panic(err)
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}