package main

import (
	"context"
	"github.com/mchirico/grocery/pkg"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
	"testing"
	"time"
)

//var ctx context.Context
//var cancel context.CancelFunc
//
//func TestMain(m *testing.M) {
//
//	ctx, cancel = pkg.Initilize()
//	os.Exit(0)
//}

func TestInsert(t *testing.T) {
	ctx, cancel := pkg.Initilize()
	defer cancel()

	db, _ := pkg.ConfigDB(ctx)
	collection := db.Collection("numbers")

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)

	type Test struct {
		ProductName string    `bson:"product_name"`
		Price       int       `bson:"price"`
		SaleDate    time.Time `bson:"sale_date"`
	}

	s := Test{}
	s.ProductName = "Turkey Test Parts"
	s.Price = 1323
	s.SaleDate = time.Now()

	_, err := collection.InsertOne(ctx, s)

	if err != nil {
		t.FailNow()
	}

	result, err := collection.DeleteMany(ctx, bson.M{"product_name": "Turkey Test Parts"})

	if err != nil {
		log.Printf("error deleting: %v", err)
	}

	if result.DeletedCount != 1 {
		t.FailNow()
	}
	log.Printf("Delete count: %d\n", result.DeletedCount)

}
