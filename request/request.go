package request

import (
	"encoding/json"
	"net/http"
	"log"
	"github.com/pkg/errors"
)

func Decode(r *http.Request, w http.ResponseWriter, v interface{}) error{
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	if err != nil {
		RespondWithError(w, err, "Error occurred in parsing request.", http.StatusBadRequest)
		return err
	}
	return nil
}

func RespondWithError(w http.ResponseWriter, err error, message string, code int){
	log.Printf("%+v", errors.Wrap(err, message))
	http.Error(w, message, code)
}
