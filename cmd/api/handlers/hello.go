package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello exec")
	message, err := io.ReadAll(r.Body)
	if err != nil {
		h.l.Println("Error on hello")
		return
	}

	fmt.Fprintf(w, "%s", message)
}
