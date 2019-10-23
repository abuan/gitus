# proto_gitus

*** Setup du projet ***

- Vérifier que la variable d'environnement GOPATH indique le path vers "%USERPROFILE%\Go" ( suivre le tuto de https://www.java.com/fr/download/help/path.xml, pour regarder le contenu de la variable et la modifier si besoin)
- Dans le dossier pointer par GOPATH, créer si ce n'est pas déjà fait le dossier "src"
- Dans ce dossier, créer un autre dossier "github.com" puis dans ce dossier créer le dossier "la_ruche_thales"
- Clôner le dépot dans ce dossier, soit dans : "%USERPROFILE%\Go\src\github.com\la_ruche_thales"

*** Instalation BDD MySQL ***

- Télécharger MySQL sur la page : https://dev.mysql.com/downloads/windows/installer/8.0.html
- Suivre le tutoriel openCLassroom pour l'instalation : https://openclassrooms.com/fr/courses/1959476-administrez-vos-bases-de-donnees-avec-mysql/1959969-installez-mysql
- Via le terminal Visual Studio Code, dans le répertoire pointé par la variable GOPATH exécuter la commande : "go get -u github.com/go-sql-driver/mysql"

*** Instalation SQL Developer ***

SQL developer est une interface pour dialoguer avec les bases de données type SQL évitant l'utilisation du terminal
L'outil n'est pas indispensable mais aide grandement pour la visualisation des tables et les informations des tables.

- Télécharger le logiciel SQL Developer : https://www.softpedia.com/get/Internet/Servers/Database-Utils/Oracle-SQL-Developer.shtml
    * SoftPedia Secure Download(US) - x64 JRE
- Télécharger le connecteur MySQL : https://www.softpedia.com/get/Internet/Servers/Database-Utils/MySQL-Connector-J.shtml
    * External miror 1 - v8.0.17
- Suivre le tutoriel suivant pour l'instalation : http://logic.edchen.org/how-sql-developer-connect-mysql/
    * Dans le champ d'édition du port rajouter à la suite du numéro : "/?serverTimezone=UTC#" exemple : "3306/?serverTimezone=UTC#"


