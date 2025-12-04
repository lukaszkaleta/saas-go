create sequence if not exists job_sequence;

CREATE TABLE if not exists job (
  id bigint not null primary key default nextval('job_sequence'),
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
  rating int not null default 10,
  status_draft timestamp not null default now(),
  status_published timestamp,
  status_occupied timestamp,
  status_closed timestamp,
  action_created_by_id bigint not null references users,
  action_created_at timestamp not null default now()
 );

CREATE TABLE if not exists job_filesystem (
  job_id bigint not null references job,
  filesystem_id bigint not null references filestore_filesystem
);

CREATE TABLE if not exists job_offer (
  id bigint not null primary key default nextval('job_sequence'),
  job_id bigint not null references job,
  price_value int not null default 0,
  price_currency text not null default 'NOK',
  description_value text not null default '',
  description_image_url text not null default '',
  rating int not null default 0,
  action_created_by_id bigint not null references users,
  action_created_at timestamp not null default now(),
  action_accepted_by_id bigint references users,
  action_accepted_at timestamp,
  action_rejected_by_id bigint references users,
  action_rejected_at timestamp
);

CREATE TABLE messages (
  id bigint not null primary key default nextval('job_sequence'),
  job_id bigint not null references job,
  body        TEXT NOT NULL,
  action_created_by_id bigint not null references users(id),
  action_created_at timestamp not null default now(),
);