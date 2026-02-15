package httpx

import (
	"errors"
	"net/http"
	"rest_waka/pkg/res"
)

type ErrMap struct {
	Err     error
	Status  int
	Message any
}

func WriteMappedError(w http.ResponseWriter, err error, m []ErrMap, defaultStatus int, defaultMessage any) {
	for _, x := range m {
		if errors.Is(err, x.Err) {
			res.Json(w, x.Message, x.Status)
			return
		}
	}
	res.Json(w, defaultMessage, defaultStatus)
}
