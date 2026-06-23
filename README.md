# sqlterm

Интерактивный TUI-браузер для SQLite баз данных. Открывай любой `.db` файл и работай с ним прямо в терминале — без установки тяжёлых GUI-инструментов.

## Что умеет

- Список таблиц в левой панели
- Просмотр данных таблицы (постраничный)
- Произвольные SQL-запросы в строке ввода
- Показ схемы таблицы (CREATE TABLE)
- Навигация клавишами

## Управление

| Клавиша | Действие |
|---------|----------|
| `↑/↓`, `j/k` | Навигация по таблицам |
| `Enter` | Выбрать таблицу |
| `/` | Режим SQL-запроса |
| `s` | Показать схему таблицы |
| `n/p` | Следующая/предыдущая страница |
| `q`, `Ctrl+C` | Выход |

## Установка

```bash
go install github.com/you/sqlterm@latest
```

Или собрать из исходников:

```bash
go build -o sqlterm .
```

## Запуск

```bash
sqlterm mydb.sqlite
sqlterm chinook.db
```

## Зависимости

- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI фреймворк
- [lipgloss](https://github.com/charmbracelet/lipgloss) — стили
- [go-sqlite3](https://github.com/mattn/go-sqlite3) — SQLite драйвер (cgo)
