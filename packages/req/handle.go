package req

import (
	"golang/packages/res"
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, req *http.Request) (*T, error) {
	body, err := Decode[T](req.Body)
	if err != nil {
		res.Json(*w, err.Error(), 402)
		return nil, err
	}
	err = IsValid[T](body)
	if err != nil {
		res.Json(*w, err.Error(), 402)
		return nil, err
	}
	return &body, nil
}
