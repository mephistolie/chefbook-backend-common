package firebase

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"time"
)

const (
	signInEmailEndpoint = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword"
	collectionUsers     = "users"
)

type Profile struct {
	Id                string
	Email             string
	Username          *string
	IsPremium         bool
	CreationTimestamp time.Time
}

func (c *Client) GetProfile(context context.Context, userId string) (Profile, error) {
	record, err := c.auth.GetUser(context, userId)
	if err != nil {
		return Profile{}, err
	}
	return c.collectProfileData(*record)
}

func (c *Client) GetProfileByEmail(context context.Context, email string) (Profile, error) {
	record, err := c.auth.GetUserByEmail(context, email)
	if err != nil {
		return Profile{}, err
	}
	return c.collectProfileData(*record)
}

func (c *Client) collectProfileData(record auth.UserRecord) (Profile, error) {
	profile := Profile{
		Id:                record.UID,
		Email:             record.Email,
		CreationTimestamp: time.UnixMilli(record.UserMetadata.CreationTimestamp),
	}

	snapshot, err := c.firestore.Collection(collectionUsers).Doc(profile.Id).Get(context.Background())
	if err != nil {
		return Profile{}, err
	}
	doc := snapshot.Data()

	if username, ok := doc["name"].(string); ok && len(username) > 0 {
		profile.Username = &username
	}

	premium, ok := doc["isPremium"].(bool)
	if ok && premium {
		profile.IsPremium = true
	}

	return profile, nil
}
