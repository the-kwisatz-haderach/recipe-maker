-- Write your migrate up statements here
CREATE extension IF NOT EXISTS "uuid-ossp";

CREATE TABLE
  recipe (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    recipe_name VARCHAR(255) NOT NULL
  );

CREATE TABLE
  recipe_user (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    recipe_id UUID NOT NULL,
    user_id UUID NOT NULL,
    relation VARCHAR(255) DEFAULT 'viewer' NOT NULL
  );

CREATE TABLE
  "user" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    username VARCHAR(255),
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
  );

ALTER TABLE recipe_user ADD FOREIGN KEY (user_id) REFERENCES "user" (id);

ALTER TABLE recipe_user ADD FOREIGN KEY (recipe_id) REFERENCES recipe (id);

---- create above / drop below ----
DROP TABLE IF EXISTS recipe CASCADE;

DROP TABLE IF EXISTS "user" CASCADE;

DROP TABLE IF EXISTS recipe_user CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.