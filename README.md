
# Часть сервиса аутентификации 
### Эндпоинты:
- GET /api/ping 
- GET /api/v1/access выдает пару access и refresh токенов 
- POST /api/v1/refresh производит refresh операцию для пары access и refresh токенов 



### Описание переменных окружения, и их значений по умолчанию.
go 1.24

| Переменная        | Описание                               | Значение по умолчанию                                                                               |
|-------------------|----------------------------------------|-----------------------------------------------------------------------------------------------------|
| PQ_DSN            | Строка подключения к PostgreSQL        | `host=auth_postgres port=5432 user=auth_user password=auth_password dbname=auth_db sslmode=disable` |
| JWT_SECRET        | Секретный ключ для подписи JWT-токенов | `secret`                                                                                            |
| ACCESS_TOKEN_EXP  | Время жизни Access Token               | `5m`                                                                                                |
| REFRESH_TOKEN_EXP | Время жизни Refresh Token              | `1h`                                                                                                |
| BCRYPT_COST       | Сложность хеширования паролей          | `6`                                                                                                 |
| TLS_KEY           | Путь к файлу с TLS-ключом              | `/path/to/tls/key`                                                                                  |
| HOST              | Порт поднятия сервиса                  | `:1235`                                                                                             |
| TLS_PEM           | Путь к файлу с TLS-сертификатом        | `/path/to/tls/pem`                                                                                  |

### Запуск 
### `go run ./cmd/app/main.go`

### Сборка 
### `go build ./cmd/app/main.go`

### Запуск в docker 
### `cd .\docker\`
### `docker compose up`