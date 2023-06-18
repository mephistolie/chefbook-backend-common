package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	collectionRecipes = "recipes"
)

type Recipe struct {
	Name              string
	IsFavourite       bool
	Categories        []string
	Servings          *int64
	Time              *string
	Calories          *int64
	Ingredients       []Ingredient
	Cooking           []Step
	CreationTimestamp *time.Time
}

type Ingredient struct {
	Text    string
	Section bool
}

type Step struct {
	Text    string
	Section bool
}

func (c *Client) GetRecipes(localId string) ([]Recipe, error) {
	iterator := c.firestore.Collection(collectionUsers).Doc(localId).Collection(collectionRecipes).Documents(context.Background())
	snapshots, err := iterator.GetAll()
	if err != nil {
		return []Recipe{}, err
	}

	var recipes []Recipe

	for _, snapshot := range snapshots {
		if recipe, err := parseRecipe(snapshot); err == nil {
			recipes = append(recipes, recipe)
		}
	}

	return recipes, nil
}

func parseRecipe(snapshot *firestore.DocumentSnapshot) (Recipe, error) {
	var ok bool
	var err error
	doc := snapshot.Data()
	recipe := Recipe{}

	if recipe.Name, ok = doc["name"].(string); !ok {
		return Recipe{}, errors.New(fmt.Sprintf("error during get name of recipe"))
	}
	recipe.IsFavourite, _ = doc["favourite"].(bool)
	if servings, ok := doc["servings"].(int64); ok && servings > 0 {
		recipe.Servings = &servings
	}
	if cookingTime, ok := doc["time"].(string); ok && len(cookingTime) > 0 {
		recipe.Time = &cookingTime
	}
	if calories, ok := doc["calories"].(int64); ok && calories > 0 {
		recipe.Calories = &calories
	}
	if creationDate, ok := doc["creationDate"].(time.Time); ok {
		recipe.CreationTimestamp = &creationDate
	}

	if categories, ok := doc["categories"].([]interface{}); ok {
		for _, category := range categories {
			recipe.Categories = append(recipe.Categories, category.(string))
		}
	}

	if recipe.Ingredients, err = parseIngredients(doc); err != nil {
		return Recipe{}, err
	}
	if recipe.Cooking, err = parseCooking(doc); err != nil {
		return Recipe{}, err
	}

	return recipe, nil
}

func parseIngredients(doc map[string]interface{}) ([]Ingredient, error) {
	rawIngredients, ok := doc["ingredients"].([]interface{})
	if !ok {
		return []Ingredient{}, errors.New(fmt.Sprintf("error during parse recipe ingredients"))
	}

	var ingredients []Ingredient
	for _, rawIngredient := range rawIngredients {
		if ingredient, ok := parseIngredient(rawIngredient); ok {
			ingredients = append(ingredients, ingredient)
		}
	}

	return ingredients, nil
}

func parseIngredient(item interface{}) (Ingredient, bool) {
	var text string
	var section bool
	var ok bool
	rawIngredient, ok := item.(map[string]interface{})
	if !ok {
		return Ingredient{}, false
	}

	rawName := rawIngredient["item"]
	if rawName == nil {
		rawName = rawIngredient["text"]
		if rawName == nil {
			return Ingredient{}, false
		}
	}
	if text, ok = rawName.(string); !ok {
		return Ingredient{}, false
	}

	rawSection := rawIngredient["selected"]
	if rawSection == nil {
		rawSection = rawIngredient["section"]
	}
	section, _ = rawSection.(bool)

	return Ingredient{
		Text:    text,
		Section: section,
	}, true
}

func parseCooking(doc map[string]interface{}) ([]Step, error) {
	rawCooking, ok := doc["cooking"].([]interface{})
	if !ok {
		return []Step{}, errors.New(fmt.Sprintf("error during parse recipe ingredients"))
	}

	var cooking []Step
	for _, rawStep := range rawCooking {
		if step, ok := parseStep(rawStep); ok {
			cooking = append(cooking, step)
		}
	}

	return cooking, nil
}

func parseStep(item interface{}) (Step, bool) {
	var text string
	var section bool
	var ok bool
	rawStep, ok := item.(map[string]interface{})
	if ok {
		if text, ok = rawStep["item"].(string); !ok {
			return Step{}, false
		}
		section, _ = rawStep["selected"].(bool)
	} else {
		if text, ok = item.(string); !ok {
			return Step{}, false
		}
	}

	return Step{
		Text:    text,
		Section: section,
	}, true
}
