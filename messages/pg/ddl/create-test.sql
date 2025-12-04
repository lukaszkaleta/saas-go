create sequence if not exists user_sequence;
CREATE TABLE if not exists users (
  id bigint not null primary key default nextval('user_sequence')
);

-- test messages are related to user.
CREATE TABLE message (
  id bigint not null primary key default nextval('user_sequence'),
  user_id bigint not null references users,
  body        TEXT NOT NULL,
  action_created_by_id bigint not null references users(id),
  action_created_at timestamp not null default now(),
);
