CREATE TABLE IF NOT EXISTS films (
		film_id SERIAL PRIMARY KEY,
		title VARCHAR(150) UNIQUE NOT NULL, 
		description VARCHAR(1000),
		date DATE,
		rating INT CHECK(rating BETWEEN 1 AND 10),
		actors_list VARCHAR
		);