WITH
    points AS (
        SELECT
            season_year,
            home_team as team,
            CASE
                WHEN home_score > away_score THEN 3
                WHEN home_score < away_score THEN 0
                ELSE 1
            END as points
        FROM
            matches
        UNION ALL
        SELECT
            season_year,
            away_team as team,
            CASE
                WHEN home_score < away_score THEN 3
                WHEN home_score > away_score THEN 0
                ELSE 1
            END as points
        FROM
            matches
    )
SELECT
    season_year,
    team,
    SUM(points) as total_points
FROM
    points
GROUP BY
    season_year,
    team
ORDER BY
    season_year DESC,
    total_points DESC;
