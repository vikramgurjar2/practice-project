package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vikramgurjar2/practice-project/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup
	//setup routes we will use inbuilt router mathods for routing and path
	router := http.NewServeMux()
	router.HandleFunc("GET /api/students", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to the students api"))
	})

	//setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	//start server
	fmt.Println("server started at", cfg.Addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("error starting server:", err)
	}

	fmt.Println("server started at", cfg.Addr)

}
