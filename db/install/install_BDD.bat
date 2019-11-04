@ECHO OFF
ECHO Creation de la base de donnee "gitus" et de ses tables 
ECHO Veuillez rentrer le mdp MySQL de l'utilisateur "root"
mysql --user=root --password -s < setup_database.sql
ECHO La base de donnee a ete correctement instalee
@ECHO ON
pause