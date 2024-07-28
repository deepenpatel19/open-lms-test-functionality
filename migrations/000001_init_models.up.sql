BEGIN;

CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL DEFAULT '',
    last_name VARCHAR(50) NOT NULL DEFAULT '',
    email VARCHAR(100) UNIQUE,
    password TEXT NOT NULL DEFAULT '',
    type INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS tests(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS questions(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    type INTEGER,
    question_data JSONB,
    answer_data JSONB
);

CREATE TABLE IF NOT EXISTS test_questions(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    test_id BIGSERIAL,
    question_id BIGSERIAL,

    CONSTRAINT test_id
        FOREIGN KEY(test_id)
            REFERENCES tests(id) ON DELETE CASCADE,

    CONSTRAINT question_id
        FOREIGN KEY(question_id)
            REFERENCES questions(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS test_question_submissions(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id BIGSERIAL,
    test_id BIGSERIAL,
    question_id BIGSERIAL,
    submitted_data JSONB,
    answer_status BOOLEAN DEFAULT false,

    CONSTRAINT test_id
        FOREIGN KEY(test_id)
            REFERENCES tests(id) ON DELETE CASCADE,

    CONSTRAINT question_id
        FOREIGN KEY(question_id)
            REFERENCES questions(id) ON DELETE CASCADE,

    CONSTRAINT user_id
        FOREIGN KEY(user_id)
            REFERENCES users(id) ON DELETE CASCADE

);

COMMIT;