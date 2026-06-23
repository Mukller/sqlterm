# Contributing

## Как помочь

Принимаю PR с:
- Улучшением навигации (mouse support)
- Новыми режимами просмотра (hex для BLOB, форматирование JSON)
- Экспортом результатов (CSV, JSON)
- Поддержкой PostgreSQL/MySQL через драйверы

## Как делать

```bash
git clone https://github.com/Mukller/sqlterm
cd sqlterm
go mod tidy
go build .
./sqlterm test.db
```

## Стиль

- `gofmt` обязательно
- TUI логика — через bubbletea, без прямых terminal escape codes
- Без glue-кода: каждый режим (таблица, схема, запрос) — отдельная модель
