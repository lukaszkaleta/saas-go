CREATE TABLE if not exists filestore_record (
                                                id serial primary key,
                                                name_value text not null default '',
                                                name_slug text not null default '',
                                                description_value text not null default '',
                                                description_image_url text not null default ''
);

CREATE TABLE if not exists filestore_container (
                                                   id serial primary key,
                                                   name_value text not null default '',
                                                   name_slug text not null default '',
                                                   description_value text not null default '',
                                                   description_image_url text not null default ''
);

CREATE TABLE if not exists filestore_filesystem (
                                                    id serial primary key,
                                                    name_value text not null default '',
                                                    name_slug text not null default ''
);

CREATE TABLE if not exists filesystem_record (
                                                 filesystem_id int not null references filestore_filesystem,
                                                 record_id int not null references filestore_record
);

CREATE TABLE if not exists filesystem_container (
                                                    filesystem_id int not null references filestore_filesystem,
                                                    container_id int not null references filestore_container
);

CREATE TABLE if not exists filesystem_container_record (
                                                           container_id int not null references filestore_container,
                                                           record_id int not null references filestore_record
);
