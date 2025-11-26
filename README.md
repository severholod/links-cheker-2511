# Тестовое задание для Workmate
## Links checker

### БД для хранения ссылок - sqlite, протокол - http, формат - JSON

### Эндпоинты:
`/links Method POST
 request - {"links": []string}
 response - {"links": { ${url}: string }, links_num: int }
`


`/report Method POST
request - {"links_list": []int}
response - report[timestamp].pdf
`

### Запуск
Установить в ENV-переменную `CONFIG_PATH` путь до файла конфигурации

`go mod tidy`

`go run ./cmd/main.go`