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
response - report.pdf
`

Тк на момент формирования отчета сервис может быть не доступен, хранить его статус в БД нет смысла. Получение статуса происходит непосредственно в момент сохранения ссылок для передачи в теле ответа клиенту и в момент формирования отчета.

### Запуск
Установить в ENV-переменную `CONFIG_PATH` путь до файла конфигурации

`go mod tidy`

`go run ./cmd/main.go`