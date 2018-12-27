package main

import (
	"context"
	"github.com/mchirico/grocery/pkg"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
	"time"
)

func main() {

	ctx, cancel := pkg.Initilize()
	defer cancel()

	db, _ := pkg.ConfigDB(ctx)
	collection := db.Collection("numbers")

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	type Sale struct {
		ProductName string    `bson:"product_name"`
		Price       int       `bson:"price"`
		SaleDate    time.Time `bson:"sale_date"`
	}

	s := Sale{}
	s.ProductName = "Turkey"
	s.Price = 1323
	s.SaleDate = time.Now()

	//res, err := collection.InsertOne(ctx, bson.M{"name": "pi",
	//	"value": "ISODate(\"2014-11-04T11:22:19.589Z\")"})

	res, err := collection.InsertOne(ctx, s)
	id := res.InsertedID

	log.Printf("id: %s\n", id)

	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.M{
		"price": bson.M{"$gte": 500}})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		//var result bson.M
		var result Sale
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("result.ProductName: %v\n", result.ProductName)
		log.Printf("result.SaleDate: %v\n", result.SaleDate)

	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	result, err := collection.DeleteMany(ctx, bson.M{
		"price": bson.M{"$gte": 500}})

	if err != nil {
		log.Printf("error deleting: %v", err)
	}

	log.Printf("Delete count: %d\n", result.DeletedCount)

	//client.Database("testing").Collection("numbers").Drop(ctx)

}
