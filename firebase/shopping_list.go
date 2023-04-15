package firebase

import (
	"context"
)

type ShoppingList []string

func (c *Client) GetShoppingList(localId string) (*ShoppingList, error) {
	snapshot, err := c.firestore.Collection(collectionUsers).Doc(localId).Get(context.Background())
	if err != nil {
		return nil, err
	}
	doc := snapshot.Data()

	shoppingList := ShoppingList{}

	if purchases, ok := doc["shoppingList"].([]interface{}); ok {
		for _, purchase := range purchases {
			shoppingList = append(shoppingList, purchase.(string))
		}
	}

	return &shoppingList, nil
}
