**grabber**

**Сборка сервиса**

`make build`

**Запуск**

Для запуска сервиса требуется Docker с поддержкой swarm. Перед
первым запуском необходимо выполнить `make init` для инициализации окружения.

`make run` запустит сервер Postgres и сервис сбора новостей

**Остановка сервиса**

`make stop`

**Запуск тестов**

`make local-test`

**Работа с сервисом**

Добавить источник парсинга (lenta.ru):

```
curl --location --request POST 'http://127.0.0.1:8701/api/v1/grabber/jobs' \
 --header 'Content-Type: application/json' \
 --data-raw '{
   "url": "https://lenta.ru",
   "query": {
     "itemsSelector": {
       "xpath": "//div[@class='\''span4'\'']/div[@class='\''item'\'']/a"
     },
     "titleSelector": {
       "value": "text"
     },
     "linkSelector": {
       "value": "attr",
       "attr": "href"
     },
     "descriptionSelector": {
       "xpath": "//div[@itemprop='\''articleBody'\'']",
       "value": "inner_text"
     },
     "followLinkForDescription": true
   },
   "periodSec": 60
 }'

 ```

Добавить источник парсинга (Rambler):

```
curl --location --request POST 'http://127.0.0.1:8701/api/v1/grabber/jobs' \
--header 'Content-Type: application/json' \
--data-raw '{
  "url": "https://news.rambler.ru",
  "query": {
    "itemsSelector": {
      "xpath": "//div[@class='\''top-main__news-item'\'']"
    },
    "titleSelector": {
      "xpath": "//div[@class='\''top-card__title'\'']",
      "value": "text"
    },
    "linkSelector": {
      "xpath": "//a",
      "value": "attr",
      "attr": "href"
    },
    "descriptionSelector": {
      "xpath": "//meta[@itemprop='\''articleBody'\'']",
      "value": "attr",
      "attr": "content"
    },
    "followLinkForDescription": true
  },
  "periodSec": 60
}'

```

Получить список активных источников (джобов):
```
curl --location --request GET 'http://127.0.0.1:8701/api/v1/grabber/jobs'
```

```
[
    {
        "id": "1c66cf03-b450-4f9a-9810-58b8ba4ae8ae",
        "periodSec": 60,
        "query": {
            "descriptionSelector": {
                "value": "inner_text",
                "xpath": "//div[@itemprop='articleBody']"
            },
            "followLinkForDescription": true,
            "itemsSelector": {
                "xpath": "//div[@class='span4']/div[@class='item']/a"
            },
            "linkSelector": {
                "attr": "href",
                "value": "attr"
            },
            "titleSelector": {
                "value": "text"
            }
        },
        "url": "https://lenta.ru"
    },
    {
        "id": "84c066e3-adf7-4ac3-8395-02373a0c06f4",
        "periodSec": 60,
        "query": {
            "descriptionSelector": {
                "attr": "content",
                "value": "attr",
                "xpath": "//meta[@itemprop='articleBody']"
            },
            "followLinkForDescription": true,
            "itemsSelector": {
                "xpath": "//div[@class='top-main__news-item']"
            },
            "linkSelector": {
                "attr": "href",
                "value": "attr",
                "xpath": "//a"
            },
            "titleSelector": {
                "value": "text",
                "xpath": "//div[@class='top-card__title']"
            }
        },
        "url": "https://news.rambler.ru"
    }
]
```

Удалить источник (джоб):
```
curl --location --request DELETE 'http://127.0.0.1:8701/api/v1/grabber/jobs/1c66cf03-b450-4f9a-9810-58b8ba4ae8ae'
```

Получить список статей:
```
curl --location --request GET 'http://127.0.0.1:8701/api/v1/grabber/articles?page=1&rowsPerPage=20'
```

Поиск статьи:
```
curl --location --request GET 'http://127.0.0.1:8701/api/v1/grabber/articles?page=1&rowsPerPage=20&q=%D0%BC%D0%B0%D0%BB%D1%8B%D1%88%D0%B5%D0%B2%D0%B0%20%D1%83%D1%81%D1%82%D1%80%D0%B8%D1%86%D0%B0'
```

```
[
    {
        "link": "https://lenta.ru/news/2020/10/23/ustric/",
        "title": "Врачи оценили совет Малышевой есть устрицы в пандемию"
    },
    ...
]
```

```
curl --location --request GET 'http://127.0.0.1:8701/api/v1/grabber/articles?page=1&rowsPerPage=20&q=%D0%B2%D0%B0%D0%BA%D1%86%D0%B8%D0%BD%D0%B0%20%D0%BF%D0%BE%D1%82%D1%80%D0%B5%D0%B1%D0%BE%D0%B2%D0%B0%D0%BB%D0%B8'
```

```
[
    {
        "link": "https://news.rambler.ru/world/45078829-poka-vy-spali-rossiyskuyu-vaktsinu-potrebovali-cherez-sud/",
        "title": "Пока вы спали: российскую вакцину потребовали через суд"
    },
   ...
]
```