CREATE TABLE "activities" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" VARCHAR(255) NOT NULL,
    "notes" TEXT NULL,
    "day" VARCHAR(255) NOT NULL CHECK (
        "day" IN ('monday','tuesday','wednesday','thursday','friday','saturday','sunday')
    ),
    "start_time" TIME NOT NULL,
    "end_time" TIME NOT NULL,
    "last_notified_date" DATE NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT "chk_time_order" CHECK ("start_time" < "end_time"),
    PRIMARY KEY ("id")
);