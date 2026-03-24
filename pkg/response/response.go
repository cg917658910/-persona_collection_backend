package response

import "github.com/gin-gonic/gin"

type Envelope struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

func OK(c *gin.Context, data interface{}, meta interface{}) {
	resp := Envelope{
		Code:    0,
		Message: "ok",
		Data:    data,
	}
	if meta != nil {
		resp.Meta = meta
	}
	c.JSON(200, resp)
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Envelope{
		Code:    status,
		Message: message,
		Data:    gin.H{},
	})
}
