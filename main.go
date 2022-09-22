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
	"github.com/robfig/cron/v3"
	"peocchiproject.it/m/api"
)

//go:embed public/*
var vue embed.FS

var ctx = context.Background()
var RedisCtx = context.Background()

type redisCtxKey interface{}

func main() {
	// :D
	splash(os.Getenv("PORT"))

	// init Gin router
	r := gin.Default()

	// redis top-level client, no pooling
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	key := redisCtxKey("client")
	RedisCtx = context.WithValue(RedisCtx, key, rdb)

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

	// API route
	apiRoute := r.Group("/api")
	api.ApplyMessagesRoute(apiRoute, &RedisCtx)

	// Background scheduler
	go func() {
		scheduler := cron.New()
		scheduler.AddFunc("0 0 * * *", LTrim)
		scheduler.Run()
	}()

	// run blocking server
	r.Run()
}
