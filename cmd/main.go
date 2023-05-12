package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fabricioandreis/ports-app/internal/infra/config"
	"github.com/fabricioandreis/ports-app/internal/infra/db"
	"github.com/fabricioandreis/ports-app/internal/usecase/store"
)

var ErrUnableOpenFile = errors.New("unable to open file")

func main() {
	log.Println("Running Ports service")

	done := make(chan bool)
	config := config.Load()

	redis, err := db.NewClient(config.RedisAddress, config.RedisPassword)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1) // cannot proceed without a database
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go gracefullyShutdown(cancel, redis, done)

	run(ctx, config, redis, done)

	<-done
	log.Println("Shutting down Ports service")
}

func run(ctx context.Context, config config.Config, dbClient db.Client, done chan<- bool) {
	defer func() {
		done <- true
	}()

	inputFile, err := os.Open(config.InputJSONFilePath)
	if err != nil {
		err = errors.Join(ErrUnableOpenFile, err)
		log.Println(err.Error())

		return
	}

	repoPort := db.NewPortRepository(dbClient)
	storeUsecase := store.NewUseCase(repoPort)
	count, err := storeUsecase.Store(ctx, inputFile)
	log.Printf("Finished store use case after storing %v ports\n", count)

	if err != nil {
		log.Printf("Did not process all items in the input JSON file due to error: %s\n", err.Error())
	}
}

func gracefullyShutdown(cancel context.CancelFunc, dbClient db.Client, done chan bool) {
	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case <-trap:
		log.Println("Received OS signal to terminate process. Cancelling context...")
		cancel()
		<-done
	case <-done:
	}
	log.Println("Gracefully shutting down application")

	dbClient.Close()
	done <- true
}
