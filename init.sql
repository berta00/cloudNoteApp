CREATE USER 'auth'@'localhost' IDENTIFIED BY 'auth123';

CREATE DATABASE auth;

GRANT ALL ON auth.* TO 'auth'@'localhost';

USE auth;

CREATE TABLE users (
	id       int          not null auto_increment primary key,
	name     varchar(255) not null,
	email    varchar(255) not null unique,
	password varchar(255) not null,
	admin    boolean      not null default false,
	date     timestamp    not null default current_timestamp()
);

INSERT INTO users (name,email,password) VALUES ("root","root@root.com","root123");
