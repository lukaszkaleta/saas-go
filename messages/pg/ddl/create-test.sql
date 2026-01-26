create sequence if not exists user_sequence;
CREATE TABLE if not exists users (
  id bigint not null primary key default nextval('user_sequence')
);

-- test messages are related to user.
-- those are direct messages (only one user is recipient)
CREATE TABLE message (
  id bigint not null primary key default nextval('user_sequence'),
  recipient_id bigint not null references users,
  owner_id bigint not null references users,
  value TEXT NOT NULL,
  action_created_by_id bigint not null references users(id),
  action_created_at timestamp not null default now(),
  action_read_by_id bigint references users(id),
  action_read_at timestamp
);
