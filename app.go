package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gbodra/pricing-api/controller"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type App struct {
	Router *mux.Router
	Port   string
	Redis  *redis.Client
}

func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
		log.Println("$PORT not set. Falling back to default " + port)
	}

	return port
}

func (a *App) Initialize() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env")
	}

	a.Redis = redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "",
		DB:       0,
	})

	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.injectRedisClient()
}

func (a *App) Run() {
	port := getPort()
	log.Fatal(http.ListenAndServe(":"+port, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/health", controller.HealthCheck).Methods("GET")
	a.Router.HandleFunc("/price", controller.GetPrice).Methods("GET")
	a.Router.HandleFunc("/price", controller.SetPrice).Methods("POST")
}

func (a *App) injectRedisClient() {
	controller.RedisClient = a.Redis
}
