CREATE TABLE "accounts" (
                            "id" bigserial PRIMARY KEY,
                            "username" varchar NOT NULL,
                            "balance" bigint NOT NULL,
                            "currency" varchar NOT NULL,
                            "location" varchar NOT NULL,
                            "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "records" (
                           "id" bigserial PRIMARY KEY,
                           "account_id" bigint NOT NULL,
                           "amount" bigint NOT NULL,
                           "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
                                "id" bigserial PRIMARY KEY,
                                "from_account_id" bigint NOT NULL,
                                "to_account_id" bigint NOT NULL,
                                "amount" bigint NOT NULL,
                                "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "records" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

CREATE INDEX ON "accounts" ("username");

CREATE INDEX ON "records" ("account_id");

CREATE INDEX ON "transactions" ("from_account_id");

CREATE INDEX ON "transactions" ("to_account_id");

CREATE INDEX ON "transactions" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "records"."amount" IS 'Can be negative or positive';

COMMENT ON COLUMN "transactions"."amount" IS 'Must be positive';