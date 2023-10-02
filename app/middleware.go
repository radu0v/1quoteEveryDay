package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// csfr protection for forms
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	csrfHandler.SetFailureHandler(failureFunc())
	return csrfHandler
}

func failureFunc() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request Failed. Reason:%v", nosurf.Reason(r))
		http.Error(w, http.StatusText(nosurf.FailureCode), nosurf.FailureCode)
	})
}
