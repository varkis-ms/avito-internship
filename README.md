# Сервис динамического сегментирования пользователей

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

Микросервис для работы с сегментированием пользователей, позволяющий создавать и удалять сегменты, добавлять и исключать пользователей из сегментов, хранить историю сегментирования пользователей и получать отчёт за конкретный период.


---
Table of Contents
---
- [Technology stack](#technology_stack)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Examples](#examples)
- - [Создание сегмента](#create_segment)
- - [Создание сегмента с добавлением N% случайных пользователей](#create_segment_with_random_users)
- - [Удаление сегмента](#delete_segment)
- - [Добавление пользователя в сегменты](#add_user_to_segments)
- - [Добавление пользователя в сегменты на ограниченное время](#add_user_to_segments_with_ttl)
- - [Удаление пользователя из сегментов](#remove_user_from_segment)
- - [Отчёт с экспортом в Google Drive](#report_link)
- - [Отчёт в формате csv файла](#report_file)
- - [Отчёт в формате json](#report_json)
- [Decisions](#decisions)
- [Additional notes](#additional_notes)


---
# Technology stack <a name="technology_stack"></a>
* [![Gin](https://img.shields.io/badge/Gin_Web_Framework-blue?style=plastic&logoColor=yellow&logo=mocha)](https://gin-gonic.com/)
* [![Postgres](https://img.shields.io/badge/PostgreSQL-4169E1?style=plastic&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
* [![Docker](https://img.shields.io/badge/Docker-white?style=plastic&logo=docker&logoColor=2496ED)](https://www.docker.com/)
* [![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=plastic&logo=swagger&logoColor=white)](https://swagger.io/)
* [![pgx](https://img.shields.io/badge/Driver_pgx-blue?style=plastic&logo=adminer&logoColor=white)](https://pkg.go.dev/github.com/jackc/pgx)

Для гибкости и удобства в тестировании и масштабировании проекта, была выбрана Clean архитектура.


# Getting Started <a name="getting_started"></a>

Для запуска сервиса с интеграцией с Google Drive необходимо предварительно:
* Зарегистрировать приложение в Google Cloud Platform: [Документация](https://developers.google.com/workspace/guides/create-project)
* Создать сервисный аккаунт и его секретный ключ: [Документация](https://developers.google.com/workspace/guides/create-credentials)
* Добавить секретный ключ в директорию secrets
* Добавить .env файл в директорию с проектом и заполнить его данными из .env.example,
  указав `GOOGLE_DRIVE_JSON_FILE_PATH=secrets/your_credentials_file.json`

Для запуска сервиса без интеграции с Google Drive достаточно заполнить .env файл,
оставив переменную `GOOGLE_DRIVE_JSON_FILE_PATH` пустой

# Usage <a name="usage"></a>

Сгенерировать .env файл можно командой `make env`

Сервис запускается при помощи команды `make compose-up`

Swagger документацию доступна по адресу `http://localhost:8000/swagger/index.html` (порт по умолчанию 8000)

Для запуска тестов необходимо выполнить команду `make test`, для запуска тестов с покрытием `make cover`

Для запуска линтера необходимо выполнить команду `make linter`


# Examples <a name="examples"></a>

Некоторые примеры запросов
* [Создание сегмента](#create_segment)
* [Создание сегмента с добавлением N% случайных пользователей](#create_segment_with_random_users)
* [Удаление сегмента](#delete_segment)
* [Добавление пользователя в сегменты](#add_user_to_segments)
* [Добавление пользователя в сегменты на ограниченное время](#add_user_to_segments_with_ttl)
* [Удаление пользователя из сегментов](#remove_user_from_segment)
* [Отчёт с экспортом в Google Drive](#report_link)
* [Отчёт в формате csv файла](#report_file)
* [Отчёт в формате json](#report_json)


## Создание сегмента <a name="create_segment"></a>
```
curl -X 'POST' \
  'http://localhost:8000/api/v1/segment/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "segment": "AVITO_VOICE_MESSAGES"
}'
```

Пример ответа:
```
{
  "message": "created"
}
```


## Создание сегмента с добавлением N% случайных пользователей <a name="create_segment_with_random_users"></a>
```
curl -X 'POST' \
  'http://localhost:8000/api/v1/segment/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "percent": 0.5,
  "segment": "AVITO_VOICE_MESSAGES"
}'
```

Пример ответа:
```
{
  "message": "created"
}
```


## Удаление сегмента <a name="delete_segment"></a>
```
curl -X 'DELETE' \
  'http://localhost:8000/api/v1/segment/delete' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "segment": "AVITO_VOICE_MESSAGES"
}'
```

Пример ответа:
```
{
  "message": "deleted"
}
```


## Добавление пользователя в сегменты <a name="add_user_to_segments"></a>
```
curl -X 'POST' \
  'http://localhost:8000/api/v1/user/add' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "segments": [
    "AVITO_VOICE_MESSAGES",
    "AVITO_PERFORMANCE_VAS"
  ],
  "user_id": 1000
}'
```

Пример ответа:
```
{
  "message": "added"
}
```


## Добавление пользователя в сегменты на ограниченное время <a name="add_user_to_segments_with_ttl"></a>
```
curl -X 'POST' \
  'http://localhost:8000/api/v1/user/add' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "segments": [
    "AVITO_VOICE_MESSAGES",
    "AVITO_PERFORMANCE_VAS"
  ],
  "ttl": 24,
  "user_id": 1000
}'
```

Пример ответа:
```
{
  "message": "added"
}
```

Примечание к методу:
> ttl задаётся в часах, то есть для добавления пользователя на сутки, необходимо указать ttl = 24.


## Удаление пользователя из сегментов <a name="remove_user_from_segment"></a>
```
curl -X 'DELETE' \
  'http://localhost:8000/api/v1/user/remove' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "segments": [
    "AVITO_VOICE_MESSAGES",
    "AVITO_PERFORMANCE_VAS"
  ],
  "user_id": 1000
}'
```

Пример ответа:
```
{
  "message": "removed"
}
```


## Отчёт с экспортом в Google Drive <a name="report_link"></a>
```
curl -X 'GET' \
  'http://localhost:8000/api/v1/report/link?month=8&year=2023' \
  -H 'accept: application/json'
```

Пример ответа:
```
{
  "Link": "https://drive.google.com/file/d/12LrddKyrL4KEFz4xOe4rTblG35h1Dap3/view?usp=sharing"
}
```


## Отчёт в формате csv файла <a name="report_file"></a>
```
curl -X 'GET' \
  'http://localhost:8000/api/v1/report/file?month=8&year=2023' \
  -H 'accept: text/csv'
```

Пример ответа:
```
Скачивается файл csv
```


## Отчёт в формате json <a name="report_json"></a>
```
curl -X 'GET' \
  'http://localhost:8000/api/v1/report/?month=8&year=2023' \
  -H 'accept: application/json'
```

Пример ответа:
```
[
  {
    "user_id": "1000",
    "segment": "AVITO_VOICE_MESSAGES",
    "operation": "add",
    "date": "2023-08-30T19:04:52.406104+03:00"
  },
  {
    "user_id": "1",
    "segment": "AVITO_VOICE_MESSAGES",
    "operation": "add",
    "date": "2023-08-30T19:10:29.339163+03:00"
  },
  {
    "user_id": "2",
    "segment": "AVITO_VOICE_MESSAGES",
    "operation": "add",
    "date": "2023-08-30T19:10:33.352127+03:00"
  },
  {
    "user_id": "100",
    "segment": "AVITO_VOICE_MESSAGES",
    "operation": "add",
    "date": "2023-08-30T19:10:37.365765+03:00"
  },
  {
    "user_id": "1000",
    "segment": "AVITO_VOICE_MESSAGES",
    "operation": "remove",
    "date": "2023-08-30T19:31:51.908592+03:00"
  },
  {
    "user_id": "1",
    "segment": "AVITO_VOICE_MESSAGES",
    "operation": "remove",
    "date": "2023-08-30T19:31:51.908592+03:00"
  }
]
```


# Decisions <a name="decisions"></a>

В процессе выполнения данного ТЗ возникали вопросы, которые были решены следующим образом:
* Что ожидается от сервиса, идемпотентность всех запросов или же нет?
> Я пришел к выводу, что необходимо делать полную идемпотентность, так как сервис является лишь прослойкой и при вызове того же метода POST от сервиса ожидается какое-либо действие по созданию, а если такая сущность уже существует, значит ошибки не было и ответ будет OK.
---
* Может ли пользователь входить и выходить из сегмента несколько раз?
> Я решил создать более логичную, на мой взгляд, ситуацию, когда пользователя могут исключить из сегмента, но после снова добавить его в тот же сегмент. Это немного усложнило разработку, но исправить эту условность можно простым индексом UNIQUE(user_id, segment_id) в таблице users_segment и небольшим переписыванием метода добавления репозитория.
---
* Сервису нужно работать с данными о пользователях, но где их взять?
> По-хорошему, данный сервис должен ходить в другой сервис и подтягивать оттуда данные о пользователях. В моей реализации есть условность, при добавлении пользователя в сегмент сервис сохраняет пользователя в бд, при операциях удаления несуществующего пользователя будет получена ошибка. 
---
* Как корректно выполнить 1 дополнительное задание?
> Изначально сделал метод возвращающий отчет в файле, но в задании было указано, что нужна именно ссылка на файл. Думал сделать метод, который возвращает ссылку на метод для скачивания файла, но решил всё-таки не изобретать велосипед и интегрировался с cервисом Google Drive. На всякий случай решил оставить возможность скачивания файла сразу через http. Также было неясно, отдавать отчет по конкретному пользователю или по всем, как итог отчет составляется по всем пользователям.
---
* Что подразумевается под 3 пунктом основного задания, вроде написано сделать метод по добавлению пользователей в сегменты, а вроде перечислены аргументы и для удаления пользователя из сегментов.
> Всё же решил реализовать 2 разных метода, ведь такой подход, на мой взгляд, делает сервис интуитивно понятнее.

# Additional notes <a name="additional_notes"></a>

## Swagger
![Swagger](https://github.com/varkis-ms/avito-internship/raw/main/img/swagger.png)

---

## Report example
![Report](https://github.com/varkis-ms/avito-internship/raw/main/img/report_example.png)
