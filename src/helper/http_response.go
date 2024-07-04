package helper

import (
	"net/http"

	"github.com/davidyunus/sawitpro-estate/src/domain"
)

type HttpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

func Response(code int, message string, data, errors interface{}) HttpResponse {
	res := HttpResponse{}
	res.Code = code
	res.Message = message
	res.Data = data
	res.Errors = errors

	return res
}

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusCreated
	}

	switch err.Error() {
	case domain.ErrInvalidInput.Error():
		return http.StatusBadRequest
	case domain.ErrMaxSizeEstate.Error():
		return http.StatusBadRequest
	case domain.ErrEstateNotFound.Error():
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
