DROP TABLE IF EXISTS news;
CREATE TABLE news (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    link TEXT NOT NULL UNIQUE,
    pub_date INTEGER DEFAULT 0,
    author TEXT,
    guid TEXT
);




