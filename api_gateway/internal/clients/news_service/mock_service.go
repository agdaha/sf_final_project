package newsservice

import (
	"context"
	"errors"
)

var _ NewsService = &Mock_client{}

type Mock_client struct {
}

func (c *Mock_client) GetNews(ctx context.Context, page int, search string) ([]byte, error) {
	if page == 13 {
		return []byte(""), errors.New("error")
	}
	return []byte(`{
		"News": [
			{
				"Id": 169,
				"Title": "Библиотека GopherJS в Golang",
				"Description": "GopherJS позволяет переводить Go-код в JavaScript — он предоставляет полноценную совместимость с большинством пакетов стандартной библиотеки Go. Также Gopher поддерживает горутины и каналы!В статье в общих деталях рассмотрим эту за ...",
				"PubDate": 1714829201
			},
			{
				"Id": 44,
				"Title": "Паттерн Identity Map в Golang",
				"Description": "Привет, Хабр!Identity Map — это паттерн проектирования, предназначенный для управления доступом к объектам, которые загружаются из базы данных. Основная его задача — обеспечить, чтобы каждый объект был загружен только один р ...",
				"PubDate": 1714283854
			},
			{
				"Id": 45,
				"Title": "Дайджест полезных материалов из мира Golang за неделю (21.04.24 -28.04.24)",
				"Description": "Подборка полезных материалов и находок из мира Go за неделю. 🛠 Инструменты недели:&nbsp;&nbsp;• go-size-analyzer - инструмент для анализа размера зависимостей в скомпилированных бинарных файлах&nbsp;Go.&nbsp;• Go-mongox - пакет Go Mongo, поддержив ...",
				"PubDate": 1714246228
			},
			{
				"Id": 53,
				"Title": "Как мы ускорили Golang-тесты на CI",
				"Description": "Привет, Хабр 👋! Меня зовут Александр, я занимаюсь разработкой ПО. В этом посте я расскажу про свой опыт, как желание улучшить свой рабочий процесс CI, помогло ускорить все golang пайплайны в  PaaS в СберМаркета. Читать далее",
				"PubDate": 1713254461
			},
			{
				"Id": 55,
				"Title": "Дайджест полезных материалов из мира Golang за неделю (6.04.24 -13.04.24)",
				"Description": "🔥 Дайджест полезных материалов из мира Golang за неделюПодборка полезных материалов и находок из мира Go за неделю. Смотреть",
				"PubDate": 1713022408
			},
			{
				"Id": 61,
				"Title": "Дайджест полезных материалов из мира Golang за неделю",
				"Description": "Подборка полезных материалов и репозиториев представляет собой ценный ресурс для всех поклонников Go, желающих быть в курсе последних тенденций и развития языка. Давайте вместе исследуем самые актуальные и интересные н ...",
				"PubDate": 1712391853
			},
			{
				"Id": 63,
				"Title": "Beego в Golang для начинающих",
				"Description": "Привет, Хабр!Beego – это фреймворк для разработки веб-приложений на языке Go, ориентированный на быстрое развертывание и простоту использования. В его основе лежит идея создания полнофункциональных приложений с минимум ус ...",
				"PubDate": 1712336383
			},
			{
				"Id": 68,
				"Title": "Golang: Мои Открытия",
				"Description": "В этом посте мы обсудим несколько увлекательных моментов, которые я узнал в процессе своей работы. В&nbsp;нашем канале много переводов стаей и обзор инструментов&nbsp; GO, welcome.1. Как используется встраивание (embedding) в Go?В Go дирек ...",
				"PubDate": 1711866668
			}
		],
		"Pages": {
			"Total": 1,
			"Current": 1
		}
	}`), nil
}

func (c *Mock_client) GetNewsDetailed(ctx context.Context, id int) ([]byte, error) {
	if id == 13 {
		return []byte(""), errors.New("error")
	}
	return []byte(` {
        "Id": 26,
        "Title": "Как потреблять API с ограничением по RPS в .NET приложениях",
        "Description": "<a href=\"https://habr.com/ru/companies/ruvds/articles/804025/\"><img src=\"https://habrastorage.org/webt/ex/mv/k0/exmvk08mftvpjzefnkulgw_p-78.png\"></a><br>\n<br>\nОднажды каждый C# программист получает на работе задачу по разработке интеграции с внешней системой, где <b><font color=\"#3AC1EF\">ограничена максимальная частота запросов в секунду</font></b>.<br>\n<br>\nИнтернет яростно сопротивлялся предоставить мне инструкцию к написанию такого кода, закидывая туториалами по настройке ограничения RPS на сервере, а не <b><font color=\"#3AC1EF\">клиенте</font></b>.<br>\n<br>\nНо теперь на Хабре есть эта статья, которая научит отправлять запросы из <code>HttpClient</code> так, чтобы не получать <code>429 Too Many Requests</code>. <a href=\"https://habr.com/ru/articles/804025/?utm_source=habrahabr&amp;utm_medium=rss&amp;utm_campaign=804025#habracut\">Читать дальше &rarr;</a>",
        "Link": "https://habr.com/ru/companies/ruvds/articles/804025/?utm_source=habrahabr&utm_medium=rss&utm_campaign=804025",
        "PubDate": 1714640418,
        "Author": "",
        "Guid": ""
    }`), nil
}
