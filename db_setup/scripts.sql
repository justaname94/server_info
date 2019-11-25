CREATE USER
IF NOT EXISTS go_admin;
CREATE DATABASE
IF NOT EXISTS website;

GRANT ALL ON DATABASE website TO go_admin;

CREATE TABLE
IF NOT EXISTS website.site
(
  domain STRING PRIMARY KEY,
  title STRING,
  created_at TIMESTAMP,
  ssl_grade STRING,
  previous_ssl_grade STRING,
  logo STRING,
  is_down Bool
);

CREATE TABLE
IF NOT EXISTS website.server
(
  address STRING,
  ssl_grade STRING,
  country STRING,
  owner STRING
)