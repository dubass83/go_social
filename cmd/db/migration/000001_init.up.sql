CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "email" citext UNIQUE NOT NULL,
  "password" bytea NOT NULL,
  "created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE "posts" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "content" text NOT NULL,
  "user_id" integer NOT NULL,
  "tags" varchar(100) [],
  "created_at" timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  "updated_at" timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

COMMENT ON COLUMN "posts"."content" IS 'Body of the post';

ALTER TABLE "posts" ADD CONSTRAINT "user_posts" FOREIGN KEY ("user_id") REFERENCES "users" ("id");
