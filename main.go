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

	redisRepo, err := repository.NewRedisRepo("localhost:6379")
	if err != nil {
		fmt.Println("Connect redis failed", err)
		return
	}

	_usecase := usecase.NewUsecase(redisRepo)
	mux := http.NewServeMux()
	_delivery := delivery.NewHttpDelivery(mux, _usecase)

	err = _delivery.Run(":9999")
	if err != nil {
		fmt.Printf("Run server failed with reason: %v\n", err)
	}
}
