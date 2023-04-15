package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
)

const (
	collectionRecipes = "recipes"
)

type Recipe struct {
	Name        string
	IsFavourite bool
	Categories  []string
	Servings    *int
	Time        *string
	Calories    *int
	Ingredients []Ingredient
	Cooking     []Step
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
	recipe.IsFavourite = doc["favourite"].(bool)
	if servings, ok := doc["servings"].(int); ok && servings > 0 {
		recipe.Servings = &servings
	}
	if time, ok := doc["time"].(string); ok && len(time) > 0 {
		recipe.Time = &time
	}
	if calories, ok := doc["calories"].(int); ok && calories > 0 {
		recipe.Calories = &calories
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
	rawIngredient := item.(map[string]interface{})

	rawName := rawIngredient["item"]
	if rawName == nil {
		rawName = rawIngredient["text"]
		if rawName == nil {
			return Ingredient{}, false
		}
	}
	text = rawName.(string)

	rawSection := rawIngredient["selected"]
	if rawSection == nil {
		rawSection = rawIngredient["section"]
	}
	if rawSection != nil {
		section = rawSection.(bool)
	}

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
	rawStep, ok := item.(map[string]interface{})
	if ok {
		text = rawStep["item"].(string)
		section = rawStep["selected"].(bool)
	} else {
		text = item.(string)
	}

	return Step{
		Text:    text,
		Section: section,
	}, true
}
