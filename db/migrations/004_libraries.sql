drop table if exists libraries cascade;
create table libraries (
  id bigserial not null primary key,
  project_id bigint references projects(id) on delete cascade not null,
  name text not null,
  version text not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now()
);

create index libraries_project_id_idx on libraries using btree(project_id);
create unique index libraries_uniq_idx on libraries using btree(project_id, name, version);

drop table if exists error_libraries;
create table error_libraries (
  error_id bigint references errors(id) on delete cascade not null,
  library_id bigint references libraries(id) on delete cascade not null,
  primary key(error_id, library_id)
);

create index error_libraries_library_id_idx on error_libraries using btree(library_id);
