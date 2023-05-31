begin;

create table tag (
   id bigserial primary key,
   name varchar not null,
   created_ts timestamp with time zone not null default current_timestamp,
   modified_ts timestamp with time zone not null default current_timestamp,
   deleted_ts timestamp with time zone
);
create unique index tag_name_key on tag(lower(name));

create table dict_platform (
   id smallint primary key,
   name varchar not null,
   display_name varchar,
   created_ts timestamp with time zone not null default current_timestamp,
   modified_ts timestamp with time zone not null default current_timestamp,
   deleted_ts timestamp with time zone
);
create unique index dict_platform_name_key on dict_platform(lower(name));

create table application (
   id bigserial primary key,
   dict_platform_id smallint not null references dict_platform(id),
   name varchar not null,
   display_name varchar,
   created_ts timestamp with time zone not null default current_timestamp,
   modified_ts timestamp with time zone not null default current_timestamp,
   deleted_ts timestamp with time zone,
   some_flg boolean not null default false,
   some_double1 double precision not null default 0,
   some_double2 double precision,
   some_json jsonb not null default '{}'::jsonb
);
create unique index application_dict_platform_id_name_key on application(dict_platform_id, lower(name));

create table application_tag (
   id bigserial primary key,
   application_id bigint not null references application(id) on delete cascade,
   tag_id bigint not null references tag(id) on delete cascade,
   created_ts timestamp with time zone not null default current_timestamp
   -- primary key (application_id, tag_id)
);
create unique index application_tag_application_id_dict_platform_id_key on application_tag(application_id, tag_id);

commit;
