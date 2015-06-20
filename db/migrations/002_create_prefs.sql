create table prefs (
  user_id bigint references users(id) not null primary key,
  email_on_error boolean default true not null
);
