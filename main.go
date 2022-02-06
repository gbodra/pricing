package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	a := App{}
	a.Initialize()

	migrations, err := strconv.ParseBool(os.Getenv("MIGRATIONS"))
	if err != nil {
		log.Println(err)
	}
	log.Println(migrations)

	if migrations {
		log.Println("Running migrations...")
		a.RunMigrations()
	}

	a.Run()
}
