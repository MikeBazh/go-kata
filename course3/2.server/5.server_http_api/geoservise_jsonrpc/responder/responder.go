package responder

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	//ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type Respond struct {
	//log *zap.Logger
	//godecoder.Decoder
}

func NewResponder() Responder {
	return &Respond{}
}

func (r *Respond) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	resp, err := json.Marshal(responseData)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(resp)
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
}

func (r *Respond) ErrorForbidden(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
}

//func (r *Respond) ErrorUnauthorized(w http.ResponseWriter, err error) {
//	r.log.Warn("http resposne Unauthorized", zap.Error(err))
//	w.Header().Set("Content-Type", "application/json;charset=utf-8")
//	w.WriteHeader(http.StatusUnauthorized)
//	if err := r.Encode(w, Response{
//		Success: false,
//		Message: err.Error(),
//		Data:    nil,
//	}); err != nil {
//		r.log.Error("response writer error on write", zap.Error(err))
//	}
//}
