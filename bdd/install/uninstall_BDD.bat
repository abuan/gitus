@ECHO OFF
ECHO Suppression de la base de donnee "gitus" et de ses tables 
ECHO Veuillez rentrer le mdp MySQL de l'utilisateur "root"
mysql --user=root --password -s < uninstall_database.sql
ECHO La base de donnee a ete correctement desininstalee
@ECHO ON
pause