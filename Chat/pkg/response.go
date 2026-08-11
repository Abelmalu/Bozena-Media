package pkg

import (
	"time"

	ierrors "github.com/abelmalu/golang-posts/Chat/internal/errors"
	"github.com/gin-gonic/gin"
)


type APIResponse[T any] struct{

	Success bool `json:"success"`
	Data T `json:"data,omitempty"`
	RequestID string `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
	Error *ierrors.AppError `json:"error,omitempty"`
}




func SendSuccessResponse[T any](c *gin.Context,data T,requestID string,status int){
	
apiResponse := APIResponse[T]{
	Success: true,
	Data:data,
	RequestID: requestID,
	Timestamp: time.Now(),
}

c.JSON(status,apiResponse)


}
func SendErrorResponse[T any](c *gin.Context,err *ierrors.AppError,requestID string,status int){


	apiResponse := APIResponse[T]{
		Success:false,
		RequestID:requestID,
		Timestamp:time.Now().In(time.FixedZone("EAT",3*60*60)),
		Error:err,

	}

	c.JSON(status,apiResponse)
}