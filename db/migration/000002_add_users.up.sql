CREATE TABLE "users" (
                         "username" varchar PRIMARY KEY,
                         "hashed_password" varchar NOT NULL,
                         "first_name" varchar NOT NULL,
                         "last_name" varchar NOT NULL,
                         "email" varchar UNIQUE NOT NULL,
                         "contact_number" varchar DEFAULT (now()),
                         "address" varchar DEFAULT (now()),
                         "updated_at" timestamptz NOT NULL DEFAULT (now()),
                         "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

CREATE UNIQUE INDEX ON "accounts" ("username", "currency");