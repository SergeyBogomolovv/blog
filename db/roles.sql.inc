-- Сервисные учётные записи с ограниченными правами.

CREATE ROLE admin
    LOGIN
    PASSWORD :'admin_password'
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE;

CREATE ROLE "user"
    LOGIN
    PASSWORD :'user_password'
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE;

REVOKE ALL
    ON DATABASE :"db_name"
    FROM PUBLIC;

GRANT CONNECT
    ON DATABASE :"db_name"
    TO admin, "user";

REVOKE ALL
    ON SCHEMA public
    FROM PUBLIC;

GRANT USAGE, CREATE
    ON SCHEMA public
    TO admin;

GRANT USAGE
    ON SCHEMA public
    TO "user";
