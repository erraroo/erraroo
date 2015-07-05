create table password_recovers (
  token text not null primary key,
  used boolean not null default false,
  user_id bigint references users(id) not null,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now()
);
