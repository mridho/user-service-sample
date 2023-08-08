/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- users table
CREATE TABLE users (
    "id" UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    "created_at" timestamp NOT NULL DEFAULT timezone('utc', now()),
    "updated_at" timestamp NOT NULL DEFAULT timezone('utc', now()),
    "deleted_at" timestamp,
    "phone_number" VARCHAR(20) NOT NULL UNIQUE,
    "full_name" VARCHAR(100) NOT NULL,
    "password_hash" VARCHAR(100) NOT NULL,
    "salt" VARCHAR(20) NOT NULL,
    "login_count" INTEGER NOT NULL DEFAULT 0
);

-- sample data, with password: pAssW0$ds
INSERT INTO users ("phone_number", "full_name", "password_hash", "salt") 
VALUES 
  ('+62810000001', 'Sample User 1', '9996f6bb66439b2d8bae91fc8f0fd81158c9d4f91ba9a892d30e2581ec8ddb26', '486j+Is1QGia1g=='), 
  ('+62810000002', 'Sample User 2', '9996f6bb66439b2d8bae91fc8f0fd81158c9d4f91ba9a892d30e2581ec8ddb26', '486j+Is1QGia1g=='),
  ('+62810000003', 'Sample User 3', '9996f6bb66439b2d8bae91fc8f0fd81158c9d4f91ba9a892d30e2581ec8ddb26', '486j+Is1QGia1g==');
