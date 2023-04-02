package schema

// IngredientQuery represents ingredients query params
type IngredientQuery struct {
	Ingredients []string `query:"ingredients"`
}

// User models inputs admin user has to provide to create new user.
type User struct {
	Username string `json:"username" extensions:"x-order=1"`
	Password string `json:"password" extensions:"x-order=2"`
	IsAdmin  bool   `json:"admin" extensions:"x-order=3"`
}

// Login models inputs user has to provide to log in.
type Login struct {
	Username string `json:"username" example:"username" minLength:"3" extensions:"x-order=1"`
	Password string `json:"password" example:"password" minLength:"4" extensions:"x-order=2"`
}

// Ingredient models inputs user has to provide to create an ingredient.
type Ingredient struct {
	Name string `json:"name"`
}

// Password models inputs user has to provide to update its password.
type Password struct {
	Password string `json:"password" minLength:"4"`
}

// Recipe models inputs user has to provide to create recipe
type Recipe struct {
	Name        string       `json:"name" extensions:"x-order=1"`
	Ingredients []Ingredient `json:"ingredients" minLength:"1" extensions:"x-order=2"`
}
