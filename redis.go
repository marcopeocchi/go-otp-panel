package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
	socketio "github.com/googollee/go-socket.io"
)

type message struct {
	Uid       string    `json:"uid"`
	Message   string    `json:"message"`
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	OTP       string    `json:"otp"`
	Updated   time.Time `json:"updated"`
}

// subscriber takes a pointer to the socket io instance and the output of
// a channel of the pub-sub goroutine, foreach message update broadcast it to
// the default socketio namespace
func Subscriber(io *socketio.Server, c <-chan *redis.Message) {
	for msg := range c {
		io.BroadcastToNamespace("/", "message_update", msg.Payload)
	}
}

// lrange performs REDIS LRANGE command against the message stack, parse it to a struct
// and retrieves OTP from the message, then send the updated struct to the client
func LRange() (string, error) {
	client := RedisCtx.Value("client").(*redis.Client)
	stack, err := client.LRange(ctx, "message_stack", 0, 49).Result()

	if err != nil {
		return "", err
	}

	parsed := []message{}

	_regexp, _ := regexp.Compile(`[a-z0-9]*\d[a-z0-9]*`)

	for _, entry := range stack {
		res := message{}
		json.Unmarshal([]byte(entry), &res)

		otpMatch := _regexp.FindString(res.Message)
		res.OTP = otpMatch

		parsed = append(parsed, res)
	}

	json, err := json.Marshal(parsed)

	return string(json), err
}

// Ltrim performs REDIS LTRIM command against the message stack.
// Used as a scheduled job to chop the stack to a default length (50)
func LTrim() {
	client := RedisCtx.Value("client").(*redis.Client)
	err := client.LTrim(ctx, "message_stack", 0, 49).Err()
	fmt.Println("TRim")
	if err != nil {
		panic(err)
	}
}
