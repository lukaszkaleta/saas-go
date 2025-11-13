-- module tables
CREATE TABLE if not exists content (
  id serial primary key,
  name_value text not null default '',
  name_slug text not null default '',
  value text not null default '{}'
);

CREATE TABLE if not exists content_localization (
  id serial primary key,
  content_id bigint references content,
  country text not null default 'NO',
  language text not null default 'nb',
  translation_value text not null default '',
  translation_slug text not null default ''
);