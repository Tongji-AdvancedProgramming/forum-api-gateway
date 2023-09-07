package api

import (
	"api-gateway/proto/user"
	"context"
	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net/http"
	"time"
)

var conn *grpc.ClientConn

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Success        bool                   `json:"success"`
	Msg            string                 `json:"msg"`
	UserRegistered bool                   `json:"userRegistered"`
	User           map[string]interface{} `json:"user"`
}

// Login godoc
//
//		@Summary		登录
//		@Description	完成登录操作
//		@Tags			UserService
//		@Accept			json
//		@Produce		json
//	 	@Param        	request    body     loginRequest  true  "Request Body"
//		@Success		200		{object} 	loginResponse
//		@Router			/user/login [post]
func login(username string, password string) loginResponse {
	req := user.LoginRequest{
		Username: username,
		Password: password,
	}
	c := user.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := c.Login(ctx, &req)

	if err != nil {
		log.Printf("Call Failed: %v", err)
		return loginResponse{
			Success:        false,
			Msg:            "gRPC连接异常",
			UserRegistered: false,
		}
	}

	var userObj map[string]interface{}
	err = jsoniter.UnmarshalFromString(r.UserSerialized, &userObj)
	if err != nil {
		return loginResponse{
			Success:        false,
			Msg:            "gRPC用户数据反序列化异常",
			UserRegistered: false,
		}
	}

	return loginResponse{
		Success:        r.Success,
		Msg:            r.Msg,
		UserRegistered: r.UserRegistered,
		User:           userObj,
	}
}

func RouteUserService() chi.Router {
	newConn, err := grpc.Dial("localhost:7001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connect Sample RPC Server Failed: %v", err)
	}
	conn = newConn

	r := chi.NewRouter()

	r.Post("/login", func(writer http.ResponseWriter, request *http.Request) {

		var req loginRequest
		bodyBytes, _ := io.ReadAll(request.Body)
		err = jsoniter.Unmarshal(bodyBytes, &req)

		if err != nil {
			writer.WriteHeader(400)
			_, _ = writer.Write([]byte("错误的请求格式"))
			return
		}

		rsp := login(req.Username, req.Password)

		rspBytes, _ := jsoniter.Marshal(rsp)
		_, _ = writer.Write(rspBytes)

	})

	return r
}
