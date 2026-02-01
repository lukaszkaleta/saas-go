CREATE TABLE test_filesystem (
  test_id bigint not null,
  filesystem_id bigint not null references filestore_filesystem
);