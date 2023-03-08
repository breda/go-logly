package http

import (
	"log"
	"net/http"
)

func handleHttpError(err error, w http.ResponseWriter) {
	if err != nil {
		log.Println("error: ", err)

		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
