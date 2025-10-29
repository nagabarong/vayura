## Vayura Backend (Go + Gin + GORM)

A simple, production-ready REST API backend built with Go, Gin, and GORM. It provides user authentication with JWT, user profile management, and avatar upload handling.

### Features
- **Authentication**: Register and login with email/password, JWT issuance.
- **User Profile**: Fetch, update, and delete authenticated user profiles.
- **Avatar Upload**: Upload profile avatars with validation (size and type) saved to local storage.
- **Health Check**: Basic `/health` endpoint.

### Tech Stack
- **Language**: Go
- **Web Framework**: Gin
- **ORM**: GORM (PostgreSQL)
- **Auth**: JWT (HS256)

### Project Structure
```text
backend/
  cmd/server/main.go           # App entrypoint
  config/                      # Config and DB setup
  internal/
    handler/                   # HTTP handlers (auth, user)
    models/                    # GORM models
    repository/                # Data access layer
    service/                   # Business logic (auth, user, storage)
  pkg/                         # Shared utilities (jwt, middleware, responses, errors)
  routes/routes.go             # Route definitions
  Uploads/avatars/             # Uploaded avatar files
```

---

### Prerequisites
- Go 1.21+
- PostgreSQL 13+

### Environment Variables
Create a `.env` file in `backend/` directory based on `backend/env.example`:

```env
# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=vayura
DB_PORT=5432

# Application Configuration
APP_PORT=8080

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Storage
UPLOAD_DIR=Uploads/avatars
```

Notes:
- `UPLOAD_DIR` defaults to `Uploads/avatars` if not set.
- Ensure the PostgreSQL database (`DB_NAME`) exists and credentials are valid.

---

### Installation & Run
From the `backend/` directory:

```bash
# 1) Copy environment file
cp env.example .env

# 2) Run the server (Go)
go run ./cmd/server

# Or build binary
go build -o bin/server ./cmd/server
./bin/server
```

The server starts on `http://localhost:8080` (configurable via `APP_PORT`).

Database migrations run automatically for `models.User` at startup.

---

### API Overview
Base URL: `http://localhost:8080`

- `GET /health` — Health check
- `POST /api/auth/register` — Register
- `POST /api/auth/login` — Login, returns JWT
- `GET /api/user/profile` — Get own profile (auth)
- `PUT /api/user/profile` — Update own profile (auth)
- `DELETE /api/user/profile` — Delete own profile (auth)
- `POST /api/user/avatar` — Upload avatar (auth, multipart)

Authentication: Send `Authorization: Bearer <token>` header for protected endpoints.

Standard response shape:

```json
{
  "success": true,
  "message": "string",
  "data": {},
  "error": "string"
}
```

---

### Auth Endpoints

#### Register
`POST /api/auth/register`

Request JSON:
```json
{
  "full_name": "John Doe",
  "username": "johnd",
  "email": "john@example.com",
  "password": "secretPass1",
  "phone": "",
  "role": "user",
  "gender": "male",
  "birthday": "1990-01-01"
}
```

Responses:
- 201: user created
- 400: validation error or duplicate email/username

#### Login
`POST /api/auth/login`

Request JSON:
```json
{
  "email": "john@example.com",
  "password": "secretPass1"
}
```

Response 200:
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "<jwt>",
    "user": {
      "id": 1,
      "full_name": "John Doe",
      "username": "johnd",
      "email": "john@example.com",
      "phone": "",
      "role": "user",
      "gender": "male",
      "birthday": "1990-01-01T00:00:00Z"
    }
  }
}
```

---

### User Endpoints (Protected)

Send header: `Authorization: Bearer <token>`

#### Get Profile
`GET /api/user/profile`

Response 200: user object

#### Update Profile
`PUT /api/user/profile`

Request JSON (any fields optional):
```json
{
  "full_name": "Johnathan Doe",
  "username": "johnny",
  "phone": "+6212345678",
  "gender": "male",
  "birthday": "1991-02-03"
}
```

Responses:
- 200: updated user
- 400: username taken or invalid birthday format

#### Delete Profile
`DELETE /api/user/profile`

Response 200: confirmation

#### Upload Avatar
`POST /api/user/avatar` (multipart form)

Form fields:
- `avatar`: file (.jpg, .jpeg, .png), max 2 MB

Response 200: updated user with `avatar` path

Notes:
- Files are saved under `Uploads/avatars/` with pattern `<userId>_<timestamp>.<ext>`.
- Response `avatar` includes a leading slash (e.g. `/Uploads/avatars/1_1700000000.jpg`).

---

### Development Tips
- Switch GORM logger level in `config/config.go` if you need SQL logs.
- Ensure `.env` is in `backend/` as `godotenv.Load()` looks there.
- JWT expiration is 72 hours; keep `JWT_SECRET` strong in production.

### License
MIT


