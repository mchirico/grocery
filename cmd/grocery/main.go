package main

import (
	"context"
	"github.com/mchirico/grocery/pkg"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
	"time"
)

func main() {

	a := pkg.App{}
	a.CollectionName = "numbers"
	ctx, cancel := a.Initilize()
	defer cancel()

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	type Test struct {
		ProductName string    `bson:"product_name"`
		Price       int       `bson:"price"`
		SaleDate    time.Time `bson:"sale_date"`
	}

	s := Test{}
	s.ProductName = "Turkey"
	s.Price = 1323
	s.SaleDate = time.Now()

	a.AddItem(ctx, s)

	collection := a.Collection

	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.M{
		"price": bson.M{"$gte": 500}})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		//var result bson.M
		var result Test
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
