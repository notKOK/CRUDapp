Создать CRUD-приложение, предоставляющее Web API к данным

"Бизнес"-сущность придумывается самостоятельно (на примере books).

API 
Приложение должно уметь посредтсвом http-запросов:
1. Получить по id (Get \book)
2. Получить список всех сущностей (GET \books)
3. Создать сущность (POST/PUT)
4. Обновить сущность по id (POST/PUT)
5. Удалить сущность по id (DELETE)

Go
Использовать gofmt, goimports для code style.
Логировать ошибки любым логгером на Ваш выбор. 
Для передачи данных использовать JSON. 
Для роутинга использовать любую библиотеку на Ваш выбор. 
Хранилище данных на Ваш выбор: PostgreSQL, MSSQL, MongoDB 
Для работы с реляционными БД использовать SQL/ORM

Git 
Ваш код должен находится в git-репозитории
В main ветке должна быть рабочая версия
Файл README обязателен. Добавить в него инструкцию по запуску Вашего приложения


Дополнительные задания со звездочкой
1. Кэширование запросов 
2. Докеризация
3. Docker-compose
4. Авторизация
5. Автоматическое тестирование через Postman
6. Unit-тесты
7. Linter



### Инструкция по установке
1. установить golang. 1.19.8

`apt isntall golang` 
или 
`pacman -S golang`

2. установить Postgresql 15.
для debian:
`apt install postgresql` 
для arch:
`pacman -S postgresql`

3. Запустить Posgresql
`systemctl enable postgresql`
`sytstemct start postgresql`

4. Создание пользователя и базы данных в posgtesql

`psql -U postgres`
`CREATE USER cruder WITH PASSWORD 'jw8';`
`CREATE DATABASE crudapp OWNER cruder;`

5. Скомпилировать программу:
`go build main.go`

6. Запустить программу:
`./main`