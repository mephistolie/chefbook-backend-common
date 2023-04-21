package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"google.golang.org/api/option"
	"net/http"
	"time"
)

type Client struct {
	http      http.Client
	auth      auth.Client
	firestore firestore.Client
	apiRoute  string
}

func NewClient(credentials []byte, googleApiKey string) (*Client, error) {
	opt := option.WithCredentialsJSON(credentials)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	db, err := firebaseApp.Firestore(context.Background())
	if err != nil {
		return nil, err
	}

	apiRoute := ""
	if len(googleApiKey) > 0 {
		apiRoute = fmt.Sprintf("%s?key=%s", signInEmailEndpoint, googleApiKey)
	}

	return &Client{
		http:      http.Client{Timeout: 5 * time.Second},
		auth:      *authClient,
		firestore: *db,
		apiRoute:  apiRoute,
	}, nil
}
