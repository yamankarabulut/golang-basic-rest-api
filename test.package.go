package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func testFunc(c *gin.Context) {
	var params Params

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if client == nil {
	// 	var err error
	// 	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("config.mongoDBURL")))
	// 	fmt.Println("client gelmemiş")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// // Initialize the MongoDB collection if not done already
	// if collection == nil {
	// 	collection = client.Database(viper.GetString("config.dbname")).Collection("users")
	// }

	filter := bson.D{
		{Key: "$or",
			Value: bson.A{
				bson.M{"version": bson.M{"$gte": params.VersionStartValue, "$lt": params.VersionEndValue}, "language": params.Lang1},
				bson.M{"version": bson.M{"$gte": params.VersionStartValue, "$lt": params.VersionEndValue}, "language": params.Lang2},
			}},
	}

	fmt.Println(filter)

	filter2 := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "version", Value: bson.M{"$gte": params.VersionStartValue, "$lt": params.VersionEndValue}}},
				bson.D{
					{Key: "$or",
						Value: bson.A{
							bson.M{"language": params.Lang1},
							bson.M{"language": params.Lang2},
						}},
				},
			},
		},
	}
	//fmt.Println(filter2)

	cursor, err := users.Find(context.Background(), filter2)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer cursor.Close(context.Background())

	var results []Data
	for cursor.Next(context.Background()) {
		var data Data
		err := cursor.Decode(&data)
		if err != nil {
			fmt.Println(err)
			// Handle error as needed
		}
		results = append(results, data)
	}

	// Send the results as JSON
	c.JSON(http.StatusOK, results)
}
