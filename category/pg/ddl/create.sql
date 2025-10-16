
CREATE TABLE if not exists category (
  id serial primary key,
  parent_category_id bigint references category,
  name_value text not null default '',
  name_slug text not null default '',
  description_value text not null default '',
  description_image_url text not null default ''
);

CREATE TABLE if not exists category_localization (
                                                     id serial primary key,
                                                     category_id bigint references category,
                                                     country text not null default 'NO',
                                                     language text not null default 'nb',
                                                     translation_value text not null default '',
                                                     translation_slug text not null default ''
);