# Réflexion sur les moyens de stocker et de partager notre base de données

Conversion d'une base de donnée Sqlite en fichier texte

4 manières plus ou moins complexes :

- Version Gitus avec la philosophie Git intégrée dans l'architecture logiciel (pull, push, merge...)
- Implémentation de gitus dans git comme le logiciel git-bug
- Version avancée avec échange de notre BDD en format texte (intégration avec Git)
- Version simple avec génération d'un fichier texte à partir de BDD, 
ce fichier texte sera ensuite partager simplement à d'utres utilisateurs ou sur Git.

Le principe général serait de faire une exportation/importation d'un fichier BDD vers/depuis un fichier texte. 

Lorsqu'on utilise le logiciel Gitus, l'interaction de données se ferait à partir d'une BDD. Lorsqu'on souhaite partager cette BDD, on la transforme en fichier texte pour pouvoir la transmettre plus facilement notamment avec Git.

Le fihier de texte peut être récupéré par n'importe quel utilisateur. Lorsqu'il utilse Gitus, une conversion de ce fichier est faite en BDD pour pouvoir modifier les données.

Pistes :
- Fichier texte en format JSON ??


# Gestion d'un historique des actions ou des données

2 manières de gérer les différentes versions de la base de données sont ressortis :

- Une méthode "Undo/Redo" qui intéragirait à partir des actions effectuées par les différents utilisateurs de la base de données. Chaque action spécifie l'auteur, date, etc des modifications.
(Design pattern vu en C++ : Undo and Redo)

Pour cette implémentation, la création de 2 fichiers distincts pour spécifier la base de données et
les actions seraient plus simple et compréhensible.


- Historique sur chaque donnée enregistrée avec plusieurs états pour les données : 
(Etats : To do, in progress, done, archived, deleted)
Plusieurs champs à ajouter : Date archivage, date création, auteur des modifications

# Autre

Utilisation et tests sur Git-bug pour comprendre son fonctionnement et les moyens de sauvegarde utilisés
