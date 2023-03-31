package schema

type IngredientQuerySchema struct {
	Ingredients []string `query:"ingredients"`
}

type CreateUserSchema struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"admin"`
}

type LoginSchema struct {
	Username string `json:"username"`
	Pasword  string `json:"password"`
}

type PasswordSchema struct {
	Password string `json:"password"`
}
