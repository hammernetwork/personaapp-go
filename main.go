package main

import (
	"log"

	"personaapp/cmd"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		err := logger.Sync() // flushes buffer, if any
		if err != nil {
			log.Println(err) // todo think about errors mapper/parser service
		}
	}()

	sugar := logger.Sugar()
	if err := cmd.Run(); err != nil {
		sugar.Fatal(err)
	}
}
