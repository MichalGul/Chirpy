## Chirpy
---
API that imitate Twitter like social App. Made in Golang

### Technology
---
Requires `postgresql` version 15 or above.
Uses `goose` for database migration and `sqlc` to handle querry code.



### Instalation
- Install postgres 15 or above. Install [goose](https://github.com/pressly/goose) and [sqlc](https://sqlc.dev/).

- Setup .env file to have fields below
    ```env
    DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable" #example database connection string
    PLATFORM="dev"
    SECRET="aYQNFfeUTsiZEbC4VAQWAlu0jLBYJe2XRv1sB3kOwDK5qb5IhBbrZvspSfCF8FYlUGtD0fQHby/Wt78vRBzPzw=="
    ```
- Run migrations from `schema` folder
    `goose up`
- `sqlc` config is already provided
- Build and run `go build -o chirpy` `./chirpy`

Test with example endpoinst

### Example endoints
---
- `POST /api/users`
    Adds user with email and password

    ```json
    {
    "email": "john@doe.com",
    "password": "123456"
    }
    ```
    Response body
    ```json
    {
        "id": "bc393a1c-86e6-463f-bb5d-81464bc250aa",
        "created_at": "2025-04-25T18:09:42.144158Z",
        "updated_at": "2025-04-25T18:09:42.144158Z",
        "email": "john@doe.com",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiYmMzOTNhMWMtODZlNi00NjNmLWJiNWQtODE0NjRiYzI1MGFhIiwiZXhwIjoxNzQ1NjAwOTg4LCJpYXQiOjE3NDU1OTczODh9.8oH2Vu0CQ-WCQvg--6yV44GHZRHYzIsYtc_CNElUARw",
        "refresh_token": "aba3b5042fa788044d25eccc9fb7b732e8a6862bfd628d8125325927d2978749"
    }
    ```
    TBD

