package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-google-cloud-storage/app/constant"
	"math/rand"
	"strconv"
	"time"
)

type Response struct {
	ApiID   string      `json:"api_id"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
}

type ResponseMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func GenerateApiCallID() string {
	nowTime := time.Now().Unix()
	currentDate := time.Now()

	hourTm := currentDate.Hour()
	minuteTm := currentDate.Minute()
	secondTm := currentDate.Second()

	startID := strconv.FormatInt(nowTime, 10) + strconv.Itoa(hourTm) + strconv.Itoa(minuteTm) + strconv.Itoa(secondTm)
	randomNum := rand.Intn(10000000) + 1
	apiCallID := fmt.Sprintf("API_CALL_%s_%d", startID, randomNum)
	return apiCallID
}

func ResponseAPI(ctx *gin.Context, response constant.ResponseMap, additionalData ...interface{}) {

	apiID, exists := ctx.Get(constant.RequestIDKey)
	if !exists {
		apiID = GenerateApiCallID()
	}

	var data interface{}
	var meta interface{}

	if len(additionalData) > 0 {
		data = additionalData[0]
	}
	if len(additionalData) > 1 {
		meta = additionalData[1]
	}

	ctx.JSON(response.Code, Response{
		ApiID:   apiID.(string),
		Status:  response.Status,
		Message: response.Message,
		Data:    data,
		Meta:    meta,
	})
}
