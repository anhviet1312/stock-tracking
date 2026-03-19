-- +goose Up
CREATE TABLE stories
(
    "id"            UUID PRIMARY KEY,
    "slug"          TEXT NOT NULL,
    "name"          TEXT NOT NULL,
    "description"   TEXT NOT NULL,
    "type"          TEXT NOT NULL,
    "status"        TEXT NOT NULL,
    "image_url"     TEXT,
    "created_by"    UUID,
    "updated_at"    TIMESTAMP WITH TIME ZONE NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT fk__stories__users FOREIGN KEY ("created_by") REFERENCES "users" ("id")
);

CREATE INDEX idx_stories__slug ON "stories" ("slug");

CREATE TABLE chapters
(
    "id"            UUID PRIMARY KEY,
    "name"          TEXT NOT NULL,
    "story_id"      UUID NOT NULL,
    "status"        TEXT NOT NULL,
    "number"        INTEGER NOT NULL,
    "description"   TEXT,
    "content"       TEXT,
    "metadata"      jsonb NOT NULL DEFAULT '{}'::jsonb,
    "created_by"    UUID,
    "updated_at"    TIMESTAMP WITH TIME ZONE NOT NULL,
    "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT fk__chapters__users FOREIGN KEY ("created_by") REFERENCES "users" ("id")
);

CREATE INDEX idx_chapters__story_id_number_status ON "chapters" ("story_id", "number", "status");

-- +goose Down
DROP TABLE "chapters";
DROP TABLE "stories";