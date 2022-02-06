package migrations

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gbodra/pricing-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

func readJson(filepath string) []byte {
	jsonFile, err := os.Open(filepath)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

func generateUserId(users *model.Users) {
	for i := range users.Users {
		users.Users[i].ID = primitive.NewObjectID()
	}
}

func CreateUsers() {
	var users model.Users
	byteValue := readJson("./migrations//users.json")
	json.Unmarshal(byteValue, &users)
	generateUserId(&users)

	// Needed to convert []Users{} to []Interface{}
	usersInterface := make([]interface{}, len(users.Users))
	for i := range users.Users {
		usersInterface[i] = users.Users[i]
	}

	collection := MongoClient.Database("pricing").Collection("users")
	result, err := collection.InsertMany(context.TODO(), usersInterface)

	if err != nil {
		log.Println(err)
	}

	log.Println("Users:", result)
}

func CreatePrices() {
	data := model.Price{
		ID:    primitive.NewObjectID(),
		Name:  "TV",
		Price: 100.,
	}

	collection := MongoClient.Database("pricing").Collection("products")
	result, err := collection.InsertOne(context.TODO(), data)

	if err != nil {
		log.Println(err)
	}

	log.Println("Prices:", result)
}
