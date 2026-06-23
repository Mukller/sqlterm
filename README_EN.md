# sqlterm

Interactive TUI browser for SQLite databases. Open any .db file and work with it in the terminal — no heavy GUI tools needed.

## Controls

| Key | Action |
|-----|--------|
| j/k | Navigate tables |
| Enter | Select table |
| / | SQL query mode |
| s | Show schema |
| n/p | Next/prev page |
| q | Quit |

## Usage

```bash
go install github.com/Mukller/sqlterm@latest
sqlterm mydb.sqlite
```
