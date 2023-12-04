-- Write your migrate up statements here
CREATE TABLE
  ingredient (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(255) NOT NULL
  );

CREATE TABLE
  pantry_item (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    ingredient_id UUID NOT NULL,
    quantity INT DEFAULT 1,
    unit VARCHAR(255) DEFAULT '' NOT NULL
  );

CREATE TABLE
  pantry_item_user (
    user_id UUID NOT NULL,
    pantry_item_id UUID NOT NULL,
    PRIMARY KEY (user_id, pantry_item_id)
  );

CREATE TABLE
  ingredient_user (
    user_id UUID NOT NULL,
    ingredient_id UUID NOT NULL,
    PRIMARY KEY (user_id, ingredient_id)
  );

ALTER TABLE pantry_item ADD FOREIGN KEY (ingredient_id) REFERENCES ingredient (id);

ALTER TABLE pantry_item_user ADD FOREIGN KEY (user_id) REFERENCES "user" (id);

ALTER TABLE pantry_item_user ADD FOREIGN KEY (pantry_item_id) REFERENCES pantry_item (id);

ALTER TABLE ingredient_user ADD FOREIGN KEY (user_id) REFERENCES "user" (id);

ALTER TABLE ingredient_user ADD FOREIGN KEY (ingredient_id) REFERENCES ingredient (id);

---- create above / drop below ----
DROP TABLE IF EXISTS ingredient CASCADE;

DROP TABLE IF EXISTS pantry_item CASCADE;

DROP TABLE IF EXISTS pantry_item_user CASCADE;

DROP TABLE IF EXISTS ingredient_user CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.