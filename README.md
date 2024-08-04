# open-lms-test-functionality

This project will user Go version  1.22

### Setup
Run `go install` in project directory.

#### Migration Notes
- `migrate` library for managing migrations.
- Command to create migration : `migrate create -ext sql -dir migrations -seq -digits 6 <migration_name>`. This command will generate migrations in `migrations` directory.

#### Statistics
##### Resource used for this testing
- Cloud Provider: Digital Ocean
- Server: 1 vCPU / 1 Gb RAM 
- No of server: 2
- Database : Postgres 16 - 1 vCPU / 1 Gb RAM

1) Create question: True or False type
![Alt text](./docs/images/create_questions_tf.png)
![Alt text](./docs/images/create_question_tf_server_cpu_memory.png)

2) Create question: Multiple choice type
![Alt text](./docs/images/create_questions_mc_2_iteration.png)
![Alt text](./docs/images/create_question_mc_2_iteration_server_cpu_memory.png)

3) Get question list with pagination of 10 questions
![Alt text](./docs/images/get_question_list_k6.png)
![Alt text](./docs/images/get_question_list.png)

4) Get question list with pagination of 50 questions
![Alt text](./docs/images/get_question_list_50_k6.png)
![Alt text](./docs/images/get_question_50_list_cpu_memory.png)

5) Submit question answer
![Alt text](./docs/images/submit_question_k6.png)
![Alt text](./docs/images/submit_question_server_cpu_memory.png)

#### Library reference
- JWT: `https://github.com/appleboy/gin-jwt/`
- Migrate: `https://github.com/golang-migrate/migrate`
- Gin: `https://github.com/gin-gonic/gin`


#### Tools for testing concurrent requests
- K6: `https://k6.io`