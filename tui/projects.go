package tui

import (
	"gitlab_tui/internal/server"
	"gitlab_tui/internal/style"
	tbl "gitlab_tui/tui/components/table"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	title, id string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.id }
func (i item) FilterValue() string { return i.title }

//
// type itemDelegate struct{}
//
// func (d itemDelegate) Height() int                             { return 1 }
// func (d itemDelegate) Spacing() int                            { return 0 }
// func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
// func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
// 	var title string
//
// 	if i, ok := listItem.(item); ok {
// 		title = i.Title()
// 		// desc = i.Description()
// 	} else {
// 		return
// 	}
//
// 	str := fmt.Sprintf("%s", title)
//
// 	fn := style.ListItemStyle.Render
// 	if index == m.Index() {
// 		fn = func(s ...string) string {
// 			return style.ListSelectedItemStyle.Render(strings.Join(s, " "))
// 		}
// 	}
//
// 	fmt.Fprint(w, fn(str))
// }

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

	// l := list.New(items, itemDelegate{}, 362, 50)
	l := list.New(items, list.NewDefaultDelegate(), 362, 50)
	l.Title = "Disney Projects"
	l.Styles.Title = style.ListItemStyle
	l.Styles.PaginationStyle = style.ListPaginationStyle
	l.Styles.HelpStyle = style.ListHelpStyle

	return l
}

// View cmds
func (m *Model) viewMergeReqs(window tea.WindowSizeMsg) tea.Cmd {
	s := m.Projects.List.SelectedItem()
	i, ok := s.(item)
	var c tea.Cmd
	if ok {
		m.Projects.ProjectID = i.id
		r, err := server.GetMergeRequests(m.Projects.ProjectID)
		c = func() tea.Msg {
			if err != nil {
				return err
			}

			return "success_mergeReqs"
		}
		m.MergeRequests.List = tbl.InitModel(tbl.InitModelParams{
			Rows:      r,
			Colums:    tbl.GetMergeReqsColums(window.Width - 10),
			StyleFunc: tbl.StyleIconsColumns(style.Table(), tbl.MergeReqsIconCols),
		})
	}
	return c
}
