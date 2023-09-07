package main

import (
	"api-gateway/api"
	_ "api-gateway/docs"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// @title						同济高程论坛
// @version					1.0
// @description				同济高程论坛 API
// @host						localhost:8080
// @BasePath					/api/
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				请在这里填入JWT，以供测试用途
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

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	r.Mount("/api", api.Router())

	fmt.Println("Server started. Check Swagger at: http://localhost:8080/swagger/index.html")

	_ = http.ListenAndServe(":8080", r)
}
