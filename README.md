## Docker образ RabbitMQ + http server

### Запуск образа

```bash
docker-compose up -d
```

### Сервисы

### Сервисы

- **rabbitmq**

    port: 5673

    Сервис RabbitMQ 

- **rabitmq management**

    port: 15673

    Админка RabbitMQ. Доступ [http://localhost:15673/](http://localhost:15673/)

- **HTTP API**

    port: 3001

    API для отправки в очереди

### Отправка стоков  в очередь

Формат JSON

```json
{
    "exchange": "stocks",
    "key": "offline",
    "data": [
        {
            "barcode": "20007464733",
            "store_id": 113,
            "available": 1,
            "reserved": 0
        },
        {
            "barcode": "20007449231",
            "store_id": 199,
            "available": 6,
            "reserved": 0
        }
    ]
}
```

- Описание полей выгрузки
    - exchange - название Exchange в RabbitMQ
    - key - ключ сообщения (для отправки подписчикам)
    - data - массив обновленных остатков

Для отправки стоков необходимо сделать запрос

```json
POST http://0.0.0.0:3001/stocks 
Content-Type: application/json
```

### Настройка "федерации"

Для включения "федерации" необходимо настроить политики через админку RabbitMQ