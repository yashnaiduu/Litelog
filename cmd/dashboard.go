package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/yashnaiduu/Litelog/models"
	"github.com/yashnaiduu/Litelog/storage"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			MarginBottom(1)

	statNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A8A8A8")).
			Width(20)

	statValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87"))

	fadedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

type DashboardStats struct {
	LogsPerSec     int
	TotalLogs      int
	ActiveServices int
	ErrorsPerMin   int
	TopServices    []ServiceStat
	RecentErrors   []models.LogEntry
}

type ServiceStat struct {
	Name  string
	Count int
}

type tickMsg time.Time

type model struct {
	stats DashboardStats
	err   error
	store *storage.Store
}

func fetchStats(store *storage.Store) DashboardStats {
	var s DashboardStats
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	store.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM logs").Scan(&s.TotalLogs)

	var logsLast5 int
	store.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM logs WHERE timestamp >= datetime('now', '-5 seconds')").Scan(&logsLast5)
	s.LogsPerSec = logsLast5 / 5

	store.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM logs WHERE level = 'ERROR' AND timestamp >= datetime('now', '-1 minute')").Scan(&s.ErrorsPerMin)

	store.DB.QueryRowContext(ctx, "SELECT COUNT(DISTINCT service) FROM logs WHERE timestamp >= datetime('now', '-5 minutes')").Scan(&s.ActiveServices)

	rows, err := store.DB.QueryContext(ctx, "SELECT service, COUNT(*) as c FROM logs GROUP BY service ORDER BY c DESC LIMIT 5")
	if err == nil {
		for rows.Next() {
			var ss ServiceStat
			rows.Scan(&ss.Name, &ss.Count)
			s.TopServices = append(s.TopServices, ss)
		}
		rows.Close()
	}

	rows, err = store.DB.QueryContext(ctx, "SELECT id, timestamp, level, service, message FROM logs WHERE level = 'ERROR' ORDER BY timestamp DESC LIMIT 5")
	if err == nil {
		for rows.Next() {
			var e models.LogEntry
			var ts string
			rows.Scan(&e.ID, &ts, &e.Level, &e.Service, &e.Message)

			if parsedTs, err := time.Parse(time.RFC3339, ts); err == nil {
				ts = parsedTs.Format("15:04:05")
			} else if parsedTs, err := time.Parse("2006-01-02 15:04:05", ts); err == nil {
				ts = parsedTs.Format("15:04:05")
			}
			s.RecentErrors = append(s.RecentErrors, models.LogEntry{
				Service: ts + " " + e.Service,
				Message: e.Message,
			})
		}
		rows.Close()
	}

	return s
}

func initialModel(store *storage.Store) model {
	return model{
		stats: fetchStats(store),
		store: store,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tickMsg:
		m.stats = fetchStats(m.store)
		return m, tickCmd()
	}
	return m, nil
}

func buildDivider(width int) string {
	return fadedStyle.Render(strings.Repeat("─", width)) + "\n"
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	var view strings.Builder

	view.WriteString(titleStyle.Render("LiteLog Dashboard"))
	view.WriteString("\n")
	view.WriteString(buildDivider(36))

	renderStat := func(name string, val int) {
		view.WriteString(statNameStyle.Render(name))
		view.WriteString(statValueStyle.Render(fmt.Sprintf("%d", val)))
		view.WriteString("\n")
	}

	renderStat("Logs/sec:", m.stats.LogsPerSec)
	renderStat("Total Logs:", m.stats.TotalLogs)
	renderStat("Active Services:", m.stats.ActiveServices)
	renderStat("Errors/min:", m.stats.ErrorsPerMin)

	view.WriteString("\n")
	view.WriteString(headerStyle.Render("Top Services"))
	view.WriteString("\n")
	view.WriteString(buildDivider(24))

	for _, s := range m.stats.TopServices {
		view.WriteString(fmt.Sprintf("%-20s %d logs\n", s.Name, s.Count))
	}
	if len(m.stats.TopServices) == 0 {
		view.WriteString(fadedStyle.Render("No data\n"))
	}

	view.WriteString("\n")
	view.WriteString(headerStyle.Render("Recent Errors"))
	view.WriteString("\n")
	view.WriteString(buildDivider(24))

	for _, e := range m.stats.RecentErrors {
		view.WriteString(errorStyle.Render(fmt.Sprintf("%s   %s\n", e.Service, e.Message)))
	}
	if len(m.stats.RecentErrors) == 0 {
		view.WriteString(fadedStyle.Render("No errors logged.\n"))
	}

	view.WriteString("\n" + fadedStyle.Render("Press q to quit"))

	return lipgloss.NewStyle().Margin(1, 2).Render(view.String())
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Open live terminal dashboard analytics",
	Run: func(cmd *cobra.Command, args []string) {
		store, err := storage.InitDB(dbPath)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		p := tea.NewProgram(initialModel(store), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatalf("Error running dashboard: %v", err)
		}
	},
}

func init() {
	dashboardCmd.Flags().StringVar(&dbPath, "db", "litelog.db", "Path to SQLite database")
	rootCmd.AddCommand(dashboardCmd)
}
