package handler

import (
	"net/http"
)

func GetToken(r *http.Request) string {
	t := r.URL.Query().Get("token")
	if t != "" {
		return t
	}
	t = r.Header.Get("Authorization")
	if t != "" {
		return t
	}

	if err := r.ParseMultipartForm(32 << 20); err == nil {
		if v := r.FormValue("token"); v != "" {
			return v
		}
	}
	return ""
}
