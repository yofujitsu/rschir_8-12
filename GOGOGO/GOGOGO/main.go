package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"GOGOGO/cookies"
)

// Declare a global variable to hold the secret key.
var secretKey []byte

func main() {
	var err error

	secretKey, err = hex.DecodeString("13d6b4dff8f84a10851021ecf814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/set", setCookieHandler)
	mux.HandleFunc("/get", getCookieHandler)

	log.Print("Listening...")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

type RequestData struct {
	Username string
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	var requestData RequestData
	erri := json.NewDecoder(r.Body).Decode(&requestData)
	if erri != nil {
		log.Println(erri)
	}
	cookieName := "cookie1_" + time.Now().Format("20060102150405")
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    requestData.Username,
		Path:     "/",
		MaxAge:   3000,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	err := cookies.WriteSigned(w, cookie, secretKey)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Timestamp:", requestData.Username)
	w.Write([]byte("cookie set!"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {

	value, err := cookies.ReadSigned(r, "exampleCookie", secretKey)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		case errors.Is(err, cookies.ErrInvalidValue):
			http.Error(w, "invalid cookie", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	w.Write([]byte(value))
}
