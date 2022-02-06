package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gbodra/pricing-api/controller"
	"github.com/gbodra/pricing-api/migrations"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Router *mux.Router
	Port   string
	Redis  *redis.Client
	Mongo  *mongo.Client
}

func (a *App) Initialize() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env")
	}

	a.Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS"),
		Password: "",
		DB:       0,
	})

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO"))
	a.Mongo, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.injectClients()
}

func (a *App) Run() {
	port := getPort()
	log.Fatal(http.ListenAndServe(":"+port, a.Router))
}

func (a *App) RunMigrations() {
	migrations.CreateUsers()
	migrations.CreatePrices()
}

func (a *App) initializeRoutes() {
	// Management routes
	a.Router.HandleFunc("/health", controller.HealthCheck).Methods("GET")

	// Auth routes
	a.Router.HandleFunc("/signin", controller.Signin).Methods("POST")
	a.Router.HandleFunc("/refreshtoken", controller.RefreshToken).Methods("GET")

	// App routes
	a.Router.HandleFunc("/price/{id}", controller.GetPrice).Methods("GET")
}

func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
		log.Println("$PORT not set. Falling back to default " + port)
	}

	return port
}

func (a *App) injectClients() {
	controller.RedisClient = a.Redis
	controller.MongoClient = a.Mongo
	migrations.MongoClient = a.Mongo
}
