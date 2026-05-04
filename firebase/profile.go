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

func (c *Client) GetProfile(ctx context.Context, userId string) (Profile, error) {
	record, err := c.auth.GetUser(ctx, userId)
	if err != nil {
		return Profile{}, err
	}
	return c.collectProfileData(ctx, *record)
}

func (c *Client) GetProfileByEmail(ctx context.Context, email string) (Profile, error) {
	record, err := c.auth.GetUserByEmail(ctx, email)
	if err != nil {
		return Profile{}, err
	}
	return c.collectProfileData(ctx, *record)
}

func (c *Client) collectProfileData(ctx context.Context, record auth.UserRecord) (Profile, error) {
	profile := Profile{
		Id:                record.UID,
		Email:             record.Email,
		CreationTimestamp: time.UnixMilli(record.UserMetadata.CreationTimestamp),
	}

	snapshot, err := c.firestore.Collection(collectionUsers).Doc(profile.Id).Get(ctx)
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
