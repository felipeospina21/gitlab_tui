package tui

import "github.com/charmbracelet/bubbles/key"

type GlobalKeyMap struct {
	Help key.Binding
	Quit key.Binding
}

func (k GlobalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k GlobalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// {k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit}, // second column
	}
}

var GlobalKeys = GlobalKeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type MergeReqsKeyMap struct {
	Comments      key.Binding
	Pipelines     key.Binding
	OpenInBrowser key.Binding
	Refetch       key.Binding
	NavigateBack  key.Binding
	Description   key.Binding
	GlobalKeyMap
}

func (k MergeReqsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Comments, k.Pipelines, k.Description, k.Refetch, k.OpenInBrowser, k.NavigateBack}
}

func (k MergeReqsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Comments, k.Pipelines, k.Description, k.Refetch, k.OpenInBrowser, k.NavigateBack}, // first column
		{k.Help, k.Quit}, // second column
	}
}

var MergeReqsKeys = MergeReqsKeyMap{
	Comments: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "view comments"),
	),
	Pipelines: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "view pipelines"),
	),
	OpenInBrowser: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "open in browser"),
	),
	Refetch: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refetch"),
	),
	NavigateBack: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "navigate back"),
	),
	Description: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view description"),
	),
	GlobalKeyMap: GlobalKeys,
}

type CommentsKeyMap struct {
	GlobalKeyMap
}

func (k CommentsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k CommentsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// {k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit}, // second column
	}
}

var CommentsKeys = CommentsKeyMap{
	GlobalKeyMap: GlobalKeys,
}

type PipelineKeyMap struct {
	GlobalKeyMap
}

func (k PipelineKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k PipelineKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// {k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit}, // second column
	}
}

var PipelinKeys = PipelineKeyMap{
	GlobalKeyMap: GlobalKeys,
}
