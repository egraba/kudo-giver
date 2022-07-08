DROP TABLE IF EXISTS persons CASCADE;
CREATE TABLE persons (
	id SERIAL PRIMARY KEY NOT NULL,
	first_name VARCHAR(32) NOT NULL,
	last_name VARCHAR(32) NOT NULL,
	email VARCHAR(254) UNIQUE NOT NULL,
	nb_kudos INTEGER DEFAULT 0
);

DROP TABLE IF EXISTS kudos CASCADE;
CREATE TABLE kudos (
	id SERIAL PRIMARY KEY NOT NULL,
	sender_id SERIAL REFERENCES persons,
	receiver_id SERIAL REFERENCES persons,
	message VARCHAR(1024)
);
