package main

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed input.json
var inputJson string

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		_, cancel := context.WithTimeout(c, 1*time.Second)

		defer cancel()

		time.Sleep(2 * time.Second)
		v := true

		if v {
			fmt.Println("pong")
			c.JSON(http.StatusOK, map[string]interface{}{
				"message": "pong",
			})
			return
		}
		sendAlert(c, fmt.Errorf("failed to process request"))
		resp := Resp{
			TransactionID:       "123",
			LinkedTransactionID: "456",
			CreatedAt:           "2021-09-01",
			RequestDocumentData: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			ResponseData: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Internal Server Error",
		})
		c.JSON(http.StatusOK, resp)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func sendAlert(ctx *gin.Context, err error) {
	_ = ctx.Error(err)
}

type Resp struct {
	TransactionID       string                 `json:"transaction_id"`
	LinkedTransactionID string                 `json:"linked_transaction_id"`
	CreatedAt           string                 `json:"created_at"`
	RequestDocumentData map[string]interface{} `json:"request_document_data"`
	ResponseData        map[string]interface{} `json:"response_document_data"`
}

type FailedError struct {
	Err error `json:"error"`
}

func (f *FailedError) Error() string {
	return fmt.Sprintf("%v", f.Err)
}
