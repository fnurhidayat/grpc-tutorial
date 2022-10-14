package postgres

const POSTGRES_MOVIE_REPOSITORY_LIST_SQL = `
SELECT movies.id,
	   movies.title,
	   movies.summary,
	   movies.rating
FROM movies;
`

const POSTGRES_MOVIE_REPOSITORY_SAVE_INSERT_SQL = `
INSERT INTO movies (title, summary, rating)
VALUES (:title, :summary, :rating)
RETURNING id;
`

const POSTGRES_MOVIE_REPOSITORY_SAVE_UPDATE_SQL = `
UPDATE movies
SET title = :title,
    summary = :summary,
    rating = :rating
WHERE id = :id;
`
