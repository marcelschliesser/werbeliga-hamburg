# Hamburger Werbeliga Analysis

Analyse data from werbeliga.de

## Example Queries with DuckDB
```sql
SELECT * FROM read_json('data.json') ORDER BY homeScore DESC LIMIT 5;
```
```sql
SELECT * FROM read_json('data.json') ORDER BY awayScore DESC LIMIT 5;
```
