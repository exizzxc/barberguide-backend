CREATE TABLE haircuts (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    style       VARCHAR(50),
    length      VARCHAR(20),
    popularity  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at  TIMESTAMP WITH TIME ZONE
);

CREATE TABLE haircut_face_shapes (
    haircut_id  INTEGER NOT NULL REFERENCES haircuts(id) ON DELETE CASCADE,
    face_shape  VARCHAR(50) NOT NULL,
    PRIMARY KEY (haircut_id, face_shape)
);

CREATE TABLE haircut_hair_types (
    haircut_id  INTEGER NOT NULL REFERENCES haircuts(id) ON DELETE CASCADE,
    hair_type   VARCHAR(50) NOT NULL,
    PRIMARY KEY (haircut_id, hair_type)
);

CREATE TABLE haircut_images (
    id          SERIAL PRIMARY KEY,
    haircut_id  INTEGER NOT NULL REFERENCES haircuts(id) ON DELETE CASCADE,
    url         VARCHAR(500) NOT NULL,
    angle       VARCHAR(50),
    is_main     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_haircuts_deleted_at ON haircuts(deleted_at);
CREATE INDEX idx_haircuts_popularity ON haircuts(popularity DESC);
CREATE INDEX idx_haircut_images_haircut_id ON haircut_images(haircut_id);