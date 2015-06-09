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
  account_id bigint references accounts(id) not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now()
);

create unique index projects_token_unique on projects using btree(token);
create index project_account_ids on projects using btree(account_id);

create table errors (
  id         bigserial not null primary key,
  payload    json,
  checksum   text not null,
  project_id bigint references projects(id) NOT NULL,
  created_at timestamp without time zone NOT NULL DEFAULT now(),
  updated_at timestamp without time zone NOT NULL DEFAULT now()
);

create index errors_project_id on errors (project_id);
create index errors_checksum on errors (checksum);

create table groups (
  id bigserial not null primary key,
  message text not null,
  checksum text not null,
  occurrences integer not null default 0,
  resolved boolean not null default false,
  last_seen_at timestamp without time zone not null default now(),

  project_id bigint references projects(id) not null,

  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now()
);

create unique index groups_project_id_checksum_idx on groups(project_id, checksum);
create index groups_project_id_idx on groups(project_id);
create index groups_checksum_idx on groups(checksum);

create table timings (
  id         bigserial not null primary key,
  created_at timestamp without time zone NOT NULL DEFAULT now(),
  project_id bigint references projects(id) NOT NULL,
  payload    json
);

create index timings_project_id_idx on timings(project_id);
create index timings_created_at_idx on timings(created_at);
