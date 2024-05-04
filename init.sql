-- Удаление объектов в обратном порядке зависимостей
DROP TRIGGER IF EXISTS trigger_update_comments_structure ON public.comments;
DROP FUNCTION IF EXISTS refresh_comments_structure();
DROP MATERIALIZED VIEW IF EXISTS public.comments_structure;
DROP TABLE IF EXISTS public.comments;

-- Основная таблица
CREATE TABLE IF NOT EXISTS public.comments
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    author text DEFAULT 'anonim'::text,
    text text,
    news_id bigint,
    parent_id bigint REFERENCES public.comments (id)
);

-- VIEW для отображения вложенной иерархии
CREATE MATERIALIZED VIEW comments_structure AS
    WITH RECURSIVE comments_cte(id, author, text, parent_id, news_id,  path) AS (
        SELECT comments.id, comments.author, comments.text, comments.parent_id, news_id, ARRAY [comments.id]
            FROM comments
            WHERE comments.parent_id IS NULL
        UNION ALL
        SELECT comments.id, comments.author, comments.text, comments.parent_id, comments.news_id, array_append(comments_cte.path, comments.id)
            FROM comments_cte,
                 comments
            WHERE comments.parent_id = comments_cte.id
    )
    SELECT *
        FROM comments_cte;

-- функция для обновления VIEW 
CREATE FUNCTION refresh_comments_structure() RETURNS TRIGGER
    LANGUAGE plpgsql AS
$$
BEGIN
    REFRESH MATERIALIZED VIEW comments_structure;
    RETURN new;
END;
$$;

--триггер запуска функции обновления view по любому изменению в основной таблице
CREATE OR REPLACE TRIGGER trigger_update_comments_structure
    AFTER INSERT OR DELETE OR TRUNCATE OR UPDATE 
    ON public.comments
    FOR EACH STATEMENT
    EXECUTE FUNCTION public.refresh_comments_structure();


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