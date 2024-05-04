INSERT INTO public.comments(
	author, text, news_id, parent_id)
	VALUES ('author 1', 'comment 1', 1, null),
    ('author 2', 'comment 2', 1, null),
    ('author 3', 'comment 3', 1, 1),
    ('author 2', 'comment 4', 1, 3),
    ('author 1', 'comment 5', 2, null),
    ('author 2', 'comment 6', 2, null),
    ('author 1', 'comment 7', 3, null),
    ('author 2', 'comment 8', 3, null),
    ('author 3', 'comment 9', 3, null),
    ('author 2', 'comment 10', 3, 9),
    ('author 1', 'comment 11', 3, 10);