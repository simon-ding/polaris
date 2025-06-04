package server

import (
	"polaris/log"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Coder interface {
	Code() int
}


func HttpHandler(f func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		r, err := f(ctx)
		if err != nil {
			log.Errorf("url %v return error: %v", ctx.Request.URL, err)
			cc, ok := err.(Coder)
			if ok {
				ctx.JSON(200, Response{
					Code:    cc.Code(),
					Message: fmt.Sprintf("%v", err),
				})
				return
	
			}
			ctx.JSON(200, Response{
				Code:    1,
				Message: fmt.Sprintf("%v", err),
			})
			return
		}
		log.Debugf("url %v return: %+v", ctx.Request.URL, r)

		ctx.JSON(200, Response{
			Code:    0,
			Message: "success",
			Data:    r,
		})

	}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
