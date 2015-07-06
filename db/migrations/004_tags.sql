drop table if exists project_tags;

create table project_tags (
  id  bigserial not null primary key,
  key text not null,
  project_id bigint references projects(id) not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now(),
  occurrences integer not null default 1
);

create index project_tags_keys on project_tags using btree(key);
create index project_tags_project_id on project_tags using btree(project_id);
create unique index project_unique_keys_project_id on project_tags using btree(key, project_id);

drop table if exists error_tag_values;
create table error_tag_values (
  id  bigserial not null primary key,
  key text not null,
  value text not null,
  project_id bigint references projects(id) not null,
  error_id bigint references errors(id) not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now(),
  occurrences integer not null default 1
);

create unique index error_tag_values_unique on error_tag_values using btree(key, project_id, error_id, value);
