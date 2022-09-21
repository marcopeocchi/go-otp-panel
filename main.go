package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"peocchiproject.it/m/api"
)

//go:embed public/*
var vue embed.FS

var ctx = context.Background()
var redisCtx = context.Background()

type redisCtxKey interface{}

func main() {
	// init Gin router
	r := gin.Default()

	// redis top-level client, no pooling
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	key := redisCtxKey("client")
	redisCtx = context.WithValue(redisCtx, key, rdb)

	defer rdb.Close()

	// init socket io instance
	io := socketio.NewServer(nil)

	// spawn a goroutine to handle pub-sub messages
	sub := rdb.Subscribe(ctx, "rt-messages")
	ch := sub.Channel()
	go Subscriber(io, ch)

	// set socket io context to ensure data transfer to client
	io.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	io.OnEvent("/", "message_stack_req", func(s socketio.Conn) {
		res, _ := LRange()
		s.Emit("message_stack_res", res)
	})

	// socket io error detection goroutine
	go func() {
		if err := io.Serve(); err != nil {
			log.Fatalf("[WS] error detected: %s\n", err)
		}
	}()
	defer io.Close()

	// Gin router handlers
	r.Use(CORS("*"))
	r.Static("/web", "./public")
	r.Static("/assets", "./public/assets")

	r.Any("/public/*f", func(ctx *gin.Context) {
		staticServer := http.FileServer(http.FS(vue))
		staticServer.ServeHTTP(ctx.Writer, ctx.Request)
	})

	r.GET("/socket.io/*any", gin.WrapH(io))
	r.POST("/socket.io/*any", gin.WrapH(io))

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/api/publish", api.PublishMessage(&redisCtx))

	// run blocking server
	r.Run()
}
