create sequence if not exists user_sequence;
CREATE TABLE if not exists users (
  id bigint not null primary key default nextval('user_sequence')
);