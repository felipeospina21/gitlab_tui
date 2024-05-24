package tui

import (
	"gitlab_tui/internal/style"

	"github.com/charmbracelet/bubbles/list"
)

type item struct {
	title, id string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.id }
func (i item) FilterValue() string { return i.title }

type ProjectsModel struct {
	List      list.Model
	ProjectID string
}

func InitProjectsList() list.Model {
	items := []list.Item{
		item{title: "Spellbook", id: "17050"},
		item{title: "Radar", id: "98211"},
		item{title: "Planning-Tool", id: "58799"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 362, 50)
	l.Title = "Disney Projects"
	l.Styles.Title = style.TitleStyle
	l.Styles.PaginationStyle = style.PaginationStyle
	l.Styles.HelpStyle = style.HelpStyle

	return l
}
