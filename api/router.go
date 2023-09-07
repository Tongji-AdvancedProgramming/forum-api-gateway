package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Router() chi.Router {
	r := chi.NewRouter()

	// 健康检查服务
	r.Get("/check", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Service OK"))
	})

	r.Mount("/user", RouteUserService())

	return r
}
