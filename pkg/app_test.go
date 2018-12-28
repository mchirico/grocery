package pkg

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {

	a := App{}
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
	s.ProductName = "Turkey Test Parts"
	s.Price = 1323
	s.SaleDate = time.Now()

	a.AddItem(ctx, s)
	a.DeleteMany(ctx, bson.M{"product_name": "Turkey Test Parts"})

	if a.DeleteResult.DeletedCount != 1 {
		t.FailNow()
	}
	log.Printf("Delete count: %d\n", a.DeleteResult.DeletedCount)

}

func TestFind(t *testing.T) {

	type B struct {
		A string `bson:"a_final""`
	}

	type A struct {
		B    B      `bson:"b"`
		Name string `bson:"name"`
	}

	type PriceStruct struct {
		Price int `bson:"price"`
		Tree  A   `bson:"a"`
	}

	type Test struct {
		ProductName string        `bson:"product_name"`
		Price       int           `bson:"price"`
		SaleDate    time.Time     `bson:"sale_date"`
		Prices      []PriceStruct `bson:"prices"`
	}

	a := App{}
	a.CollectionName = "numbers"
	ctx, cancel := a.Initilize()

	defer cancel()

	tt := Test{}
	p := PriceStruct{Price: 34}
	p.Price = 340

	tt.ProductName = "Stuff"
	tt.SaleDate = time.Now()
	tt.Price = 3045
	tt.Prices = []PriceStruct{PriceStruct{Price: 34}, PriceStruct{Price: 21}}

	a.AddItem(ctx, tt)
	tt.ProductName = "Stuff2"
	tt.Price = 890
	a.AddItem(ctx, tt)

	cur, err := a.Collection.Find(ctx, bson.M{
		"price": bson.M{"$gte": 500}})

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	records := make([]Test, 0)

	for cur.Next(ctx) {
		//var result bson.M
		var result Test
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("result.ProductName: %v\n", result.ProductName)
		records = append(records, result)

	}

	log.Printf("records[0].Prices[1].Price: %d\n",
		records[0].Prices[1].Price)

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	a.DeleteMany(ctx, bson.M{"product_name": bson.M{"$regex": "^Stuff.*"}})
	log.Printf("Delete count: %d\n", a.DeleteResult.DeletedCount)

	if a.DeleteResult.DeletedCount != 2 {
		t.FailNow()
	}

}
