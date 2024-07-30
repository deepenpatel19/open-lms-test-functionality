# open-lms-test-functionality

This project will user Go version  1.22

### Setup
Run `go install` in project directory.

#### Migration Notes
- `migrate` library for managing migrations.
- Command to create migration : `migrate create -ext sql -dir migrations -seq -digits 6 <migration_name>`. This command will generate migrations in `migrations` directory.


#### Library reference
- JWT: `https://github.com/appleboy/gin-jwt/`
- Migrate: `https://github.com/golang-migrate/migrate`
- Gin: `https://github.com/gin-gonic/gin`