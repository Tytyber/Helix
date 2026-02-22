# Документация API сервиса

Этот документ содержит обзор структуры сервиса API, его компонентов и способов их взаимодействия.

## Обзор

API — это бэкенд-сервис на Go, построенный с использованием фреймворка [Gin](https://gin-gonic.com/) для маршрутизации и PostgreSQL через [pgx](https://github.com/jackc/pgx) для доступа к базе данных. Он реализует стандартные RESTful-эндпоинты для аутентификации пользователей, социальных функций (друзья, группы, хабы), реального общения (чат, звонки) и ленты контента.

Приложение следует чистой архитектуре с разделением на:
- **Обработчики** (`hundler.go`): логика обработки HTTP-запросов/ответов
- **Сервисы** (`service.go`): бизнес-логика и взаимодействие с базой данных
- **Модели/Типы**: структуры данных, используемые на всех уровнях

Аутентификация осуществляется через JWT-токены, а refresh-токены хранятся безопасно в базе данных.

## Структура проекта

```
api/
├── go.mod
├── go.sum
├── main.go
└── internal/
    ├── auth/           # Аутентификация и управление пользователями
    ├── calls/          # Голосовые/видео звонки
    ├── chat/           # Личные сообщения
    ├── config/         # Загрузка конфигурации
    ├── database/       # Подключение к базе данных
    ├── feed/           # Персонализированная лента контента
    ├── friend/         # Система друзей
    ├── group/          # Групповые чаты
    ├── hub/            # Сообщества с постами
    └── middleware/     # Общее middleware (например, аутентификация)
```

## Основные компоненты

### Основное приложение (main.go)

Файл `main.go` инициализирует:
- Конфигурацию через переменные окружения
- Пул подключений к базе данных
- Все экземпляры сервисов
- Маршруты HTTP, сгруппированные по функциональности

Маршруты организованы под `/api/v1` с:
- Публичными маршрутами под `/auth`
- Защищенными маршрутами (требуется JWT) под `/protected`

### Конфигурация (internal/config)


Использует `godotenv` для загрузки переменных окружения со значениями по умолчанию:
- `PORT` (по умолчанию: ":8080")
- `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`
- `JWT_SECRET` (по умолчанию: "supersecret")

### База данных (internal/database)

Управляет пулом подключений с использованием `pgxpool`. Пул передается всем сервисам, которым требуется доступ к базе данных.

### Аутентификация (internal/auth)

#### Эндпоинты
- `POST /api/v1/auth/register` — Регистрация нового пользователя
- `POST /api/v1/auth/login` — Вход и получение токенов
- `POST /api/v1/auth/refresh` — Обновление JWT-токена
- `POST /api/v1/auth/logout` — Невалидация refresh-токена
- `GET /api/v1/protected/me` — Получить профиль текущего пользователя
- `PUT /api/v1/protected/profile` — Обновить профиль
- `POST /api/v1/protected/reputation` — Добавить репутацию другому пользователю

#### Реализация
- Пароли хешируются с помощью `bcrypt`
- JWT access-токены (срок действия 15 минут)
- Refresh-токены хранятся в виде SHA256-хешей в БД
- Middleware проверяет `Bearer`-токены

### Хабы (internal/hub)

Сообщества, где пользователи могут публиковать контент.

#### Эндпоинты
- `GET /api/v1/protected/hubs` — Список всех хабов
- `GET /api/v1/protected/hubs/:id` — Получить конкретный хаб
- `POST /api/v1/protected/hubs/:id/join` — Присоединиться к хабу
- `POST /api/v1/protected/hubs/:id/leave` — Покинуть хаб
- `GET /api/v1/protected/hubs/:id/posts` — Получить посты хаба
- `POST /api/v1/protected/hubs/:id/posts` — Создать пост

### Лента (internal/feed)

Персонализированная лента контента на основе хабов пользователя.

#### Эндпоинты
- `GET /api/v1/protected/feed` — Получить ленту пользователя
- `GET /api/v1/protected/feed/top-hubs` — Получить топ хабов по числу участников

### Друзья (internal/friend)

Система запросов в друзья и социального графа.

#### Эндпоинты
- `POST /api/v1/protected/friends/request` — Отправить запрос в друзья
- `POST /api/v1/protected/friends/accept` — Принять запрос
- `GET /api/v1/protected/friends` — Список друзей

### Группы (internal/group)

Приватные групповые чаты с несколькими участниками.

#### Эндпоинты
- `POST /api/v1/protected/groups` — Создать группу
- `POST /api/v1/protected/groups/:id/join` — Присоединиться к группе
- `POST /api/v1/protected/groups/:id/leave` — Покинуть группу
- `GET /api/v1/protected/groups/:id/messages` — Получить сообщения
- `POST /api/v1/protected/groups/:id/messages` — Отправить сообщение

### Чат (internal/chat)

Прямое 1-на-1 общение.

#### Эндпоинты
- `POST /api/v1/protected/chats` — Создать чат с пользователем
- `GET /api/v1/protected/chats/:id/messages` — Получить историю чата
- `POST /api/v1/protected/chats/:id/messages` — Отправить сообщение

### Звонки (internal/calls)

Управление голосовыми/видео звонками.

#### Эндпоинты
- `POST /api/v1/protected/calls` — Инициировать звонок
- `PUT /api/v1/protected/calls/:id/status` — Обновить статус звонка (ожидание, принят, отклонен, завершен)

## Модели данных

### Пользователь
- id, username, email, password_hash, avatar_url, bio, status, reputation, created_at

### Хаб
- id, name, description, avatar, created_at

### HubPost
- id, hub_id, author_id, content, likes, created_at

### FeedPost (только для чтения)
- id, hub_id, hub_name, author_id, content, likes, created_at

### TopHub (только для чтения)
- id, name, description, avatar, members

### Группа
- id, name, owner_id, created_at

### GroupMessage
- id, group_id, sender_id, content, created_at

### Чат
- id, user1_id, user2_id, created_at

### Message
- id, chat_id, sender_id, content, created_at

### Звонок
- id, caller_id, callee_id, status, created_at

## Обработка ошибок

Большинство эндпоинтов возвращают JSON-ошибки:
- `400 Bad Request` — Неверный ввод
- `401 Unauthorized` — Отсутствует/недействителен токен
- `404 Not Found` — Ресурс не найден
- `500 Internal Server Error` — Ошибка на стороне сервера

## Зависимости

См. `go.mod` для полного списка. Основные зависимости:
- `github.com/gin-gonic/gin` — Веб-фреймворк
- `github.com/jackc/pgx/v5` — Драйвер PostgreSQL
- `github.com/golang-jwt/jwt/v5` — Реализация JWT
- `golang.org/x/crypto/bcrypt` — Хеширование паролей
- `github.com/joho/godotenv` — Переменные окружения

## Как запустить

1. Скопируйте `.env.example` в `.env` и настройте
2. Убедитесь, что PostgreSQL запущен
3. Запустите `go run main.go`

## Будущие улучшения

- Добавить валидацию ввода
- Реализовать ограничение частоты запросов
- Добавить middleware для логирования
- Поддержка пагинации в ленте и сообщениях
- Добавить модульные и интеграционные тесты
- Реализовать поддержку WebSocket для обновлений в реальном времени
