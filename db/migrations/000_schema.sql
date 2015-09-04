create function now_utc() returns timestamp as $$
  select now() at time zone 'utc';
$$ language sql;

create table accounts (
  id bigserial not null primary key,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now()
);

create table users (
  id bigserial not null primary key,
  email text not null,
  encrypted_password text not null,
  account_id bigint references accounts(id) not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now()
);

create unique index users_email_unique on users using btree(email);
create index users_account_id on users using btree(account_id);

create table projects (
  id bigserial not null primary key,
  name text not null,
  token text not null,
  account_id bigint references accounts(id) on delete cascade not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now()
);

create unique index projects_token_unique on projects using btree(token);
create index project_account_ids on projects using btree(account_id);

create table events (
  id         bigserial not null primary key,
  checksum   text not null,
  kind       text not null,
  project_id bigint references projects(id) on delete cascade not null,
  created_at timestamp without time zone not null default now()
);

create index events_project_id on events (project_id);
create index events_checksum on events (checksum);

create table errors (
  id bigserial not null primary key,
  name text not null default 'Error',
  message text not null,
  checksum text not null,
  occurrences integer not null default 0,
  resolved boolean not null default false,
  muted boolean default false not null,
  last_seen_at timestamp without time zone not null default now(),
  project_id bigint references projects(id) on delete cascade not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now()
);

create unique index errors_project_id_checksum_idx on errors(project_id, checksum);
create index errors_project_id_idx on errors(project_id);
create index errors_checksum_idx on errors(checksum);

create table timings (
  id         bigserial not null primary key,
  created_at timestamp without time zone NOT NULL DEFAULT now(),
  project_id bigint references projects(id) on delete cascade not null,
  payload    json
);

create index timings_project_id_idx on timings(project_id);
create index timings_created_at_idx on timings(created_at);

create table prefs (
  user_id bigint references users(id) on delete cascade not null primary key,
  email_on_error boolean default false not null
);
