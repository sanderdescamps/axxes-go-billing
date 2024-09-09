package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type httpResult[T any] struct {
	statusCode int
	msg        string
	changed    bool
	failed     bool
	results    []T
}

type HttpResult[T any] interface {
	Changed() HttpResult[T]
	Status(int) HttpResult[T]
	Msg(m string, a ...any) HttpResult[T]
	AddResult(item ...T) HttpResult[T]
	Error(msg string, status int) HttpResult[T]
	MarshalJSON() ([]byte, error)
	Send(w http.ResponseWriter)
}

func NewResult[T any]() HttpResult[T] {
	new := httpResult[T]{statusCode: 200, msg: "", changed: false, failed: false}
	var result HttpResult[T] = &new
	return result
}

func NewResultForResults[T any](items ...T) HttpResult[T] {
	result := NewResult[T]()
	result.AddResult(items...)
	return result
}

func (r *httpResult[T]) Changed() HttpResult[T] {
	r.changed = true
	var result HttpResult[T] = r
	return result
}

func (r *httpResult[T]) Status(c int) HttpResult[T] {
	r.statusCode = c
	var result HttpResult[T] = r
	return result
}

func (r *httpResult[T]) Msg(m string, a ...any) HttpResult[T] {
	r.msg = fmt.Sprintf(m, a...)
	var result HttpResult[T] = r
	return result
}

func (r *httpResult[T]) AddResult(item ...T) HttpResult[T] {
	if r.results == nil {
		r.results = []T{}
	}
	r.results = append(r.results, item...)
	var result HttpResult[T] = r
	return result
}

func (r *httpResult[T]) Error(msg string, status int) HttpResult[T] {
	r.statusCode = status
	r.failed = true
	r.msg = msg
	var result HttpResult[T] = r
	return result
}

func (r *httpResult[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		StatusCode int    `json:"status_code"`
		Msg        string `json:"msg,omitempty"`
		Changed    bool   `json:"changed"`
		Failed     bool   `json:"failed"`
		Results    []T    `json:"results,omitempty"`
	}{
		StatusCode: r.statusCode,
		Msg:        r.msg,
		Changed:    r.changed,
		Failed:     r.failed,
		Results:    r.results,
	})
}

func (r *httpResult[T]) Send(w http.ResponseWriter) {
	w.WriteHeader(r.statusCode)
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(r)
	// json.NewEncoder(w).Encode(r)
}
