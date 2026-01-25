create table if not exists users (
    id serial primary key,
    account_token text not null default '',
    person_first_name text not null default '',
    person_last_name text not null default '',
    person_email text not null default '',
    person_phone text not null default '',
    address_line_1 text not null default '',
    address_line_2 text not null default '',
    address_city       text not null default '',
    address_postal_code text not null default '',
    address_district   text not null default '',
    avatar_description_value text not null default '',
    avatar_description_image_url text not null default '',
    settings_radar_perimeter int not null default 10000,
    settings_radar_position_latitude double precision not null default 0,
    settings_radar_position_longitude double precision not null default 0
);

create table if not exists user_filesystem (
    id serial primary key,
    name text not null default '',
    user_id bigint not null references users,
    filesystem_id bigint not null references filestore_filesystem
);

create table if not exists user_rating (
  id serial primary key,
  user_id bigint references users
);

