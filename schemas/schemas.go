package schemas

type IngredientQuerySchema struct {
	Ingredients []string `query:"ingredients"`
}
