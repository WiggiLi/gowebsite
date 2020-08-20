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


CREATE TABLE Comments (
	page integer,
	name varchar(255),
	content text
);

CREATE TABLE "accounts" (
	"id" serial, 
	"email" text,
	"password" text,
	"token" text , 
	PRIMARY KEY ("id")
); 


insert into content_page (
       page,
    title,
    content
)
values 
(1, 'Fox', '<p><b>Foxes</b> are small to medium-sized, omnivorous mammals belonging to several genera of the family Canidae. Foxes have a flattened skull, upright triangular ears, a pointed, slightly upturned snout, and a long bushy tail (or brush).</p>'),
(2, 'Bear',  '<p><b>Bears</b> are carnivoran mammals of the family Ursidae.</p> They are classified as caniforms, or doglike carnivorans. Although only eight species of bears are extant, they are widespread, appearing in a wide variety of habitats throughout the Northern Hemisphere and partially in the Southern Hemisphere. Bears are found on the continents of North America, South America, Europe, and Asia. Common characteristics of modern bears include large bodies with stocky legs, long snouts, small rounded ears, shaggy hair, plantigrade paws with five nonretractile claws, and short tails.</p>'),
(3, 'Hare', '<p><b>Hares</b> and <b>jackrabbits</b> are leporids belonging to the genus Lepus.</p> <p>Hares are classified in the same family as rabbits. They are similar in size and form to rabbits and have similar herbivorous diets, but generally have longer ears and live solitarily or in pairs. Also unlike rabbits, their young are able to fend for themselves shortly after birth rather than emerging blind and helpless. Most are fast runners. Hare species are native to Africa, Eurasia and North America. </p>');


