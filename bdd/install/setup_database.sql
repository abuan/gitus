-- Création de notre base de données 'gitus'
CREATE DATABASE IF NOT EXISTS gitus;
-- Indication de changement de de BDD
USE gitus;

-- Creation de la table projet
CREATE TABLE IF NOT EXISTS Projet (
	id INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	creation_date DATE,
	description TEXT,
	us_list BLOB,
	PRIMARY KEY (id)
);

-- Creation de la table user_story
CREATE TABLE IF NOT EXISTS User_story(
	id INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	creation_date DATE,
	description TEXT,
	tache_list BLOB,
	PRIMARY KEY (id)
);

-- Creation de la table Tache
CREATE TABLE IF NOT EXISTS Tache (
	id INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	creation_date DATE,
	description TEXT,
	PRIMARY KEY (id)
);

-- Affichage des BDD disppnibles
SHOW DATABASES;

-- Affichage des tables diponibles
SHOW TABLES;