create database web_pages;
\c web_pages;
alter user postgres with encrypted password 'qwerty';
GRANT USAGE ON SHEMA PUBLIC TO POSTGRES;
grant all privileges on database web_pages to postgres; 

CREATE TABLE content_page ( 
	page integer, 
	title varchar(255), 
	content text
);


CREATE TABLE comments (
	page integer,
	name varchar(255),
	content text
);

CREATE TABLE "accounts" (
	"id" serial, 
	"email" text,
	"password" text,
	"token" text , PRIMARY KEY ("id")
); 


