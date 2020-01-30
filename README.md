# Gitus

## Motivation

Ce petit logiciel en ligne de commande a été réalisé dans le cadre d'un projet industriel en 5 ème année à l'Insa Rennes.

Gitus est un prototype d'un logiciel visant à aider la gestion de projet en méthode Agile. Comme des applications telles que Taiga ou Jira, Gitus permet de créer des User Story.
L'objectif est de pouvoir avoir un suivi de ces Stories et de pouvoir travailler en colaboration avec d'autres membre de l'équipe.

## Fonctionnement

Gitus utilise git comme une base de données NOSQL. Le principe est d'englober dans un objet git un fichier JSON décrivant notre Story. Par la suite des opérations sont affectées à cet objet et sont stockées également. A chaque action réalisée sur l'objet un commit est réalisé permettant ainsi d'avoir une gestion de version de notre Story.

Le principe d'utiliser Git comme base de données NOSQL est expliqué dans le lien suivant : <https://www.kenneth-truyers.net/2016/10/13/git-nosql-database/>

Gitus est à l'origine un fork du projet du projet Open Source Git-Bug de Micheal Mure : <https://github.com/MichaelMure/git-bug>
Le fork du projet est disponible sur Github également : <https://github.com/abuan/git-bug>

Gitus utilise donc les packages de Git-Bug pour la gestion des repository et de l'architecture de Git en tant que base de données NOSQL.
Gitus à défini le modèle de donnée Story en se basant sur le modèle de donnée "bug" de Git-Bug, nottament pour l'architecture avec des opérations.

## Etat actuel

Pour l'instant Gitus permet seulement de créer des Stories modifier les quelques paramètres de son modèle et de push les modifications ou les pulls dans repo à un autre.
L'objectif serait d'implémenter le modèle complet avec la création de Projets, de Tâches et de pouvoir lier des Stories à des Projets. (voir la specification pour plus d'informations)
