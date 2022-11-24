CREATE USER 'app'@'localhost' IDENTIFIED BY 'pass123';
CREATE DATABASE app;
GRANT ALL ON app.* TO 'app'@'localhost';+

USE app;

CREATE TABLE users (
	id       int          not null auto_increment primary key,
	name     varchar(255) not null,
	email    varchar(255) not null unique,
	password varchar(255) not null,
	admin    boolean      not null default false,
	date     timestamp    not null default current_timestamp()
);
CREATE TABLE emailConf (
	id       int          not null auto_increment primary key,
	name     varchar(255) not null,
	email    varchar(255) not null unique,
	token    varchar(255) not null unique,
    sndDate  datetime     not null,
	expDate  datetime     not null,
	done     boolean      not null default false
);
CREATE TABLE basicNote (
	id       int          not null auto_increment primary key,
	creator  varchar(255) not null,
	content  varchar(255) not null,
	crDate   timestamp    not null default current_timestamp(),
	mfDate   timestamp    not null default current_timestamp()
);

INSERT INTO users (name,email,password) VALUES ("root","root@root.com","cm9vdDEyMw==");
INSERT INTO emailConf (name,email,token,sndDate,expDate) VALUES ("root","root@root.com","dG9rZW4=",current_timestamp(),current_timestamp() + INTERVAL 1 DAY);

