-- Création de notre base de données 'gitus'
CREATE DATABASE IF NOT EXISTS gitus;
-- Indication de changement de de BDD
USE gitus;

-- Creation de la table projet
CREATE TABLE IF NOT EXISTS Project (
	id INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	creation_date DATETIME,
	descript TEXT,
	PRIMARY KEY (id)
);

-- Creation de la table user_story
CREATE TABLE IF NOT EXISTS UserStory(
	id INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	creation_date DATETIME,
	descript TEXT,
	effort INTEGER UNSIGNED,
	us_list BLOB,
	tache_list BLOB,
	PRIMARY KEY (id)
);

-- Creation de la table Tache
CREATE TABLE IF NOT EXISTS Task (
	id INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	creation_date DATETIME,
	descript TEXT,
	PRIMARY KEY (id)
);

-- Creation de la table Project_structure
-- Lie des US à un Projet
CREATE TABLE IF NOT EXISTS Project_structure (
	project_id INTEGER UNSIGNED NOT NULL,
	userstory_id INTEGER UNSIGNED NOT NULL,
	FOREIGN KEY (project_id) REFERENCES Project(id),
	FOREIGN KEY (userstory_id) REFERENCES UserStory(id)
);

-- Ajoute des Items dans les tables
INSERT INTO Project VALUES (1,"Default_Project",NOW(),"Projet par défaut créé automatiquement");
INSERT INTO UserStory VALUES (1,"Default_US1",NOW(),"User Story par défaut créée automatiquement",0,NULL,NULL);
INSERT INTO UserStory VALUES (2,"Default_US2",NOW(),"User Story par défaut créée automatiquement",0,NULL,NULL);
INSERT INTO UserStory VALUES (3,"Default_US3",NOW(),"User Story par défaut créée automatiquement",0,NULL,NULL);
INSERT INTO Project_structure VALUES(1,1),(1,2),(1,3);

-- Affichage des BDD disppnibles
SHOW DATABASES;

-- Affichage des tables diponibles
SHOW TABLES;