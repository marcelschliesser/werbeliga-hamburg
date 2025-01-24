# Hamburger Werbeliga Analysis

Analyse 2188 match results (24.01.2025) from werbeliga.de

## Example Queries with DuckDB
```sql
SELECT * FROM read_json('data.json') ORDER BY homeScore DESC LIMIT 5;
```
```sql
SELECT * FROM read_json('data.json') ORDER BY awayScore DESC LIMIT 5;
```
```sql
SELECT COUNT(*) FROM read_json('data.json');
```
