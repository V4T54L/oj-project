package main

import (
	"net/http"
	"time"
)

// Helper functions for password hashing
func hashPassword(pw string) (string, error) {
	// Just a mock hash function, use bcrypt or a better hashing algorithm in production.
	return pw + "_hash", nil
}

// func comparePassword(hash, pw string) error {
// 	if hash != pw+"_hashed" {
// 		return fmt.Errorf("wrong password")
// 	}
// 	return nil
// }

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}
