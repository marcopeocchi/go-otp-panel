package api

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"peocchiproject.it/m/api/dto"
)

func ApplyMessagesRoute(r *gin.RouterGroup, redisCtx *context.Context) {
	r.GET("/range/:range", getRange())
	r.POST("/publish", publishMessage(redisCtx))
}

func publishMessage(redisCtx *context.Context) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		messageReq := dto.MessagePublishRequest{}
		if err := ctx.Bind(&messageReq); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		detectOTP, _ := regexp.Compile(`[a-z0-9]*\d[a-z0-9]*`)
		otp := detectOTP.FindString(messageReq.Message)
		id, _ := uuid.NewV4()

		if len(otp) <= 3 {
			otp = ""
		}

		toRedis := map[string]interface{}{
			"uid":       id.String(),
			"message":   messageReq.Message,
			"otp":       otp,
			"sender":    messageReq.Sender,
			"recipient": messageReq.Recipient,
			"updated":   time.Now(),
		}

		json, _ := json.Marshal(toRedis)

		client := (*redisCtx).Value("client").(*redis.Client)
		_, err := client.LPush(context.TODO(), "message_stack", json).Result()
		client.Publish(context.TODO(), "rt-messages", json)

		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, toRedis)
	}
}

func getRange() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusSeeOther)
	}
}
