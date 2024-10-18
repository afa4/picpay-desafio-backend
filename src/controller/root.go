package controller

import (
	"net/http"
)

type RootController struct {
}

func NewRootController() *RootController {
	return &RootController{}
}

func (c *RootController) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello World!"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
