create table repositories (
  id bigserial not null primary key,
  project_id bigint references projects(id) not null,
  provider text not null,
  github_org text not null default '',
  github_repo text not null default '',
  github_scope text not null default '',
  github_token text not null default '',
  github_token_type text not null default ''
);
create unique index repositories_project_id_provider_uniq_idx on repositories using btree(project_id, provider);
