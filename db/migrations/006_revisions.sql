create table revisions (
  id bigserial not null primary key,
  project_id bigint references projects(id) not null,
  sha text not null,
  dependencies json not null,
  created_at timestamp without time zone not null default now_utc(),
  updated_at timestamp without time zone not null default now_utc()
);
create unique index revisions_uniq_project_id_sha on revisions using btree(project_id, sha);
