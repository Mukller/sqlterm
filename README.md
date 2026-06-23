[English](README_EN.md)

# sqlterm

Интерактивный TUI-браузер для SQLite. Открывай любой `.db` файл и работай с ним
прямо в терминале — без тяжёлых GUI-инструментов.

## Запуск

```bash
sqlterm mydb.sqlite
sqlterm chinook.db
```

Собрать из исходников:
```bash
go build -o sqlterm .
```

## Управление

| Клавиша | Действие |
|---------|----------|
| `↑/↓`, `j/k` | Навигация по таблицам |
| `Enter` | Выбрать таблицу |
| `/` | Ввести SQL-запрос |
| `s` | Показать схему таблицы |
| `n/p` | Следующая/предыдущая страница |
| `q`, `Ctrl+C` | Выход |

## Что умеет

- Список таблиц в левой панели
- Просмотр данных (постраничный)
- Произвольные SQL-запросы
- Схема таблицы (`CREATE TABLE`)

## Зависимости

- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI
- [lipgloss](https://github.com/charmbracelet/lipgloss) — стили
- [go-sqlite3](https://github.com/mattn/go-sqlite3) — SQLite (cgo)
