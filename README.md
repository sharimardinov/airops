# AirOps API

REST API для работы с авиационной базой данных. Предоставляет endpoints для работы с рейсами, аэропортами, самолетами, бронированиями и статистикой.

## Технологический стек

- **Go 1.21+** - основной язык
- **PostgreSQL 15+** - база данных с расширениями (btree_gist, cube, earthdistance)
- **pgx/pgxpool** - PostgreSQL драйвер и connection pooling
- **net/http** - HTTP сервер (stdlib, без фреймворков)
- **Swagger** - API документация

## Особенности

- Clean Architecture (domain → usecase → repositories → handlers)
- Многоязычная поддержка через JSONB (en/ru)
- Транзакционные операции с бронированиями
- Middleware: logging, metrics, rate limiting, request ID, timeouts, panic recovery
- Connection pooling для PostgreSQL
- Swagger UI для документации API

## Быстрый старт

### Предварительные требования

- Go 1.21 или выше
- Docker и Docker Compose
- Make (опционально)

### Запуск через Docker Compose

```bash
# Клонировать репозиторий
git clone <repo-url>
cd airops

# Скопировать .env.example в .env и настроить при необходимости
cp .env.example .env

# Запустить PostgreSQL и API
docker-compose up -d

# Проверить что сервисы запущены
docker-compose ps

# API будет доступен на порту указанном в .env (по умолчанию 8080)
# Swagger UI: /swagger/index.html
```

### Локальный запуск (для разработки)

```bash
# Запустить только PostgreSQL
docker-compose up -d postgres

# Применить миграции (если используются)
# psql -h localhost -U airops -d airops -f schema.sql

# Запустить API
go run ./api/main.go

# Или через Make
make run
```

## API Endpoints

### Health & Monitoring

- `GET /health` - Базовая проверка работоспособности
- `GET /ready` - Проверка готовности (включая подключение к БД)
- `GET /metrics` - Prometheus метрики
- `GET /pool-stats` - Статистика connection pool

### Flights (Рейсы)

- `GET /flights` - Список рейсов с фильтрацией по дате
  - Query params: `date` (required, format: 2006-01-02), `limit` (default: 100)
- `GET /flights/{id}` - Детальная информация о рейсе
- `POST /flights/search` - Поиск рейсов по параметрам
  ```json
  {
    "departure_airport": "SVO",
    "arrival_airport": "LED",
    "departure_date": "2025-06-15",
    "passengers": 2,
    "fare_conditions": "Economy"
  }
  ```

### Airports (Аэропорты)

- `GET /airports` - Список всех аэропортов
- `GET /airports/{code}` - Информация об аэропорте (IATA код)
- `GET /airports/search` - Поиск аэропортов
  - Query params: `city` или `country`

### Airplanes (Самолеты)

- `GET /airplanes` - Список всех типов самолетов
- `GET /airplanes/{code}` - Информация о самолете
- `GET /airplanes/{code}/seats` - Схема посадочных мест
- `GET /airplanes/{code}/stats` - Статистика по самолету

### Passengers (Пассажиры)

- `GET /flights/{id}/passengers` - Список пассажиров на рейсе
  - Query params: `limit` (default: 100)

### Bookings (Бронирования)

- `POST /bookings` - Создать бронирование
  ```json
  {
    "passenger_id": "1234 567890",
    "passenger_name": "IVAN IVANOV",
    "contact_phone": "+79001234567",
    "contact_email": "ivan@example.com",
    "flights": [
      {
        "flight_id": 12345,
        "fare_conditions": "Economy",
        "seat_no": "10A"
      }
    ]
  }
  ```
- `GET /bookings/{book_ref}` - Информация о бронировании
- `GET /bookings/passenger/{passenger_id}` - Бронирования пассажира
- `DELETE /bookings/{book_ref}` - Отменить бронирование

### Statistics (Статистика)

- `GET /stats/routes/top` - Топ популярных маршрутов
  - Query params: `date` (required), `limit` (default: 10)

## Аутентификация

API использует простую аутентификацию через API ключ в заголовке:

```bash
curl -H "X-API-Key: your-secret-key" http://localhost:8080/airports
```

API ключ задается через переменную окружения `API_KEY`.

## Конфигурация

Все настройки задаются через переменные окружения (см. `.env.example`):

```bash
# Database
DATABASE_URL=postgres://airops:secret@localhost:5432/airops

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=60s

# API Security
API_KEY=your-secret-api-key-change-in-production

# Logging
LOG_LEVEL=info

# Rate Limiting
RATE_LIMIT_RPS=100
```

## Разработка

### Структура проекта

