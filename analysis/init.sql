INSTALL sqlite;

LOAD sqlite;

ATTACH '../src/crawler/app.db' AS db (TYPE sqlite);

USE db;
