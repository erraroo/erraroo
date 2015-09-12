create table rate_limit_notifications (
  id bigserial not null primary key,
  account_id bigint references accounts(id) not null,
  created_at timestamp without time zone not null default now()
);
create index rate_limit_notifications_account_id on rate_limit_notifications using btree(account_id);
