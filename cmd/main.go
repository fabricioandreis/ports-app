package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/fabricioandreis/ports-app/internal/infra/config"
	"github.com/fabricioandreis/ports-app/internal/infra/db"
	"github.com/fabricioandreis/ports-app/internal/usecase/store"
)

func main() {
	log.Println("Running Ports service")

	config := config.Load()
	inputFile, err := os.Open("./ports.json")
	if err != nil {
		err = errors.Join(errors.New("unable to open file"), err)
		log.Fatalf(err.Error())
		return
	}

	log.Println("Storing data from input file into database")

	repoPort := db.NewPortRepository(config.RedisAddress, config.RedisPassword)
	storeUsecase := store.NewStoreUsecase(repoPort)
	storeUsecase.Store(context.Background(), inputFile)

	log.Println("Finished storing data")
}
