package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
	"github.com/robfig/cron/v3"
	"peocchiproject.it/m/api"
)

//go:embed public
var vue embed.FS
var vueFS fs.FS

var ctx = context.Background()
var RedisCtx = context.Background()

type redisCtxKey interface{}

func main() {
	// :D
	splash(os.Getenv("PORT"))
	vueFS, _ = fs.Sub(vue, "public")

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
	ws := socketio.NewServer(nil)

	// spawn a goroutine to handle pub-sub messages
	sub := rdb.Subscribe(ctx, "rt-messages")
	ch := sub.Channel()
	go Subscriber(ws, ch)

	// set socket io context to ensure data transfer to client
	ws.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	ws.OnEvent("/", "message_stack_req", func(s socketio.Conn) {
		res, _ := LRange()
		s.Emit("message_stack_res", res)
	})

	// socket io error detection goroutine
	go func() {
		if err := ws.Serve(); err != nil {
			log.Fatalf("[WS] error detected: %s\n", err)
		}
	}()
	defer ws.Close()

	// Gin router handlers
	r.Use(CORS("*"))

	// ember vue app
	r.GET("/", staticHandler("/", true))
	r.GET("/assets", staticHandler("/assets", false))

	r.GET("/socket.io/*any", gin.WrapH(ws))
	r.POST("/socket.io/*any", gin.WrapH(ws))

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
