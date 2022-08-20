DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id SERIAL NOT NULL,
    line_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS weight_histories;
CREATE TABLE weight_histories (
    id SERIAL NOT NULL,
    user_id SERIAL NOT NULL,
    weight_num NUMERIC(5, 1) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);