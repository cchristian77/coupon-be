## Coupon API

Coupon API is a robust backend service designed to facilitate coupon management and claim process.
This system provides RESTful API endpoints which enables efficient functionalities,
utilizing Go and PostgreSQL for reliable data storage.
The study case is to **solve race condition when concurrent users try to claim a single coupon with
a limited quota**.

### Technology
1. Backend: Go v1.24.10
2. Database : Postgres 14.7
3. Redis : latest

### Prerequisites 
1. Docker 

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
- `sqlmock`: Mock SQL for database testing. **(Mandatory to Install)**
- `mock`: Mocking framework. **(Mandatory to Install)**
- `goose`: database migration library. **(Mandatory to Install)**

### Installation
Before running the application, you need to setup the necessary prerequisites, as following :
1. Clone the repository
   ```bash
   git clone git@github.com:cchristian77/coupon-be.git
   ```

2. Configure environment variable </br>
   Use **'localhost'** instead on **database.host** and **redis.host**, if backend is not run on Docker. </br>
   The port of application is set to **9000** as default (or setup based on your preferred configuration).
     ```bash
    copy env.json.example 
    ```

2. Initialize services (database, redis, backend app)
    ```bash
    docker compose up -d
    ```

3. Install dependencies
    ```bash
    go mod download
    ```

4. Run database migrations
    ```bash
     goose postgres "user=admin password=password dbname=coupon_db sslmode=disable" up
    ```

5. Configure environment variable
    ```bash
    copy .env.json.example and setup based on your preferred configuration
    ```

6. Migrate the database
   (Alternatively, you can use coupon_db.sql in the documentation directory.)
   ```bash
   goose -dir ./migrations  postgres "user=admin password=password dbname=coupon_db sslmode=disable" up
   ```

7. Run application
    ```
    docker compose up -d
    ```

### How to Test

1. Check healthcheck endpoint
    ```bash
    curl http://localhost:9000/healthcheck
    ```

2. Populate user database
   Populate user data
   ```bash
    curl http://localhost:9000/api/users/register
    ```
   Populate coupon data
   ```bash
   curl --request POST \
      --url http://localhost:9000/api/coupons \
      --header 'content-type: application/json' \
      --data '{
        "name": "COUPON_TEST",
        "amount": 10
      }'
   ```

3. Run stress test
   ```
   go run ./stress_test/
    ```

### Architecture Notes
![Database Schema.png](documentation/Database%20Schema.png)

- `users` table contains field for a single user such as username, password, etc.
- `coupons` table contains `name`, `amount`, and `remaining_amount` for a single coupon.
- `user_claims` table is a pivot table between coupons and users table. Since multiple users can claim multiple coupons, with unique constraint for
  `user_id` and `coupon_id` pairs.

#### Locking Strategy
In a concurrent environment, multiple users may attempt to claim the same coupon simultaneously,
which is introduced a race condition issue. </br>

To prevent this, I serialize access to the claim process using a distributed lock using Redis.
When a user attempts to claim a coupon, the system tries to set a unique key `claim:coupon:{coupon_name}` in Redis;
if this key is successfully set, the process "wins" the lock and proceeds the claim process by inserting an entry in
user_claims table and decrement the `coupon.remaining_amount `by 1. Concurrent requests will detect that the key already
exists, forcing them to pause for `500ms` and retry up to `3` times. </br>

This approach eliminates the race conditions in the coupon claim process by enforcing that only 1 claim is processed,
while maintaining the high performance of the API. Please check the documentation directory for the detailed solution. </br>

### Author
Chris Christian 
