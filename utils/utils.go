package utils

import (
	"encoding/json"
	"net/http"
)

func VContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func WriteAPIJSON(w http.ResponseWriter, s int, m interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)

	return json.NewEncoder(w).Encode(m)
}

func DecodeAPIJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	return nil
}
