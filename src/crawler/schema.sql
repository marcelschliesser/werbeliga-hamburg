CREATE TABLE IF NOT EXISTS matches (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    season_year INTEGER,
    match_datetime DATETIME,
    court INTEGER,
    home_team TEXT,
    away_team TEXT,
    home_score INTEGER,
    away_score INTEGER
)
