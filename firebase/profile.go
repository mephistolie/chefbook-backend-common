package firebase

import (
	"context"
	"time"
)

const (
	signInEmailEndpoint = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword"
	collectionUsers     = "users"
)

type Profile struct {
	Email             string
	Username          *string
	IsPremium         bool
	CreationTimestamp time.Time
}

func (c *Client) GetProfile(userId string) (*Profile, error) {
	snapshot, err := c.firestore.Collection(collectionUsers).Doc(userId).Get(context.Background())
	if err != nil {
		return nil, err
	}
	doc := snapshot.Data()

	profile := Profile{}

	if record, err := c.auth.GetUser(context.Background(), userId); err == nil {
		profile.Email = record.Email
		profile.CreationTimestamp = time.UnixMilli(record.UserMetadata.CreationTimestamp)
	}

	if username, ok := doc["name"].(string); ok && len(username) > 0 {
		profile.Username = &username
	}

	premium, ok := doc["isPremium"].(bool)
	if ok && premium {
		profile.IsPremium = true
	}

	return &profile, nil
}
