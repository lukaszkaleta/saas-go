CREATE TABLE if not exists users (
    id serial primary key,
    account_token text not null default '',
    person_first_name text not null default '',
    person_last_name text not null default '',
    person_email text not null default '',
    person_phone text not null default '',
    avatar_description_value text not null default '',
    avatar_description_image text not null default '',
    address_line1 text not null default '',
    address_line2 text not null default '',
    address_city       text not null default '',
    address_postal_code text not null default '',
    address_district   text not null default '',
    settings_radar_perimeter int not null default 10000,
    settings_radar_position_latitude int not null default 0,
    settings_radar_position_longitude int not null default 0
);
