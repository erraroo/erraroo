create table invitations (
  token text not null primary key,
  user_id bigint references users(id) not null,
  account_id bigint references accounts(id) not null,
  address text not null,
  accepted boolean not null default false,
  created_at timestamp without time zone not null default now(),
  updated_at timestamp without time zone not null default now()
);

create index invitations_account_id on invitations using btree(account_id);
