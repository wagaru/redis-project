package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wagaru/redis-project/delivery"
	"github.com/wagaru/redis-project/repository"
	"github.com/wagaru/redis-project/usecase"
)

func main() {
	log.Println("start...")
	defer log.Println("end...")

	redisUserRepo, err := repository.NewRedisUserRepo("localhost:6379")
	if err != nil {
		fmt.Println("Connect redis failed", err)
		return
	}

	redisPostRepo, err := repository.NewRedisPostRepo("localhost:6379")
	if err != nil {
		fmt.Println("Connect redis failed", err)
		return
	}

	userusecase := usecase.NewUserUsecase(redisUserRepo)
	postusecase := usecase.NewPostUsecase(redisPostRepo)
	mux := http.NewServeMux()
	_delivery := delivery.NewHttpDelivery(mux, userusecase, postusecase)

	err = _delivery.Run(":9999")
	if err != nil {
		fmt.Printf("Run server failed with reason: %v\n", err)
	}
}
