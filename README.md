# Gitus

## Setup du projet

- Vérifier que la variable d'environnement GOPATH indique le path vers "%USERPROFILE%\Go" ( suivre le tuto de <https://www.java.com/fr/download/help/path.xml,> pour regarder le contenu de la variable et la modifier si besoin)
- Dans le dossier pointer par GOPATH, créer si ce n'est pas déjà fait le dossier "src"
- Dans ce dossier, créer un autre dossier "github.com" puis dans ce dossier créer le dossier "abuan"
- Clôner le dépot dans ce dossier, soit dans : "%USERPROFILE%\Go\src\github.com\abuan"

## Instalation BDD SQLite

- Télécharger SQLite avec le tuto suivant : <https://www.sqlitetutorial.net/download-install-sqlite/>
- Télécharger  GCC en suivant le tuto suivant : <https://medium.com/@yaravind/go-sqlite-on-windows-f91ef2dacfe>
- Ajouter dans votre variable path le chemin vers le dossier bin du logiciel installé
- Relancer Visual Studio Code, dans le répertoire du projet gitus exécuter les commandes :
- "go get -u github.com/mattn/go-sqlite3@v1.12.0"
- "go install github.com/mattn/go-sqlite3"

Une fois l'instalation de SQLite effectuée il faut installer notre base de donnée "gitus"

- Aller dans le dossier du projet puis dans le dossier "bdd" puis dans le dossier "install"
- Vous trouverez des fichier .bat qui sont des exécutables
- Lancer le fichier install_BDD.bat

La base de donnée gitus doit être créée dans C:\sqlite\

## Modification des champs d'une table ou ajout d'une table

Pour modifier les champs d'une table ou ajouter une table il faut le faire dans le fichier "setup_database.sql". Ce fichier contient toutes les instructions nécessaires pour la création de la base de données et de ces tables.
Si des modifications sont apportées à ce fichier (par vous ou lorqsque que vous faites un "pull" sous GIT) cela signifie que votre BDD actuelle n'est pas à jour et qu'il faut la mettre à jour.

## Mise à jour de la BDD

La mise à jour de votre BDD se fait en deux étapes.

- Supprimer votre BDD avec l'exécutable "uninstall_database.bat" dans le dossier "gitus/bdd/install"
- Réinstaller la BDD avec l'exécutable "setup_database.bat" dans le dossier "gitus/bdd/install"

Cette opération vous fera perdre toutes les données contenu dans votre base. Une fois le projet bien avancé la base sera peuplée lors de la création supprimant ce problème.

## Utilisation du Makefile pour build et run des tests

Le Makefile est un script regroupant un ensemble de commandes permettant de build notre projet. Pour utiliser le Makefile, c'est à dire lancer l'ensemble des commandes d'une étape, comme par exemple l'étape "build", il faut dans le terminal à la racine de notre projet taper "make build". La commande "make" est disponible directement sous Linux mais pas pour Windows, il faut télacharger le logiciel et l'installer.

- Dans la partie "Download" prendre le premier "setup" : <http://gnuwin32.sourceforge.net/packages/make.htm>
- Ajouter le chemin vers "make.exe" dans votre variable d'environnement "PATH"

Si vous ne souhaitez pas utiliser le Makefile, vous pouvez simplement taper dans votre terminal les commandes décrites dans chaque étape en les adaptants légerment pour cibler les bonnes target.
