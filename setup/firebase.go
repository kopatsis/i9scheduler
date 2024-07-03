package setup

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func InitFirebase() *auth.Client {
	firebaseConfigBase64 := os.Getenv("FIREBASE_CONFIG_BASE64")
	if firebaseConfigBase64 == "" {
		log.Fatal("FIREBASE_CONFIG_BASE64 environment variable is not set.")
	}

	configJSON, err := base64.StdEncoding.DecodeString(firebaseConfigBase64)
	if err != nil {
		log.Fatalf("Error decoding FIREBASE_CONFIG_BASE64: %v", err)
	}
	opt := option.WithCredentialsJSON(configJSON)
	fmt.Println(opt)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Failed to initialize Firebase app")
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		panic("Failed to initialize Firebase auth")
	}

	return auth
}
