## Summary

Implement Todo CRUD operations with authentication

## 🎯 What Does This PR Do?

### Summary

This PR implements complete CRUD (Create, Read, Update, Delete) operations for Todo items. Users can create new todos, list all their todos, update existing todos, and delete todos. All endpoints are protected by JWT authentication middleware to ensure users can only access their own todos. The implementation follows the handler → service → repository architecture pattern with proper error handling and input validation.

### Key Changes

- Created Todo model
- Implemented database migration for `todos` table
- Added TodoRepository with methods: Create, FindByUserID, FindByID, Update, Delete
- Added TodoService with business logic for todo operations and validation
- Implemented REST handlers for all CRUD endpoints with authentication
- Added input validation for todo creation and updates

### Request Flow & Architecture Layers

#### Request Flow in This PR

Request → Auth Middleware → Todo Handler → Todo Service → Todo Repository → Database

#### Layers Implemented

- **Handler Layer** (`internal/handlers/rest/todo_handler.go`) - HTTP request/response handling, JWT token extraction, input validation, error responses
- **Service Layer** (`internal/services/todo_service.go`) - Business logic for todo operations, ownership verification, status validation
- **Repository Layer** (`internal/repository/todo_repository.go`) - Data access abstraction, GORM queries for CRUD operations

### Database & Models

#### Models Defined

- **Todo** (`internal/models/todo.go`)
  - ID: Primary key (uint)
  - UserID: Foreign key to users table (uint)
  - Title: Todo title, required, max 200 chars (string)
  - Description: Optional description, max 1000 chars (string)
  - Status: Enum ("pending", "in_progress", "completed") (string)
  - CreatedAt, UpdatedAt: Timestamps (time.Time)

#### Migrations

- `migrations/000003_create_todos_table.up.sql` - Creates todos table with indexes on user_id and status
- `migrations/000003_create_todos_table.down.sql` - Drops todos table

#### Database Operations

- Create new todo with user ownership
- List all todos for authenticated user with optional status filter
- Get single todo by ID with ownership check
- Update todo title, description, or status with ownership verification
- Soft delete todo (sets deleted_at timestamp)

## 📋 API Endpoints

### Endpoints Added/Modified

| Method | Endpoint          | Handler Function | Middleware | Description           |
| ------ | ----------------- | ---------------- | ---------- | --------------------- |
| POST   | /api/v1/todos     | `CreateTodo`     | Auth       | Create new todo       |
| GET    | /api/v1/todos     | `GetTodos`       | Auth       | List all user's todos |
| GET    | /api/v1/todos/:id | `GetTodo`        | Auth       | Get single todo by ID |
| PUT    | /api/v1/todos/:id | `UpdateTodo`     | Auth       | Update todo           |
| DELETE | /api/v1/todos/:id | `DeleteTodo`     | Auth       | Delete todo           |

## Review Focus

1. **Ownership verification** - Ensure users can only access/modify their own todos
2. **Input validation** - Check if title/description validation is sufficient
3. **Error handling** - Verify all edge cases are handled (todo not found, unauthorized access)
4. **Status enum** - Is the status validation in service layer correct?

## Testing

### Create Todo

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "title": "Learn Golang",
    "description": "Complete Todo app project",
    "status": "in_progress"
  }'

**Expected Result**:
{
  "status": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "Learn Golang",
    "description": "Complete Todo app project",
    "status": "in_progress",
    "created_at": "2024-12-03T10:30:00Z",
    "updated_at": "2024-12-03T10:30:00Z"
  }
}
```

### List Todos

```bash
curl -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

**Expected Result**: Array of all todos belonging to authenticated user
```

### Update Todo

```bash
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "status": "completed"
  }'

**Expected Result**: Updated todo with new status
```

### Delete Todo

```bash
curl -X DELETE http://localhost:8080/api/v1/todos/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

**Expected Result**:
{
  "status": "success",
  "message": "Todo deleted successfully"
}
```

### Error Cases Tested

- ✅ Creating todo without authentication → 401 Unauthorized
- ✅ Accessing other user's todo → 403 Forbidden
- ✅ Invalid status value → 400 Bad Request
- ✅ Empty title → 400 Bad Request
- ✅ Todo not found → 404 Not Found

## Notes

**Decisions Made**:

- Used soft delete (GORM's deleted_at) instead of hard delete to preserve data
- Status field is validated in service layer with allowed values: "pending", "in_progress", "completed"
- Default status is "pending" if not provided
- Todos are ordered by created_at DESC when listing

**Known Issues**:

- No pagination implemented yet for listing todos (will add in future PR if needed)
- No filtering by status yet (can be added easily)

**Questions**:

- Should we add priority field (low, medium, high) to todos?
- Do we need due_date field for todos?
