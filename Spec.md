Specification du projet GITUS
==============================

Définition des objets
------------------------------


#Projet 
Un projet est le point de départ. Un projet a un nom, un ID, un objectif, des noms de personnes qui lui sont associés, des user stories.

Fonctions : créer, supprimer, modifier, afficher.

#User story 
Au lancement d’un projet, une réunion « kick off » a lieu entre le testeur, les développeurs, le product owner et l’expert métier. Après discussion entre l’équipe, des spécifications simples appelées User stories sont établies.
Une user story est la formalisation d’un besoin associé à une fonctionnalité à réaliser. 
Une user story comporte un ID, une description, un effort, (un statut ?), des tâches et/ou des users stories. Certaines user stories peuvent s’emboiter les unes dans les autres tant qu’elles ne sont pas découplables en simple tâche.
* L’ID permet d’identifier facilement l’user story. Chiffre allant de 0 à 100 et s’implémente en fonction de l’ordre de création des user story.
* La description doit être courte et doit avoir le format suivant : En tant que « fonction de la personne », je souhaiterais + verbe à l’impératif…
* L’effort permet de prioriser les user stories entre elles  et peut prendre les valeurs suivantes 0,2,3,5,8,13.
Chaque user story peut passer individuellement à l’étape supérieure.
Des comportements ou fonctionnalités que l’équipe a oublié de spécifier sont notées pour la prochaine itération.

Fonctions : créer, supprimer, modifier, afficher.

#Tâche
Une tâche une courte user story correspondant à exigence technique dérivée d’une user story. Elle permet au développeur de mettre de la solution technique.
Une tâche comporte un ID, un nom, un effort, un ou plusieurs exécuteurs, un statut (TODO, IN PROGRESS, TO VERIFY).
Une tâche ne peut pas contenir une autre tâche.

Fonctions : créer, supprimer, modifier, afficher.
