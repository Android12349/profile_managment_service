package models

type BJU struct {
	Protein int32 `json:"protein"`
	Fat     int32 `json:"fat"`
	Carbs   int32 `json:"carbs"`
}

type User struct {
	ID           int32  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Height       *int32 `json:"height,omitempty"`
	Weight       *int32 `json:"weight,omitempty"`
	BJU          *BJU   `json:"bju,omitempty"`
	Budget       *int32  `json:"budget,omitempty"`
	Preferences  string `json:"preferences,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type Product struct {
	ID        int32  `json:"id"`
	UserID    int32  `json:"user_id"`
	Name      string `json:"name"`
	Calories  *int32 `json:"calories,omitempty"`
	Protein   *int32 `json:"protein,omitempty"`
	Fat       *int32 `json:"fat,omitempty"`
	Carbs     *int32 `json:"carbs,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type Meal struct {
	ID         int32   `json:"id"`
	UserID     int32   `json:"user_id"`
	Name       string  `json:"name"`
	ProductIDs []int32 `json:"product_ids"`
	CreatedAt  string  `json:"created_at,omitempty"`
}

