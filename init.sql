CREATE USER 'auth'@'localhost'     IDENTIFIED BY 'auth123';
CREATE USER 'download'@'localhost' IDENTIFIED BY 'download123';

CREATE DATABASE auth;
CREATE DATABASE soundCloudDownloads;

GRANT ALL ON auth.*                TO 'auth'@'localhost';
GRANT ALL ON soundCloudDownloads.* TO 'auth'@'localhost';

USE auth;

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

#                                                                       base64
INSERT INTO users (name,email,password) VALUES ("root","root@root.com","cm9vdDEyMw==");
#                																		 base64
INSERT INTO emailConf (name,email,token,sndDate,expDate) VALUES ("root","root@root.com","dG9rZW4=",current_timestamp(),current_timestamp() + INTERVAL 1 DAY);


# ---------------------------------------------------------------
USE soundCloudDownloads;

CREATE TABLE songs (
	id       int          not null auto_increment primary key,
	name     varchar(255) not null unique,
	author   varchar(255) not null,
	length   int          not null,
	request  int          not null
)
