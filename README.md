# VK

**Это проект на собеседование в VK** 👾

Примечание:

* Перед запуском сделать миграции бд (пакет migrate Go)

  `$ migrate -source file://path/to/migrations -database postgres://localhost:5432/database up 2`

**Стек:**

* Database: PostgreSql
* Language: Golang
* REST API
* net/http

**Запросы :**

* На все post, delete запросы нужен токен за правами админа
* На get любой валидный токен
* Токен ставится в заголовок Authorization

**Идеи:**

* Можно было бы реализовать проверку уникальности токена через БД, для улучшения безопасности в случае утечки токенов
* Сброс пароля для пользователя

![1710531162041](images/README/1710531162041.png)
