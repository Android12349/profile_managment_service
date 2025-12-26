package profile_management_storage

// Users table constants
const (
	usersTableName          = "users"
	usersIDColumn           = "id"
	usersUsernameColumn     = "username"
	usersPasswordHashColumn = "password_hash"
	usersHeightColumn       = "height"
	usersWeightColumn       = "weight"
	usersBJUColumn          = "bju"
	usersBudgetColumn       = "budget"
	usersPreferencesColumn  = "preferences"
	usersCreatedAtColumn    = "created_at"
)

// Products table constants
const (
	productsTableName       = "products"
	productsIDColumn        = "id"
	productsUserIDColumn    = "user_id"
	productsNameColumn      = "name"
	productsCaloriesColumn  = "calories"
	productsProteinColumn   = "protein"
	productsFatColumn       = "fat"
	productsCarbsColumn     = "carbs"
	productsCreatedAtColumn = "created_at"
)

// Meals table constants
const (
	mealsTableName        = "meals"
	mealsIDColumn         = "id"
	mealsUserIDColumn     = "user_id"
	mealsNameColumn       = "name"
	mealsProductIDsColumn = "product_ids"
	mealsCreatedAtColumn  = "created_at"
)
