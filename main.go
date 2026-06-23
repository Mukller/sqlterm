package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
)

const pageSize = 20

// ─── стили ────────────────────────────────────────────────────────────────────

var (
	styleTitle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	styleSelected = lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("39"))
	styleBorder   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("238"))
	styleHeader   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("243"))
	styleError    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	styleHelp     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	styleDim      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// ─── модель ───────────────────────────────────────────────────────────────────

type mode int

const (
	modeTable mode = iota
	modeQuery
	modeSchema
)

type model struct {
	db        *sql.DB
	dbPath    string
	tables    []string
	cursor    int
	selTable  string
	cols      []string
	rows      [][]string
	page      int
	totalRows int
	mode      mode
	input     string
	err       string
	schema    string
	width     int
	height    int
}

func newModel(path string) (model, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return model{}, err
	}
	m := model{db: db, dbPath: path}
	if err := m.loadTables(); err != nil {
		return model{}, err
	}
	if len(m.tables) > 0 {
		m.selTable = m.tables[0]
		m.loadTable(m.tables[0], 0)
	}
	return m, nil
}

func (m *model) loadTables() error {
	rows, err := m.db.Query(`SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`)
	if err != nil {
		return err
	}
	defer rows.Close()
	m.tables = nil
	for rows.Next() {
		var name string
		rows.Scan(&name)
		m.tables = append(m.tables, name)
	}
	return nil
}

