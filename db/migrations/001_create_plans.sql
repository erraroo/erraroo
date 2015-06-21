create table plans (
  account_id bigint references accounts(id) not null primary key,
  data_retention_in_days int null null default 90,
  requests_per_minute int not null default 200
);
