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
