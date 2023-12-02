-- Write your migrate up statements here
create extension if not exists "uuid-ossp";

create table
  recipes (
    id uuid primary key default uuid_generate_v4 (),
    recipe_name varchar(255) not null
  );

create table
  recipe_roles (
    id uuid primary key default uuid_generate_v4 (),
    recipe_id uuid not null,
    user_id uuid not null,
    relation varchar(255) default 'viewer' not null
  );

create table
  users (
    id uuid primary key default uuid_generate_v4 (),
    username varchar(255),
    password varchar(255) not null,
    email varchar(255) not null unique
  );

alter table recipe_roles add foreign key (user_id) references users (id);

alter table recipe_roles add foreign key (recipe_id) references recipes (id);

---- create above / drop below ----
drop table recipes;

drop table users;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.