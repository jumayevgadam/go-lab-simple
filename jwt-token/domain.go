package main

import "github.com/golang-jwt/jwt/v5"

// Create the JWT key used to create the signature
var jwtKey = []byte("gadamus's_secret_key")

// For simlpification, we're storing the users information as an memory-map in our code
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// So for now, there are only two valid users in our application: user1, and user2.
// In a real application, the user information would be stored in a database,
// and the password would be hashed and stored in a separate column. 
// We are using a hard-coded map here for simplicity.