create table prefs (
  user_id bigint references accounts(id) not null primary key,
  email_on_error boolean default true not null
);
