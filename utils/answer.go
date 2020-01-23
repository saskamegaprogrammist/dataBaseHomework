package utils

import (
	"encoding/json"
	"net/http"
)

func CreateAnswer(w http.ResponseWriter, statusCode int, value interface{})  {
	encoded , err := json.Marshal(value)
	if err != nil {
		//log.Println(err)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(encoded)
	if err != nil {
		//log.Println(err)
	}
}