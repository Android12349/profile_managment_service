# Примеры JSON для API Profile Management Service

## Users API

### POST /users - Создание пользователя

**Request:**
```json
{
  "user": {
    "username": "john_doe",
    "password": "securePassword123",
    "height": 180,
    "weight": 75,
    "bju": {
      "protein": 100,
      "fat": 70,
      "carbs": 250
    },
    "budget": 3000,
    "preferences": "{\"allergies\": [\"nuts\"], \"cuisine\": [\"italian\", \"japanese\"]}"
  }
}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "username": "john_doe",
    "passwordHash": "$2a$10$...",
    "height": 180,
    "weight": 75,
    "bju": {
      "protein": 100,
      "fat": 70,
      "carbs": 250
    },
    "budget": 3000,
    "preferences": "{\"allergies\": [\"nuts\"], \"cuisine\": [\"italian\", \"japanese\"]}",
    "createdAt": "2025-12-26T15:00:00Z"
  }
}
```

### GET /users/{id} - Получение пользователя

**Request:**
```
GET /users/1
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "username": "john_doe",
    "passwordHash": "$2a$10$...",
    "height": 180,
    "weight": 75,
    "bju": {
      "protein": 100,
      "fat": 70,
      "carbs": 250
    },
    "budget": 3000,
    "preferences": "{\"allergies\": [\"nuts\"], \"cuisine\": [\"italian\", \"japanese\"]}",
    "createdAt": "2025-12-26T15:00:00Z"
  }
}
```

### PATCH /users/{id} - Обновление пользователя

**Request:**
```json
{
  "user": {
    "username": "john_doe_updated",
    "height": 182,
    "weight": 77,
    "bju": {
      "protein": 110,
      "fat": 75,
      "carbs": 260
    },
    "budget": 3500,
    "preferences": "{\"allergies\": [\"nuts\", \"dairy\"], \"cuisine\": [\"italian\"]}"
  }
}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "username": "john_doe_updated",
    "passwordHash": "$2a$10$...",
    "height": 182,
    "weight": 77,
    "bju": {
      "protein": 110,
      "fat": 75,
      "carbs": 260
    },
    "budget": 3500,
    "preferences": "{\"allergies\": [\"nuts\", \"dairy\"], \"cuisine\": [\"italian\"]}",
    "createdAt": "2025-12-26T15:00:00Z"
  }
}
```

### DELETE /users/{id} - Удаление пользователя

**Request:**
```
DELETE /users/1
```

**Response:**
```json
{}
```

---

## Products API

### POST /products - Создание продукта

**Request:**
```json
{
  "product": {
    "userId": 1,
    "name": "Куриная грудка",
    "calories": 165,
    "protein": 31,
    "fat": 3,
    "carbs": 0
  }
}
```

**Response:**
```json
{
  "product": {
    "id": 1,
    "userId": 1,
    "name": "Куриная грудка",
    "calories": 165,
    "protein": 31,
    "fat": 3,
    "carbs": 0,
    "createdAt": "2025-12-26T15:05:00Z"
  }
}
```

### GET /products?user_id={id} - Получение продуктов пользователя

**Request:**
```
GET /products?user_id=1
```

**Response:**
```json
{
  "products": [
    {
      "id": 1,
      "userId": 1,
      "name": "Куриная грудка",
      "calories": 165,
      "protein": 31,
      "fat": 3,
      "carbs": 0,
      "createdAt": "2025-12-26T15:05:00Z"
    },
    {
      "id": 2,
      "userId": 1,
      "name": "Рис",
      "calories": 130,
      "protein": 2,
      "fat": 0,
      "carbs": 28,
      "createdAt": "2025-12-26T15:06:00Z"
    }
  ]
}
```

### PATCH /products/{id} - Обновление продукта

**Request:**
```json
{
  "product": {
    "name": "Куриная грудка (филе)",
    "calories": 170,
    "protein": 32,
    "fat": 3,
    "carbs": 0
  }
}
```

**Response:**
```json
{
  "product": {
    "id": 1,
    "userId": 1,
    "name": "Куриная грудка (филе)",
    "calories": 170,
    "protein": 32,
    "fat": 3,
    "carbs": 0,
    "createdAt": "2025-12-26T15:05:00Z"
  }
}
```

### DELETE /products/{id} - Удаление продукта

**Request:**
```
DELETE /products/1
```

**Response:**
```json
{}
```

---

## Meals API

### POST /meals - Создание блюда

**Request:**
```json
{
  "meal": {
    "userId": 1,
    "name": "Курица с рисом",
    "productIds": [1, 2]
  }
}
```

**Response:**
```json
{
  "meal": {
    "id": 1,
    "userId": 1,
    "name": "Курица с рисом",
    "productIds": [1, 2],
    "createdAt": "2025-12-26T15:10:00Z"
  }
}
```

### GET /meals?user_id={id} - Получение блюд пользователя

**Request:**
```
GET /meals?user_id=1
```

**Response:**
```json
{
  "meals": [
    {
      "id": 1,
      "userId": 1,
      "name": "Курица с рисом",
      "productIds": [1, 2],
      "createdAt": "2025-12-26T15:10:00Z"
    },
    {
      "id": 2,
      "userId": 1,
      "name": "Салат овощной",
      "productIds": [3, 4, 5],
      "createdAt": "2025-12-26T15:11:00Z"
    }
  ]
}
```

### PATCH /meals/{id} - Обновление блюда

**Request:**
```json
{
  "meal": {
    "name": "Курица с рисом и овощами",
    "productIds": [1, 2, 3]
  }
}
```

**Response:**
```json
{
  "meal": {
    "id": 1,
    "userId": 1,
    "name": "Курица с рисом и овощами",
    "productIds": [1, 2, 3],
    "createdAt": "2025-12-26T15:10:00Z"
  }
}
```

### DELETE /meals/{id} - Удаление блюда

**Request:**
```
DELETE /meals/1
```

**Response:**
```json
{}
```

---

## Примечания

1. **Поля height, weight, budget, bju** - опциональные, могут быть не указаны
2. **Поля calories, protein, fat, carbs** в продуктах - опциональные
3. **preferences** - JSON строка, может содержать любые данные в формате JSON
4. **productIds** - массив ID продуктов, может быть пустым
5. Все даты в формате ISO 8601 (RFC3339)

