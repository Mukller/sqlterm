[Русский](README.md)

# sqlterm

Interactive TUI browser for SQLite. Open any `.db` file and work with it
directly in the terminal — no heavy GUI tools needed.

## Run

```bash
sqlterm mydb.sqlite
sqlterm chinook.db
```

Build from source:
```bash
go build -o sqlterm .
```

## Controls

| Key | Action |
|-----|--------|
| `↑/↓`, `j/k` | Navigate tables |
| `Enter` | Select table |
| `/` | Enter SQL query |
| `s` | Show table schema |
| `n/p` | Next/previous page |
| `q`, `Ctrl+C` | Quit |

## Features

- Table list in left panel
- Paginated data view
- Arbitrary SQL queries
- Table schema (`CREATE TABLE`)

## Dependencies

- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) — styling
- [go-sqlite3](https://github.com/mattn/go-sqlite3) — SQLite driver (cgo)
