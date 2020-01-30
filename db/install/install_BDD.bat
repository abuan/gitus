@ECHO OFF
ECHO Creation de la base de donnee "gitus" avec SQLite et de ses tables 
sqlite3 C:\sqlite\gitus.db < setup_database.sql
@ECHO ON
pause