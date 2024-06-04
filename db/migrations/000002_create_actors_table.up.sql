CREATE TABLE IF NOT EXISTS actors (
			actor_id SERIAL PRIMARY KEY,
			name VARCHAR UNIQUE NOT NULL,
			sex  CHAR CHECK(sex='M' OR sex='W'),
			birthday DATE 
		);