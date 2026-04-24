package response

import (
	"net/http"
	"payment-service/constants"
	errConstant "payment-service/constants/error"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string  `json:"status"`
	Message any     `json:"message"`
	Data    any     `json:"data"`
	Token   *string `json:"token,omitempty"`
}

func ErrorFromApp(c *gin.Context, err error) {
	code := errConstant.GetErrorCode(err)
	if errConstant.ErrMapping(err) {
		Error(c, code, err, nil)
	} else {
		Error(c, http.StatusInternalServerError, errConstant.ErrInternalServerError, nil)
	}
}

func Success(c *gin.Context, code int, data any, token *string) {
	c.JSON(code, Response{
		Status:  constants.Success,
		Message: http.StatusText(code),
		Data:    data,
		Token:   token,
	})
}

func Error(c *gin.Context, code int, err error, msg *string, data ...any) {
	message := errConstant.ErrInternalServerError.Error()

	if msg != nil {
		message = *msg
	} else if err != nil {
		message = err.Error()
	}

	var responseData any
	if len(data) > 0 {
		responseData = data[0]
	}

	c.JSON(code, Response{
		Status:  constants.Error,
		Message: message,
		Data:    responseData,
	})
}
