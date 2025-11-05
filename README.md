# Tasks & Implementation Steps

1. Setup:

- Initialize a new Go module: go mod init todo-api
- Install the recommended packages:

```bash
go get github.com/go-chi/chi/v5
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get gorm.io/gorm
go get gorm.io/driver/mysql
```

2. Database:

- Create MySQL tables manually based on the models
- Use gorm package to connect, read, and write data

```bash
docker run -d --name todo-app-db -p 3306:3306 -e MYSQL_ROOT_PASSWORD=my-password -e MYSQL_DATABASE=todo_db mysql:8.0
```

````bash
use todo_db;
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS todos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;
```

3. User Handlers:
   - Implement the /register handler: hash the password with bcrypt, save the user to the DB.
   - Implement the /login handler: find the user, compare the password with bcrypt.CompareHashAndPassword, and if successful, generate and return a JWT.
4. Middleware:
   - Create an authentication middleware function.
   - This middleware should:
     - Extract the Authorization header.
     - Validate the JWT.
     - Extract the user_id from the token.
     - Add the user_id to the request's context (context.WithValue) so the protected handlers can access it.
     - If auth fails, return a 401 Unauthorized error.
5. Todo Handlers:
   - Implement the 5 CRUD handlers for /todos.
   - Crucially: Each handler must get the user_id from the request's context (set by the middleware).
   - All database queries must include WHERE user_id = ? to ensure users can only access their own data.
6. Main Function:
   - In main.go, set up your chi router.
   - Define a public group for /register and /login.
   - Define a protected group for /todos that uses your authentication middleware.
   - Start the http.Server.
7. Demo

```bash
go run ./cmd/api
````
