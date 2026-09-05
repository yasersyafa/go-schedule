CREATE TABLE "notif_logs" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "activity_id" UUID NOT NULL,
    "sent_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_notification_activity"
        FOREIGN KEY ("activity_id") REFERENCES "activities"("id") ON DELETE CASCADE
);