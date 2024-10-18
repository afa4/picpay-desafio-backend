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
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API PicPay API"))
	}
}
