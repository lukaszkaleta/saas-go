CREATE TABLE job_chat (
                          id bigint PRIMARY KEY DEFAULT nextval('job_sequence'),
                          job_id bigint NOT NULL REFERENCES job(id),
                          worker_id bigint NOT NULL REFERENCES users(id),
                          action_created_by_id bigint not null references users,
                          action_created_at timestamp not null default now(),
                          UNIQUE(job_id, worker_id)
);
CREATE INDEX job_chat_job_idx ON job_chat(job_id);
CREATE INDEX job_chat_worker_idx ON job_chat(worker_id);


CREATE TABLE job_message (
                             id bigint NOT NULL PRIMARY KEY DEFAULT nextval('job_sequence'),
                             chat_id bigint NOT NULL REFERENCES job_chat(id) ON DELETE CASCADE,
                             value text NOT NULL, value_generated boolean NOT NULL DEFAULT false,
                             action_created_by_id bigint NOT NULL REFERENCES users(id),
                             action_created_at timestamp NOT NULL DEFAULT now()
);
CREATE INDEX job_message_chat_idx ON job_message(chat_id, action_created_at);
CREATE INDEX job_message_created_by_idx ON job_message(action_created_by_id);

CREATE TABLE job_chat_read (
                               chat_id bigint NOT NULL REFERENCES job_chat(id) ON DELETE CASCADE,
                               last_read_message_id bigint REFERENCES job_message(id),
                               action_updated_by_id bigint NOT NULL REFERENCES users(id),
                               action_updated_at timestamp NOT NULL DEFAULT now(),
                               PRIMARY KEY(chat_id, action_updated_by_id)
);
CREATE INDEX job_chat_read_user_idx ON job_chat_read(action_updated_by_id);

CREATE TABLE job_message_filesystem (
                                        job_message_id bigint not null references job_message,
                                        filesystem_id bigint not null references filestore_filesystem
);
CREATE UNIQUE INDEX job_message_filesystem_uidx ON job_message_filesystem USING btree (job_message_id, filesystem_id);
CREATE INDEX job_message_filesystem_job_idx ON job_message_filesystem USING btree (job_message_id);
