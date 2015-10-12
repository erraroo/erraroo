create table outdated_revisions (
  id bigserial not null primary key,
  project_id bigint references projects(id) not null,
  sha text not null,
  dependencies json not null,
  created_at timestamp without time zone not null default now()
);
create unique index outdated_revisions_uniq_project_id_sha on outdated_revisions using btree(project_id, sha);