```
airops/
├── api/main.go                 # Entry point
├── internal/
│   ├── app/                    # Application layer
│   │   ├── app.go              # App initialization
│   │   └── usecase/            # Business logic
│   ├── config/                 # Configuration
│   ├── domain/                 # Domain models & interfaces
│   │   ├── models/             # Data models
│   │   └── ports.go            # Repository interfaces
│   ├── postgres/               # Data layer
│   │   ├── db.go               # Connection pooling
│   │   └── repositories/       # Repository implementations
│   └── http/                   # HTTP layer
│       ├── handlers/           # HTTP handlers
│       ├── middleware/         # HTTP middleware
│       ├── dto/                # Data transfer objects
│       └── router.go           # Route definitions
├── docs/                       # Swagger documentation
├── tools/gotree/               # Code analysis tool
└── schema.sql                  # Database schema
```

### Добавление нового endpoint

1. Добавить модель в `internal/domain/models/`
2. Добавить интерфейс репозитория в `internal/domain/ports.go`
3. Реализовать репозиторий в `internal/postgres/repositories/`
4. Создать usecase в `internal/app/usecase/`
5. Добавить handler в `internal/http/handlers/`
6. Зарегистрировать route в `internal/http/router.go`
7. Добавить Swagger аннотации

### Запуск тестов

```bash
# Unit тесты
go test ./... -v

# Unit тесты с покрытием
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out

# Integration тесты (требуется запущенная БД)
go test ./tests/integration/... -v

# Нагрузочное тестирование (k6)
k6 run tests/load/load_test.js
```

### Code quality

```bash
# Форматирование
go fmt ./...

# Линтер
go vet ./...

# Статический анализ (требуется golangci-lint)
golangci-lint run

# Обновление зависимостей
go mod tidy
go mod verify
```

## Нагрузочное тестирование

Проект включает k6 тесты для проверки производительности:

```bash
k6 run tests/load/load_test.js
```

### Последние результаты (5 минут, до 2000 VUs)

- **RPS**: 647 req/s sustained
- **Latency P95**: 6.2ms (threshold: <800ms) ✅
- **Error rate**: 0.00% (threshold: <1%) ✅
- **Success rate**: 100%

Тесты симулируют смешанную нагрузку:
- 35% - поиск рейсов
- 20% - получение рейса по ID
- 15% - поиск аэропортов
- 15% - статистика маршрутов
- 8% - схема посадочных мест
- 7% - информация о самолете

## База данных

### Схема

База использует схему `bookings` с следующими таблицами:

- `airplanes_data` - типы самолетов
- `airports_data` - аэропорты
- `routes` - маршруты
- `flights` - рейсы
- `seats` - посадочные места
- `bookings` - бронирования
- `tickets` - билеты
- `segments` - сегменты перелета
- `boarding_passes` - посадочные талоны

### Многоязычность

Названия хранятся в JSONB:
```sql
-- Пример: airplanes_data.model
{"en": "Boeing 777-300", "ru": "Боинг 777-300"}

-- Установка языка
SET bookings.lang = 'ru';
```

### Backup и восстановление

```bash
# Backup
docker-compose exec postgres pg_dump -U airops airops > backup.sql

# Restore
docker-compose exec -T postgres psql -U airops airops < backup.sql
```

## Monitoring

### Prometheus метрики

API экспортирует метрики в формате Prometheus на `/metrics`:

- `http_requests_total` - общее количество запросов
- `http_request_duration_seconds` - длительность запросов
- `http_requests_in_flight` - текущие запросы

### Logging

Логи выводятся в JSON формате для удобной агрегации:

```json
{
  "time": "2025-01-11T12:34:56Z",
  "level": "info",
  "msg": "request completed",
  "request_id": "abc123",
  "method": "GET",
  "path": "/flights/12345",
  "status": 200,
  "duration_ms": 15.3,
  "ip": "192.168.1.1"
}
```

## Production deployment

### Рекомендации

1. **Используйте HTTPS** с reverse proxy (nginx/caddy)
2. **Настройте firewall** - открыт только порт API
3. **Измените API_KEY** на сложный
4. **Настройте логирование** в файл или syslog
5. **Мониторинг** через Prometheus/Grafana
6. **Backup БД** регулярный и автоматический
7. **Resource limits** в docker-compose для production

### Пример nginx конфига

```nginx
server {
    listen 443 ssl http2;
    server_name api.example.com;
    
    ssl_certificate /etc/ssl/certs/cert.pem;
    ssl_certificate_key /etc/ssl/private/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Troubleshooting

### API не стартует

```bash
# Проверить логи
docker-compose logs api

# Проверить подключение к БД
docker-compose exec postgres psql -U airops -d airops -c "SELECT 1"

# Проверить порт
lsof -i :8080
```

### Медленные запросы

```sql
-- Включить логирование медленных запросов в PostgreSQL
ALTER DATABASE airops SET log_min_duration_statement = 1000;

-- Проверить индексы
SELECT schemaname, tablename, indexname
FROM pg_indexes
WHERE schemaname = 'bookings';
```

### High CPU usage

- Проверить connection pool limits в конфиге
- Проверить N+1 запросы в логах
- Использовать `pprof` для профилирования: `go tool pprof http://localhost:8080/debug/pprof/profile`

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

MIT License - see LICENSE file for details

## Контакты

- **Issues**: GitHub Issues
- **Документация**: `/swagger/index.html`
- **API Status**: `/health`