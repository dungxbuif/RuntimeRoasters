package errs

import (
	"encoding/json"
	"net/http"
)

// problemRenderer writes Problem as application/problem+json (RFC 9457)
type problemRenderer struct {
	p Problem
}

func (r problemRenderer) Render(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(r.p)
}

func (r problemRenderer) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", problemContentType)
}
