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

type App struct {
	DB              *mongo.Database
	Collection      *mongo.Collection
	CollectionName  string
	Error           error
	InsertOneResult *mongo.InsertOneResult
	DeleteResult    *mongo.DeleteResult
}

type MongoStruct struct {
	URI string `json:"uri"`
	DB  string `json:"db"`
}

func readFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	return string(data), err
}

func (a *App) Initilize() (context.Context, context.CancelFunc) {

	mongoSetup := MongoStruct{}
	usr, _ := user.Current()
	//file := usr.HomeDir + "/.groceryMongoDB"
	//file := usr.HomeDir + "/.freeMongoDB"
	file := usr.HomeDir + "/.freeMongoDB2"
	//file := usr.HomeDir + "/.aipiggybot"

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

	a.DB, a.Error = ConfigDB(ctx)
	a.Collection = a.DB.Collection("numbers")

	return ctx, cancel
}

func (a *App) AddItem(ctx context.Context, document interface{}) {
	a.InsertOneResult, a.Error = a.Collection.InsertOne(ctx, document)
}

func (a *App) DeleteMany(ctx context.Context, document interface{}) {
	a.DeleteResult, a.Error = a.Collection.DeleteMany(ctx, document)
}

func (a *App) Find(ctx context.Context, document interface{}, records interface{}) {

}

func ConfigDB(ctx context.Context) (*mongo.Database, error) {
	client, err := mongo.NewClient(ctx.Value(UriKey).(string))
	err = client.Connect(ctx)
	db := client.Database(ctx.Value(DBKey).(string))
	return db, err
}
