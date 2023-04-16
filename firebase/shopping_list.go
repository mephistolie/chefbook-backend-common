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
		for _, rawPurchase := range purchases {
			if purchase, ok := rawPurchase.(string); ok {
				shoppingList = append(shoppingList, purchase)
			}
		}
	}

	return &shoppingList, nil
}
