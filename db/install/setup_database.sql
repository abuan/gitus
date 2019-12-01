-- Creation de la table projet
CREATE TABLE IF NOT EXISTS `Project` (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	creation_date DATETIME,
	descript TEXT
);

-- Creation de la table user_story
CREATE TABLE IF NOT EXISTS `UserStory`(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title text NOT NULL,
	creation_date DATETIME,
	descript TEXT,
	author TEXT,
	effort INTEGER
);

-- Creation de la table Tache
CREATE TABLE IF NOT EXISTS `Task` (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	creation_date DATETIME,
	descript TEXT
);

-- Creation de la table Project_structure
-- Lie des US à un Projet
CREATE TABLE IF NOT EXISTS `Project_structure` (
	project_id INTEGER UNSIGNED NOT NULL,
	userstory_id INTEGER UNSIGNED NOT NULL,
	FOREIGN KEY (project_id) REFERENCES Project(id),
	FOREIGN KEY (userstory_id) REFERENCES UserStory(id)
);

-- Ajoute des Items dans les tables
INSERT INTO `Project` VALUES (1,"Default_Project1",datetime('now', 'localtime'),"Projet par defaut utilise pour les tests");
INSERT INTO `Project` VALUES (2,"Default_Project2",datetime('now', 'localtime'),"Projet par defaut utilise pour les tests");
INSERT INTO `UserStory` VALUES (1,"Default_US1",datetime('now', 'localtime'),"User Story par defaut utilise pour les tests","auteur",5);
INSERT INTO `UserStory` VALUES (2,"Default_US2",datetime('now', 'localtime'),"User Story par defaut utilise pour les tests","auteur 2",8);
INSERT INTO `UserStory` VALUES (3,"Default_US3",datetime('now', 'localtime'),"User Story par defaut utilise pour les tests","auteur 3",13);
INSERT INTO `Project_structure` VALUES(1,1),(1,2),(1,3),(2,1),(2,3);
