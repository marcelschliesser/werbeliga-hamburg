# Hamburger Werbeliga Analysis

Analyse 2188 match results (24.01.2025) from werbeliga.de

## Example Queries with DuckDB
```sql
INSTALL sqlite;
LOAD sqlite;
ATTACH 'app.db' AS db (TYPE sqlite);
USE db;
SHOW TABLES;
SELECT * FROM db.matches ORDER BY home_score DESC LIMIT 5;
```
