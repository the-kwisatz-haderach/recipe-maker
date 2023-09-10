-- Write your migrate up statements here
create table recipes (
  recipe_id serial primary key,
  recipe_name varchar(255)
);

---- create above / drop below ----

drop table recipes;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
