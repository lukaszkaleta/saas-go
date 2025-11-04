CREATE TABLE if not exists job (
  id serial primary key,
  description_value text not null default '',
  description_image_url text not null default '',
  address_line_1 text not null default '',
  address_line_2 text not null default '',
  address_city       text not null default '',
  address_postal_code text not null default '',
  address_district   text not null default '',
  position_latitude numeric not null default 0,
  position_longitude numeric not null default 0,
  price_value int not null default 0,
  price_currency text not null default 'NOK',
  status_draft timestamp not null default now(),
  status_published timestamp,
  status_occupied timestamp,
  status_closed timestamp
);

CREATE TABLE if not exists job_filesystem (
  job_id int not null references job,
  filesystem_id int not null references filestore_filesystem
);