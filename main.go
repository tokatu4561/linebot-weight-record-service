package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"tokatu4561/line-bot/record-service/handlers"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

type application struct {
}

func main() {
	_ = godotenv.Load(".env")

	log.SetFlags(log.Llongfile)

	app := &application{}
	err := app.serve()
	if err != nil {
		log.Fatalln(err)
	}
}


func(app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", "8000"),
		Handler:           routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	log.Printf("Starting Back end server in mode on port")

	return srv.ListenAndServe()
}

func routes() http.Handler  {
	mux := chi.NewRouter()

	mux.Post("/weight-regist", handlers.WeightRegist)

	return mux
}