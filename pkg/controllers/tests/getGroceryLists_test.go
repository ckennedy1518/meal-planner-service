package tests

import (
	"encoding/json"
	"meal-planner-service/pkg/controllers"
	"testing"
)

func TestDecodeIngredientRelationFromSnakeCaseJSON(t *testing.T) {
	payload := []byte(`[{"ingredient":{"name":"Milk","is_staple":true}}]`)

	var decoded []controllers.GroceryListItemResponse
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !decoded[0].Ingredient.IsStaple {
		t.Fatalf("expected ingredient is_staple to decode to true")
	}
}

func TestBuildGroceryListDTOs(t *testing.T) {
	groceryLists := []controllers.GroceryListResponse{
		{
			DatePlannedToGo: "2026-08-04T00:00:00Z",
			GroceryListItems: []controllers.GroceryListItemResponse{
				{
					Amount:          3,
					AmountPurchased: 1,
					Unit:            "lb",
					Ingredient:      controllers.IngredientRelation{Name: "Milk", IsStaple: false},
				},
				{
					Amount:          2,
					AmountPurchased: 0,
					Unit:            "ea",
					Ingredient:      controllers.IngredientRelation{Name: "Eggs", IsStaple: true},
				},
			},
		},
	}

	got := controllers.BuildGroceryListDTOs(groceryLists)
	want := []controllers.GroceryListDTO{{
		Date: "2026-08-04T00:00:00Z",
		Ingredients: []controllers.GroceryListIngredientDTO{
			{Name: "Milk", Quantity: 3, QuantityPurchased: 1, Unit: "lb", IsStaple: false},
			{Name: "Eggs", Quantity: 2, QuantityPurchased: 0, Unit: "ea", IsStaple: true},
		},
	}}

	if len(got) != len(want) {
		t.Fatalf("expected %d items, got %d", len(want), len(got))
	}

	for i := range want {
		if got[i].Date != want[i].Date || len(got[i].Ingredients) != len(want[i].Ingredients) {
			t.Fatalf("item %d mismatch: got %+v want %+v", i, got[i], want[i])
		}
		for j := range want[i].Ingredients {
			if got[i].Ingredients[j] != want[i].Ingredients[j] {
				t.Fatalf("ingredient %d mismatch: got %+v want %+v", j, got[i].Ingredients[j], want[i].Ingredients[j])
			}
		}
	}
}
