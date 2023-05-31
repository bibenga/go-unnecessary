begin;

create table audit (
   id uuid primary key not null default gen_random_uuid(),
   create_ts timestamp with time zone not null default current_timestamp,
   "username" varchar not null,
   ip cidr,
   request_id varchar not null,
   message text not null,
   meta jsonb
);
create index audit_create_ts_key on audit using brin(create_ts);

commit;
