CREATE USER IF NOT EXISTS go_admin;
CREATE DATABASE IF NOT EXISTS website;
GRANT ALL ON DATABASE website TO go_admin;
CREATE TABLE IF NOT EXISTS website.site (
  domain STRING PRIMARY KEY,
  title STRING,
  ssl_grade STRING,
  previous_ssl_grade STRING,
  logo STRING,
  is_down Bool,
  servers_changed bool,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
CREATE TABLE IF NOT EXISTS website.server (
  id SERIAL PRIMARY KEY,
  address STRING,
  ssl_grade STRING,
  country STRING,
  owner STRING,
  site string NOT NULL REFERENCES website.site(domain)
)