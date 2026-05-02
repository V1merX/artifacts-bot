# Artifacts Bot

Бот для автоматизации игрового процесса [Artifacts MMO](https://artifactsmmo.com).  
A bot for automating gameplay in [Artifacts MMO](https://artifactsmmo.com).

---

## На русском

### Описание

Бот автоматически управляет персонажем в игре Artifacts MMO: перемещает его на нужную локацию, сражается с монстрами и восстанавливает HP при необходимости. Цель — прокачать персонажа до заданного уровня без ручного управления.

### Требования

- Go 1.26+
- Аккаунт на [artifactsmmo.com](https://artifactsmmo.com)
- [Task](https://taskfile.dev) (опционально, для удобного запуска)

### Установка и запуск

1. Склонируй репозиторий:
   ```bash
   git clone https://github.com/V1merX/artifacts-bot.git
   cd artifacts-bot
   ```

2. Создай файл `.env` на основе примера:
   ```bash
   cp .example.env .env
   ```

3. Заполни `.env`:
   ```env
   SERVER_ADDR="https://api.artifactsmmo.com"
   BASIC_AUTH_USERNAME="your@email.com"
   BASIC_AUTH_PASSWORD="your_password"
   ```

4. Запусти бота:
   ```bash
   task run
   # или напрямую:
   go run cmd/bot/main.go
   ```

### Конфигурация

| Переменная              | Описание                        | Обязательная |
|-------------------------|---------------------------------|:------------:|
| `SERVER_ADDR`           | Адрес API сервера               | ✓            |
| `BASIC_AUTH_USERNAME`   | Email аккаунта Artifacts MMO    | ✓            |
| `BASIC_AUTH_PASSWORD`   | Пароль аккаунта Artifacts MMO   | ✓            |

### Архитектура

Проект следует принципам **Clean Architecture**. Зависимости направлены строго внутрь.

```
cmd/bot/             — точка входа
internal/
  domain/            — сущности и бизнес-объекты (не зависят ни от чего)
    character/       — агрегат персонажа (HP, Level, XP, Status)
    action/          — результаты действий (Fight, Move, Rest, Cooldown)
  usecase/           — бизнес-логика
    levelup/         — прокачка персонажа до целевого уровня
  gateway/           — адаптеры к внешним системам
    artifacts/       — реализует интерфейс use case через Artifacts API
  di/                — сборка зависимостей (composition root)
  app/               — запуск приложения
pkg/api/             — сгенерированный клиент Artifacts API (не редактировать)
configs/             — загрузка конфигурации
```

**Правило зависимостей:**
```
app → usecase → domain
          ↑
       gateway
```

`gateway` знает об API и реализует интерфейс, объявленный в `usecase`. `usecase` знает только о `domain`. `domain` не зависит ни от чего.

### Генерация API клиента

Клиент генерируется из OpenAPI-спецификации с помощью [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen):

```bash
task api-gen
```

---

## In English

### Description

A bot that automates character management in Artifacts MMO: moves to target locations, fights monsters, and rests when HP drops below threshold. The goal is to level up a character to a target level without manual input.

### Requirements

- Go 1.22+
- An account on [artifactsmmo.com](https://artifactsmmo.com)
- [Task](https://taskfile.dev) (optional, for convenient running)

### Installation & Running

1. Clone the repository:
   ```bash
   git clone https://github.com/V1merX/artifacts-bot.git
   cd artifacts-bot
   ```

2. Create `.env` from the example:
   ```bash
   cp .example.env .env
   ```

3. Fill in `.env`:
   ```env
   SERVER_ADDR="https://api.artifactsmmo.com"
   BASIC_AUTH_USERNAME="your@email.com"
   BASIC_AUTH_PASSWORD="your_password"
   ```

4. Run the bot:
   ```bash
   task run
   # or directly:
   go run cmd/bot/main.go
   ```

### Configuration

| Variable                | Description                     | Required |
|-------------------------|---------------------------------|:--------:|
| `SERVER_ADDR`           | API server address              | ✓        |
| `BASIC_AUTH_USERNAME`   | Artifacts MMO account email     | ✓        |
| `BASIC_AUTH_PASSWORD`   | Artifacts MMO account password  | ✓        |

### Architecture

The project follows **Clean Architecture** principles. Dependencies point strictly inward.

```
cmd/bot/             — entry point
internal/
  domain/            — entities and value objects (no external dependencies)
    character/       — character aggregate (HP, Level, XP, Status)
    action/          — action results (Fight, Move, Rest, Cooldown)
  usecase/           — business logic
    levelup/         — level up a character to a target level
  gateway/           — adapters to external systems
    artifacts/       — implements the use case interface via Artifacts API
  di/                — dependency wiring (composition root)
  app/               — application runner
pkg/api/             — generated Artifacts API client (do not edit)
configs/             — configuration loading
```

**Dependency rule:**
```
app → usecase → domain
          ↑
       gateway
```

`gateway` knows about the API and implements the interface declared in `usecase`. `usecase` only knows about `domain`. `domain` has no dependencies.

### Regenerating the API Client

The client is generated from the OpenAPI spec using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen):

```bash
task api-gen
```
