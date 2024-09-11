package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Create a Signin handler
func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Get the json body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// if the structure of the body is wrong, return log that error and http error
		log.Printf("Decode.creds: %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory app
	expectedPassword, ok := users[creds.Username]

	// If a password exist for a given user
	// AND, if it is the same as the password we received, then we can move ahead ->>
	// if NOT, then we return an unauthorized error
	if !ok || expectedPassword != creds.Password {
		log.Print("error occured getting expected password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 1 minutes
	expirationTime := time.Now().Add(1 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	jwtClaims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("error in creating tokenString: %v", err.Error())
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt-token",
		Value:   tokenString,
		Expires: expirationTime,
	})

}

// In this example, the jwtKey variable is used as the secret key for the JWT
// signature. This key should be kept secure on the server,
//  and should not be shared with anyone outside of the server.
// Normally, this is stored in a configuration file, and not in the source code.
// We are using a hardcoded value here for simplicity.

// Welcome handler is
func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the request cookies, which come with every request
	c, err := r.Cookie("jwt-token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Printf("no cookie found: %v", err.Error())
			// If the cookie is not set, return an unauthorized error 401
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Printf("r.Cookie(jwt-token): %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tokenStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// JWT Parsing: On each subsequent request, the token is sent back to the server (typically in the HTTP headers).
	// The server parses the token, extracts the user information,
	// and validates it to ensure the user is authenticated.

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Printf("jwt.ErrSignatureMethod: %v", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !token.Valid {
		log.Print("Token not valid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally, return the welcome message to the user, along with their
	// username given in the token
	w.Write([]byte(fmt.Sprintf("Welcome: %s", claims.Username)))
}

// Refresh token handler is
func Refresh(w http.ResponseWriter, r *http.Request) {
	// Implement Refresh logic
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Implement Refresh logic
}
