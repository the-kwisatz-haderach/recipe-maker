-- Write your migrate up statements here
create table recipes (
  id integer primary key generated always as identity,
  recipe_name varchar(255)
);

create table users (
  id integer primary key generated always as identity,
  username varchar(255),
  password varchar(255) not null,
  email varchar(255) not null

)

---- create above / drop below ----

drop table recipes;
drop table users;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