func (m *model) loadTable(table string, page int) {
	m.err = ""
	m.mode = modeTable
	m.selTable = table
	m.page = page

	countRow := m.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %q", table))
	countRow.Scan(&m.totalRows)

	rows, err := m.db.Query(
		fmt.Sprintf("SELECT * FROM %q LIMIT %d OFFSET %d", table, pageSize, page*pageSize),
	)
	if err != nil {
		m.err = err.Error()
		return
	}
	defer rows.Close()

	m.cols, _ = rows.Columns()
	m.rows = nil
	for rows.Next() {
		vals := make([]interface{}, len(m.cols))
		ptrs := make([]interface{}, len(m.cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		rows.Scan(ptrs...)
		row := make([]string, len(m.cols))
		for i, v := range vals {
			if v == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		m.rows = append(m.rows, row)
	}
}

func (m *model) runQuery(q string) {
	m.err = ""
	m.mode = modeTable
	m.selTable = ""
	m.page = 0

	rows, err := m.db.Query(q)
	if err != nil {
		m.err = err.Error()
		m.mode = modeTable
		return
	}
	defer rows.Close()

	m.cols, _ = rows.Columns()
	m.totalRows = 0
	m.rows = nil
	for rows.Next() {
		vals := make([]interface{}, len(m.cols))
		ptrs := make([]interface{}, len(m.cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		rows.Scan(ptrs...)
		row := make([]string, len(m.cols))
		for i, v := range vals {
			if v == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		m.rows = append(m.rows, row)
		m.totalRows++
	}
}

func (m *model) loadSchema(table string) {
	row := m.db.QueryRow(
		"SELECT sql FROM sqlite_master WHERE type='table' AND name=?", table,
	)
	var s string
	row.Scan(&s)
	m.schema = s
	m.mode = modeSchema
}

// ─── bubbletea ────────────────────────────────────────────────────────────────

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch m.mode {
		case modeQuery:
			return m.updateQuery(msg)
		case modeSchema:
			m.mode = modeTable
			return m, nil
		default:
			return m.updateNormal(msg)
		}
	}
	return m, nil
}

func (m model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.db.Close()
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.loadTable(m.tables[m.cursor], 0)
		}
	case "down", "j":
		if m.cursor < len(m.tables)-1 {
			m.cursor++
			m.loadTable(m.tables[m.cursor], 0)
		}
	case "enter":
		if len(m.tables) > 0 {
			m.loadTable(m.tables[m.cursor], 0)
		}
	case "n":
		if (m.page+1)*pageSize < m.totalRows {
			m.loadTable(m.selTable, m.page+1)
		}
	case "p":
		if m.page > 0 {
			m.loadTable(m.selTable, m.page-1)
		}
	case "/":
		m.mode = modeQuery
		m.input = ""
	case "s":
		if len(m.tables) > 0 {
			m.loadSchema(m.tables[m.cursor])
		}
	case "r":
		m.loadTables()
	}
	return m, nil
}

func (m model) updateQuery(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.input != "" {
			m.runQuery(m.input)
		}
		m.mode = modeTable
	case "esc":
		m.mode = modeTable
	case "backspace", "ctrl+h":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	default:
		if len(msg.Runes) == 1 {
			m.input += string(msg.Runes)
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Загрузка..."
	}

	leftW := 22
	rightW := m.width - leftW - 3

	// левая панель — список таблиц
	left := renderLeft(m, leftW)
	// правая панель
	right := renderRight(m, rightW)

	body := lipgloss.JoinHorizontal(lipgloss.Top,
		styleBorder.Width(leftW).Render(left),
		"  ",
		styleBorder.Width(rightW).Render(right),
	)

	help := renderHelp(m)
	return lipgloss.JoinVertical(lipgloss.Left, body, help)
}

func renderLeft(m model, w int) string {
	title := styleTitle.Render("  " + m.dbPath)
	lines := []string{title, ""}
	for i, t := range m.tables {
		line := "  " + t
		if i == m.cursor {
			line = styleSelected.Width(w).Render(line)
		}
		lines = append(lines, line)
	}
	if len(m.tables) == 0 {
		lines = append(lines, styleDim.Render("  (нет таблиц)"))
	}
	return strings.Join(lines, "\n")
}

func renderRight(m model, w int) string {
	if m.mode == modeQuery {
		return styleTitle.Render("SQL> ") + m.input + "█"
	}
	if m.mode == modeSchema {
		return styleTitle.Render("Схема\n\n") + m.schema
	}
	if m.err != "" {
		return styleError.Render("Ошибка: " + m.err)
	}
	if len(m.cols) == 0 {
		return styleDim.Render("Пусто")
	}

	// вычислим ширину каждой колонки
	colW := make([]int, len(m.cols))
	for i, c := range m.cols {
		colW[i] = len(c)
	}
	for _, row := range m.rows {
		for i, cell := range row {
			if len(cell) > colW[i] {
				colW[i] = len(cell)
			}
			if colW[i] > 30 {
				colW[i] = 30
			}
		}
	}

	pad := func(s string, n int) string {
		if len(s) > n {
			return s[:n-1] + "…"
		}
		return s + strings.Repeat(" ", n-len(s))
	}

	var b strings.Builder

	// заголовок
	header := make([]string, len(m.cols))
	for i, c := range m.cols {
		header[i] = styleHeader.Render(pad(c, colW[i]))
	}
	b.WriteString(strings.Join(header, "  "))
	b.WriteByte('\n')
	b.WriteString(strings.Repeat("─", min(w-4, 80)))
	b.WriteByte('\n')

	// строки
	for _, row := range m.rows {
		cells := make([]string, len(row))
		for i, cell := range row {
			cells[i] = pad(cell, colW[i])
		}
		b.WriteString(strings.Join(cells, "  "))
		b.WriteByte('\n')
	}

	// пагинация
	if m.totalRows > 0 {
		from := m.page*pageSize + 1
		to := from + len(m.rows) - 1
		b.WriteString(fmt.Sprintf("\n%s", styleDim.Render(
			fmt.Sprintf("строки %d–%d из %d", from, to, m.totalRows),
		)))
	}

	return b.String()
}

func renderHelp(m model) string {
	if m.mode == modeQuery {
		return styleHelp.Render("Enter: выполнить  Esc: отмена")
	}
	return styleHelp.Render("↑↓/jk: навигация  Enter: выбрать  /: SQL  s: схема  n/p: страницы  q: выход")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ─── main ─────────────────────────────────────────────────────────────────────

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "использование: sqlterm <файл.db>")
		os.Exit(1)
	}

	m, err := newModel(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка открытия БД:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
