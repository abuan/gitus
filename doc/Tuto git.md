
# Tuto git

## Ajouter un fichier
1. créer sa branche dès qu’on rajoute du travail, aller dans branche sur le côté puis clique droit « create branch » feature/nom_pertinent
2. faire sa modif, puis cliquer sur commit, indexer le fichier puis mettre un message de description pour l’action qu’on fait. Si on veut mettre une description plus longue faire deux espaces puis écrire le texte
3. faire push que quand la version est presque finie
4. aller dans github online
5. aller dans pull request, ou new pull request si je vois pas ma branche
6. cliquer sur "create pull request" en mettant description si on veut
7. mettre en commentaire le hash de commit si une personne nous a fait des commentaires et qu'on les a corrigé dans un commit
8. merge la branche dans github en cliquant sur le nom de notre branche
9. delete branch une fois qu'elle a été merged


## Nouvelle branche
1. se placer sur develop
2. clique droit "charger la branche"
3. faire un pullp pour se mettre à jour avec les info à distance
4. créer une nouvelle branche

## Que faire en attendant que la pull request soit accepté et qu'il y a une dépendance avec la suite
1. Créer une nouvelle branche à partir de la feature qui est en attente d'être merge.
2. Quand la merge est faite, rebase la nouvelle feature sur develop
