# Specification du projet GITUS

## Environnement du projet

 Le Business Owner (BO) a pour rôle de faire le lien entre le client et l'équipe. Il est le principal interlocuteur des clients. Il est chargé de restranscrire le besoin du client au reste de l'équipe. Il définit la priorité de chaque User stories en prenant en compte l'importance de la fonctionnalité pour le client, l'estimation donnée par l'équipe et le facteur de risque.
 Suite à son passage chez le client, son travail consiste en plusieurs tâches :
 1. Création/supprission de projet
 2. Ajout/suppression/modification de stories
 3. estimation de la Business value (de 0 à 1000) 
 4. partage des modifications apportées avec les PO

Un Product Owner (PO) créé des user stories techniques à partir des user story fonctionnelles définies par le BO après sa discussion avec le client. Il contruit dans un premier temps le "Results Backlog" qui consiste en une liste ordonnée des livrables qui seront livrés au client pour mettre en valeur l'implémentation de nouvelles fonctionnalités. Suite à cela, de par sa vision technique, il organise le travail des développeurs afin d'être en mesure de respecter les dates des livrables fixées en amont. Il aide donc à la construction du "Work Backlog", composée de l'ensemble des userstory sur lesquelles l'équipe devra travailler. 
Après une séance d'échange avec l'équipe, son travail se découpe en plusieurs tâches :
1. décomposition de stories tout en conservant la traçabilité de celles-ci
2. estimation de l'effort via des story points selon la suite de Fibonacci.

Actuellement, le BO utilise le logiciel tyga pour synthétiser sa conversation avec le client et décliner les besoins du clients en user stories fonctionnelles. Il dispose d'un PC portable sous Window 7.
Cet ordinateur lui permet :
- d'accéder au répertoire partagé entre BO. 
- de se connecter au réseau de l'entreprise.
Toutefois, il est soumis aux contraintes suivantes :
- impossibilité d'installer une nouvelle application via MSI ou INstallShield
- exécution des applications non validées par le groupe Thalès uniquement depuis un répertoire spécifique nommé "MyApp"
- impossibilité d'ouvirr un port même local.
Il est amené à se déplacer chez le client avec son ordinateur portable. 

Le PO, quant à lui,  remplit un fichier excel partagé dans lequel figure les user stories techniques sur lequel l'équipe de développeur s'appuira pour travailler. Il dispose d'un PC portable, identique à celui du BO et d'un PC fixe sous Window 10. Ce PC fixe, située sur une bulle isolée qui n'a pas accès au partage de dossier de l'entreprise, lui permet d'accéder au répertoire partagé entre PO en s'y connectant par connexion réseau et d'obtenir des droits administrateurs. 
 
La communication entre BO et PO se fait à l'oral ou bien par clé USB pour faire transiter les informations distantes d'un PC fixe de PO à un PC portable de BO.

## Objectifs du projet

L'objectif principal du projet est de faciliter la communication en entreprise entre les PO et les BO concernant la mise en commun des user stories, actuellement contrainte par les consignes de sécurité mises en place par l'entreprise. L'idée est d'avoir un outil qui permettrait de :
- rendre rapide et efficace la modification des userstories propres au BO et propres au PO tout en gérant les problèmes de conflits 
- d'accéder aux PO aux user stories des BO et inversement
- de mettre en lien les user stories de PO avec celles des BO tout en gérant leur synchronisation.

En effet, il faudrait faire en sorte que les outils de gestion de projet contenant les user stories respectives puissent communiquer et se mettre à jour automatiquement. 

## Mise en oeuvre des solutions

Pour implémanter une solution, il est nécessaire d'évaluer la faisabilité des trois pistes suivantes afin de convenir de la solution optimale :
1. Collaborer avec un logiciel opensource Git Bug pour ajouter les fonctionnalités qui s'adaptent à notre besoin
2. S'appuyer de Git Bug pour le faire évoluer de notre côté afin qu'il réponde à nos attentes
3. Coder notre propre logiciel en partant de 0

## Définition des objets

### Projet 

Un projet est le point de départ. Un projet a un nom, un ID, un objectif, des noms de personnes qui lui sont associés, des user stories.

Fonctions : créer, supprimer, modifier, afficher.

### User story 
Au lancement d’un projet, une réunion « kick off » a lieu entre le testeur, les développeurs, le product owner et l’expert métier. Après discussion entre l’équipe, des spécifications simples appelées User stories sont établies.
Une user story est la formalisation d’un besoin associé à une fonctionnalité à réaliser. 
Une user story comporte un ID, une description, un effort, (un statut ?), des tâches et/ou des users stories. Certaines user stories peuvent s’emboiter les unes dans les autres tant qu’elles ne sont pas découplables en simple tâche.
* L’ID permet d’identifier facilement l’user story. Chiffre allant de 0 à N et s’implémente en fonction de l’ordre de création des user story.
* La description doit être courte et doit avoir le format suivant : En tant que « fonction de la personne », je souhaiterais + verbe à l’impératif…
* L’effort permet de prioriser les user stories entre elles  et peut prendre les valeurs suivantes 0,1,3,5,8,13.
Chaque user story peut passer individuellement à l’étape supérieure.
Des comportements ou fonctionnalités que l’équipe a oublié de spécifier sont notées pour la prochaine itération.

Fonctions : créer, supprimer, modifier, afficher.

### Tâche
Une tâche une courte user story correspondant à exigence technique dérivée d’une user story. Elle permet au développeur de mettre de la solution technique.
Une tâche comporte un ID, un nom, un effort, un ou plusieurs exécuteurs, un statut (TODO, IN PROGRESS, TO VERIFY).
Une tâche ne peut pas contenir une autre tâche.

Fonctions : créer, supprimer, modifier, afficher.
