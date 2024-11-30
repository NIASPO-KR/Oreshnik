package writer

import (
	"encoding/json"
	"net/http"

	"static/pkg/http/header"
)

func WriteJson(w http.ResponseWriter, data any) {
	header.AddJSONContentType(w.Header())

	encoder := json.NewEncoder(w)

	if err := encoder.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
