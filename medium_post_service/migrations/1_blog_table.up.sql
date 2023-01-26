CREATE TABLE IF NOT EXISTS "categories"(
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR(100) NOT NULL UNIQUE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "posts"(
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR NOT NULL,
    "description" TEXT NOT NULL,
    "image_url" VARCHAR,
    "user_id" INTEGER NOT NULL,
    "category_id" INTEGER NOT NULL REFERENCES categories(id),
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "views_count" INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS posts_title_idx ON posts(title);

CREATE TABLE IF NOT EXISTS "comments"(
    "id" SERIAL PRIMARY KEY,
    "post_id" int not null REFERENCES posts(id),
    "user_id" int not null,
    "description" text not null,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "likes"(
    "id" SERIAL PRIMARY KEY,
    "post_id" int not null REFERENCES posts(id),
    "user_id" int not null,
    "status" boolean not null,
    UNIQUE(post_id, user_id)
);