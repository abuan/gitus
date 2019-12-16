# Compte rendu séances

## 04/11/2019

1. Mise au point concernant l'uniformation du code :
- code toujours en anglais (nom des fonctions), commentaire en français
- un fichier créé dans projet_gitus par structure (project, userStory, task). Dans chaque fichier de structure se trouve un fichier structure.go à l'intérieur.

2. Tous les fichiers de doc vont dans le fichier doc pour ne pas polluer le répertoire. 
Il faut mettre à jour le diagramme des données pour qu'il corresponde au nom des fonctions.

3. Création de userStory par ligne de commande à coder :
l'objectif de la séance prochaine est de se familiariser chacun de notre côté avec cobra puis de coder toutes les lignes de commandes définies. 
Pour cela s'inspirer de https://github.com/MichaelMure/git-bug.
Créer un fichier lignedecommande.go par ligne de commande dans fichier commands.
Il ne faut pas oublier de mettre à jour le trello en ajoutant son nom à une tâche pour ne pas travailler en doublon.


## 9/12

manipuler 5 premiers caractères pour l'ID ou plus si pas dissociables
incrémentation si identique

mettre edit plutot que modify ?

BO : vision du client
PO : vision technique

MyApp : répertoire avec des droits d'exécution
fichier .exe et base de données ne sont pas situées au même endroit
toute communication réseau impossible (possibilité usb)
chemin pour accéder

PO : même PC que BO + PC fixe sur bulle isolée qui n'a pas accès au partage de dossier de l'entreprise (=réseau séparé)
BO et PO communique entre eux par clé usb
actuellement BO ont un fichier excel partagé  : en tant que BO, je souhaite que cet attribut 
==> pour communication entre BO uniquement

BO et PO communique entre eux à l'oral, PO utilise tyga pour synthétiser conversation et décline us en technique

remettre en question SQLite éventuellement ? s'inspirer de git bug !!!

Pb : fichier excel plus à jour, ne correspond plus aux entrées de tyga

chiffrage de clé est un chiffrage propriétaire, juste faire un clone qui sera mis dans un répertoire crypté

Schéma de l'environnement :
Deux repo qui existent : Un repo PO  et un repo BO.
Communication entre les deux repo se fait par clé pour faire circuler des informations distantes d'un PC fixe de PO à un PC portable de BO.
PC fixe sous W10 relié au repo PO par connexion réseau.
PC portable sous W7 relié au repo BO par connexion réseau qui va chez le client.


cf site https://git-scm.com/book/en/v2/Distributed-Git-Distributed-Workflows

Gagner si on fait synchro entre deux repo


### 3 possibilités de chemin à prendre : 
1. participer à un logiciel opensource en travaillant en collaboration avec Git bug ==> exercice pas adapté, collaboration difficile 
2. forker et repartir de sa base, faire evoluer la chose en fonction du besoin ==> semble plus réalisable (en citant merci à nanani)
3. coder notre propre logiciel

### US du sprint prochain
- définir le workflow (défini par la clé, le réseau)
- étudier les possibilités de git bug
- essayer de gérer la gestion de conflits entre deux fichiers distants
- faire un graphique pour exprimer clairement le besoin





