package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gbodra/pricing-api/model"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.Background()
var RedisClient *redis.Client
var MongoClient *mongo.Client

func GetPrice(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	cache_active, _ := strconv.ParseBool(os.Getenv("CACHE"))

	if cache_active {
		log.Printf("Cache: Active")
		val, err := RedisClient.Get(ctx, vars["id"]).Result()

		w.Header().Set("App-Cached", "True")

		if err == redis.Nil {
			log.Print("Key expired or does not exist")
			cachePrice(vars["id"])
			val, _ = RedisClient.Get(ctx, vars["id"]).Result()
			w.Header().Set("App-Cached", "False")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))
		return
	}

	log.Printf("Cache: Inactive")
	price, _ := getPriceFromMongo(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(price))
}

func getPriceFromMongo(id string) ([]byte, primitive.ObjectID) {
	id_obj, _ := primitive.ObjectIDFromHex(id)
	collection := MongoClient.Database("pricing").Collection("products")
	var result model.Price
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id_obj}}).Decode(&result)

	if err != nil {
		log.Println(err)
	}

	price, _ := json.Marshal(result)

	return price, id_obj
}

func cachePrice(id string) {
	price, id_response := getPriceFromMongo(id)

	err := RedisClient.Set(ctx, id_response.Hex(), price, time.Minute).Err()
	if err != nil {
		log.Println(err)
	}
}
