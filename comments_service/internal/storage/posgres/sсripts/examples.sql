SELECT * FROM comments_structure ORDER BY path;

SELECT * FROM comments_structure WHERE news_id=1 ORDER BY path;

SELECT * FROM comments_structure WHERE 1 = ANY(path) ORDER BY path;;