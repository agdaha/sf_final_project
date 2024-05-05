# Итоговый проект по курсу "Backend-разработчик на Go"

### Состав проекта:

1. API Gateway
2. Сервис загрузки новостей
3. Сервис обработки комментариев
4. Сервис цензурирования комментариев

### API:

##### GET /api/news - получение новостей,

Дополнительные параметры запроса:

- s - строка поиска в наименовании новости
- page - номер отображаемрй страницы

Формат ответа:

```
{
    "News": [
            {
                "Id":  {int}
                "Title": {string}
                "PubDate": {int}
            }, ...
    ]
    "Pages": {
        "Total": {int},
        "Current": {int}
    }
}
```

##### GET /api/news/ - получение новости c id = ,

Формат ответа:

```
{
    "News": {
        "Id": {int},
        "Title": {string},
        "Description": {string},
        "Link": {string},
        "PubDate": {int},
        "Author": {string},
        "Guid" {string}:
        },
    "Comments": [
    {
        "id": {int},
        "author": {string},
        "comment_text": {string},
        "news_id": {int},
        "parent_id": {int, nullable},
        "level" {int}:
    }, ...
    ]
}
```

### POST /api/comments/ - сохранение комментария ,

формат передаваемых данных:

```
{
    "author": {string},
    "comment_text": {string},
    "news_id": {int},
    "parent_id": {int, nullable}
}
```

### Сборка:

Все основные манипуляци можно произвести через make.

`make api_gateway` - сборка  API Gateway

`make censor-service` - сборка Сервис цензурирования комментариев

`make comments-service` - сборка Сервис обработки комментариев

`make news-service` - сборка Сервис загрузки новостей

`make all` - собирает все вышеприведенные сервисы.

Собранные бинарники расположены в папке *builds*

Для запуска приложений используются файлы конфигурации передаваемые через переменную среды CONFIG_PATH. Образцы файлов конфигурации приложены в директориях {module_name}/configs

### Запуск Docker контейнеров:

Протестировать приложения в связке можно через docker, файл docker-compose.yml разположен в корневой директори.

Опции make для работы с контейнерами:

`start_api_gateway`:    			 	docker compose up

`rebuild_and_start_api_gateway`:  	docker compose up --build

`prune`: 							docker image prune #для удаления зависших образов после пересборки

`down`:							docker compose down

Настройки для БД приведены в файле .env

### Тесты для Postman:

Тестовые запросы для Postman в файле для импорта news.postman_collection.json
