# VK Graph API

Серверное REST API для работы с графовой базой данных Neo4j. API позволяет выполнять базовые операции с узлами и связями социального графа VK.

## Функциональность

- Получение всех узлов графа с атрибутами
- Получение конкретного узла со всеми его связями
- Добавление узлов и связей
- Удаление узлов и связей
- Защита POST и DELETE операций через API ключ

## Технологии

- Python 3.8+
- FastAPI
- Neo4j
- Pydantic
- pytest

## Установка и запуск

1. Клонируйте репозиторий:

```bash
git clone https://github.com/uumk0n/graphRest.git
cd graphRest
```

2. Установите зависимости:

```bash
pip install -r requirements.txt
```

3. Запустите сервер:

```bash
uvicorn main:app --reload
```
