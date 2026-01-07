## Base API

Base API is a robust backend service designed to streamline post and comment management.
This system provides RESTful API endpoints which enables efficient functionalities,
utilizing Go and PostgreSQL for reliable data storage. The purpose of this application is for my personal learning purposes.

### Technology
1. Backend: Go v1.24.10
2. Database : Postgres 14.7
3. Docker (optional)

### Project Structure
- `cmd/`: This directory contains application-specific entrypoints. It's the main program of the application.

- `domain/`: This directory holds struct object representation for each table in the database.

- `domain/enums`: This directory holds enums constants for each domain.

- `entrypoint/`: This directory is the controller layer for API endpoints. The controller layer's function is to accept and validate incoming requests before they are processed by the service layer.

- `migration/`: This directory holds SQL migrations files to create and modify tables in the database.

- `repository/`: This directory is a repository layer to handle interactions between the application and database.

- `request/`: This directory holds HTTP request structures that define the input of the application.

- `response/`: This directory contains HTTP response structures that define the output formats of the API.

- `service/`: This directory implements the core business logic layer, processing data between the controller and repository layers.

- `shared/`: This directory contains external and internal services which the application interacts with.

- `util/`: This directory holds utility functions and helper code that supports the main application functionality.

#### Library
- `gorilla/mux`: The HTTP router for web servers
- `gorm`: The ORM for database operations
- `koanf`: The config environment library
- `jwt`: JSON Web Token for authentication
- `validator`: validation library
- `uuid`: Generate and handle UUID
- `pgx`: PostgreSQL driver
- `decimal`: Handling decimal numbers with precision
- `zap`: Structured logging library
- `crypto`: Cryptography functions for password hashing
- `testify`: Unit Testing library
- `sqlmock`: Mock SQL for database testing
- `mock`: Mocking framework
- `goose`: database migration library

### Installation
Before running the application, you need to setup the necessary prerequisites, as following :
1. Clone the repository
   ```bash
   git clone git@github.com:cchristian77/payroll_be.git
   ```

2. Initialize the database
    ```bash
    docker compose up -d
    ```

3. Install dependencies
    ```bash
    go mod download
    ```

4. Run database migrations
    ```bash
     goose postgres "user=admin password=password dbname=payroll sslmode=disable" up
    ```

5. Configure environment variable
    ```bash
    copy .env.json.example and setup based on your preferred configuration
    ```

6. Run application
    ```
    go run ./cmd/web
    ```

7. Check healthcheck endpoint
    ```bash
    curl http://localhost:9000/healthcheck
    ```
8. Populate database
   ```bash
    curl http://localhost:9000/auth/v1/register
    ```
   
### Goose Command  
   ```bash
   goose -dir ./migrations  postgres "user=admin password=password dbname=base sslmode=disable" create migration_file_name sql 
   goose -dir ./migrations  postgres "user=admin password=password dbname=base sslmode=disable" up 
   goose -dir ./migrations  postgres "user=admin password=password dbname=base sslmode=disable" status
   ```

### Mockgen Command
   ```bash 
   go generate ./...
   ```

### Author
Chris Christian 
