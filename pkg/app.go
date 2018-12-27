package pkg

import (
	"context"
	"encoding/json"
	"github.com/mongodb/mongo-go-driver/mongo"
	"io/ioutil"
	"log"
	"os/user"
	"time"
)

type key string

const (
	UriKey = key("uriKey")
	DBKey  = key("dbKey")
)

type MongoStruct struct {
	URI string `json:"uri"`
	DB  string `json:"db"`
}

func readFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	return string(data), err
}

func Initilize() (context.Context, context.CancelFunc) {

	mongoSetup := MongoStruct{}
	usr, _ := user.Current()
	//file := usr.HomeDir + "/.groceryMongoDB"
	file := usr.HomeDir + "/.freeMongoDB"

	jsonData, err := readFile(file)
	if err != nil {
		log.Fatalf("Can't read file ~/.groceryMongoDB")
	}
	err = json.Unmarshal([]byte(jsonData), &mongoSetup)
	if err != nil {
		log.Fatalf("Can't parse file ~/.groceryMongoDB")
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ctx = context.WithValue(ctx, UriKey, mongoSetup.URI)
	ctx = context.WithValue(ctx, DBKey, mongoSetup.DB)
	ctx, cancel := context.WithCancel(ctx)
	return ctx, cancel
}

func ConfigDB(ctx context.Context) (*mongo.Database, error) {
	client, err := mongo.NewClient(ctx.Value(UriKey).(string))
	err = client.Connect(ctx)
	db := client.Database(ctx.Value(DBKey).(string))
	return db, err
}
