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

// function that returns the error for the no surf csrf cookie
// implemented this function because i had error bad request after
// submitting some forms and i did not know what whas the problem
// it was a typo in the go template file
func failureFunc() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request Failed. Reason:%v", nosurf.Reason(r))
		http.Error(w, http.StatusText(nosurf.FailureCode), nosurf.FailureCode)
	})
}

// function that loads and saves the session on each request
func SessionLoad(next http.Handler) http.Handler {
	return sessionManager.LoadAndSave(next)
}
