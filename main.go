package main

import (
	"api-gateway/api"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	fmt.Println("Server now starting...")

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("This service provide no frontend interface"))
	})

	r.Mount("/api", api.Router())

	_ = http.ListenAndServe(":8080", r)
}
