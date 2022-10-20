CREATE TABLE skillbox.users (
id SERIAL PRIMARY KEY NOT NULL,
name VARCHAR(50) NOT NULL,
age NUMERIC(3),
friends text[]
);