create table levels (
  id integer primary key,
  value varchar(10) not null,
  created_at timestamp default CURRENT_TIMESTAMP not null
);
