
# Person Enrichment Service

Этот проект реализует сервис на языке Go, который позволяет обогатить данные о человеке (ФИО) с использованием открытых API для определения возраста, пола и национальности. Информация сохраняется в базе данных PostgreSQL и доступна через REST API.

## 📋 Описание

Сервис принимает запросы с ФИО, обогащает их информацией о возрасте (имена на кириллице — автоматически конвертируются в латиницу), поле и национальности, а затем сохраняет данные в базе данных. В дальнейшем информация о людях может быть получена с различными фильтрами и пагинацией.

### Основные функции:
- Получение данных о человеке с различными фильтрами и пагинацией
- Удаление записи по уникальному идентификатору
- Обновление информации о человеке
- Добавление новых людей с ФИО через API
- Валидация входных данных
- Логирование на уровнях `debug` и `info`
- Генерация Swagger-документации для API

### Ожидаемый формат входных данных:
```json
{
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich"  // необязательно
}
```

### Обогащение данных:
- Возраст — https://api.agify.io/?name=Dmitriy
- Пол — https://api.genderize.io/?name=Dmitriy
- Национальность — https://api.nationalize.io/?name=Dmitriy

## 🚀 Установка

1. Клонируйте репозиторий:

```
git clone https://github.com/evgeniySeleznev/person-enrichment-service.git
cd person-enrichment-service
```

2. Скачайте зависимости:

```
go mod tidy
```

3. Создайте файл `.env` в корне проекта с настройками подключения к вашей базе данных:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=person_enrichment_db
```

4. Запустите миграции для создания базы данных:

```
go run cmd/api/main.go migrate
```

Для отката миграций используйте скрипт:

```
go run scripts/rollback.go
```

5. Запустите приложение:

```
go run cmd/api/main.go
```

Теперь сервис будет доступен по адресу `http://localhost:8080`.

## 🛠️ API

Документация API доступна через Swagger:

```
http://localhost:8080/swagger/index.html
```

### Эндпоинты:

- **GET /persons** — Получить список людей с фильтрами и пагинацией
- **POST /persons** — Добавить нового человека
- **GET /persons/{id}** — Получить данные по идентификатору
- **PUT /persons/{id}** — Обновить данные по идентификатору
- **DELETE /persons/{id}** — Удалить человека по идентификатору

### Пример запроса на добавление:

```
curl -X POST "http://localhost:8080/persons" -H "Content-Type: application/json" -d '{
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich"
}'
```

### Пример ответа:

```json
{
  "id": 1,
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich",
  "age": 44,
  "gender": "male",
  "nationality": "UA"
}
```

## 🔧 Технологии

- Go (1.24.3)
- PostgreSQL
- Swagger (для документации API)
- go-playground/validator (для валидации данных)
- golang-migrate для настройки миграций в БД
- gorilla/mux для маршрутизации запросов

## 📝 Логирование

В проекте используется логирование с уровнями:

- `debug` — для отладки и подробной информации
- `info` — для основной информации о запросах и ответах
- `error` — для ошибок и исключений
- `fatal` — для критических ошибок

## 📄 Лицензия

Этот проект распространяется под лицензией MIT.
