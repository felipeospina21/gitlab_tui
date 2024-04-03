package customlog

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Logs to debug.log file
//
//	logToFile("log", func() {
//		log.Println(strconv.Itoa(msg.Width))
//		log.Println("tw " + strconv.Itoa(m.table.Width()))
//	})
func ToFile(logPrefix string, cb func()) {
	f, err := tea.LogToFile("debug.log", logPrefix)
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	cb()
}
