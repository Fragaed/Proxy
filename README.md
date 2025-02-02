# Proxy Service

## Описание

Это сервис для обработки запросов и взаимодействия с хранилищем данных и внешним API.

## Установка и запуск

### Сборка и запуск приложения

1. Склонируйте репозиторий:
    ```bash
    git clone <repository-url>
    cd <repository-directory>
    ```

2. Сборка приложения:
    ```bash
    make build
    ```

3. Запуск приложения с помощью Docker Compose:
    ```bash
    docker-compose up 
    ```


### Конфигурация

Параметры подключения к хранилищу данных можно задавать как через флаги запуска, так и через переменные окружения.

#### Переменные окружения:

- `DB_HOST`: хост базы данных
- `DB_PORT`: порт базы данных
- `DB_USER`: пользователь базы данных
- `DB_PASSWORD`: пароль базы данных
- `DB_NAME`: имя базы данных

#### Флаги запуска:

- `--db-host`: хост базы данных
- `--db-port`: порт базы данных
- `--db-user`: пользователь базы данных
- `--db-password`: пароль базы данных
- `--db-name`: имя базы данных

Пример запуска с использованием флагов:
```bash
./proxy --db-host=localhost --db-port=5432 --db-user=user --db-password=password --db-name=database
