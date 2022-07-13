# Working with Postgres

Examples of how to use common commands. 

## Backups and Restoration
q
Create a backup.
```
pg_dump -d postgres://postgres:postgres@localhost:30010/user > ./data/user_20220704.sql
```
Drop a database.
```
dropdb -h localhost -p 30010 -U postgres -W user
```
Recreate database.
```
createdb -h localhost -p 30010 -U postgres -W user
```
Restore database.
```
psql -d postgres://postgres:postgres@localhost:30010/user < ./data/user_20220704.sql
```
Close connections to database.
```
SELECT pg_terminate_backend(pid)FROM pg_stat_activity WHERE datname = 'user'
```

## Working with Users

### Create non-root user

When testing Row Level Security the root user will always see results. Create a non-root user for the application.
Then grant all privileges to all tables.
```
CREATE USER user_a WITH PASSWORD 'postgres';
GRANT ALL ON ALL TABLES IN SCHEMA "public" TO user_a;
```
### List users

`\du`

### Switching users after connecting
Supply the database and user name.

`\c user user_a`

### Dropping a user
You can't drop a user while the user still owns something or has any granted privileges on other objects.
```
REASSIGN OWNED BY user_a TO postgres;
DROP OWNED BY user_a;
DROP USER user_a;
```