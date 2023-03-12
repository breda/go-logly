package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func handleHttpError(err error, w http.ResponseWriter, logger *zerolog.Logger) {
	if err != nil {
		logger.Err(err).Send()

		var buf bytes.Buffer
		err := HttpError{
			Error: err.Error(),
			Time:  time.Now(),
		}

		json.NewEncoder(&buf).Encode(&err)
		http.Error(w, buf.String(), http.StatusBadRequest)
	}
}
