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
