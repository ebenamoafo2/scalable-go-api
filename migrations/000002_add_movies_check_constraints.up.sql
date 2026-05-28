ALTER TABLE movies ADD CONSTRAINT movie_runtime_check CHECK (runtime >= 0);
ALTER TABLE movies ADD CONSTRAINT movie_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
ALTER TABLE movies ADD CONSTRAINT genre_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);