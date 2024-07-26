package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type GlobalKeyMap struct {
	Help       key.Binding
	Quit       key.Binding
	ThrowError key.Binding
}

func (k GlobalKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k GlobalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{},               // first column
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

	// TODO: make this available only when program is run whith certain cmd
	ThrowError: key.NewBinding(
		key.WithKeys("E"),
		key.WithHelp("E", "throw error"),
	),
}

type MergeReqsKeyMap struct {
	Comments      key.Binding
	Pipelines     key.Binding
	OpenInBrowser key.Binding
	Refetch       key.Binding
	NavigateBack  key.Binding
	Description   key.Binding
	Merge         key.Binding
	GlobalKeyMap
}

func (k MergeReqsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Comments, k.Pipelines, k.Description, k.Refetch, k.OpenInBrowser, k.NavigateBack, k.Merge}
}

func (k MergeReqsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Comments, k.Pipelines, k.Description, k.Refetch, k.OpenInBrowser, k.NavigateBack}, // first column
		{k.Merge, k.Help, k.Quit}, // second column
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
	Merge: key.NewBinding(
		key.WithKeys("M"),
		key.WithHelp("M", "merge MR"),
	),
	GlobalKeyMap: GlobalKeys,
}

type CommentsKeyMap struct {
	Refetch       key.Binding
	OpenInBrowser key.Binding
	Description   key.Binding
	NavigateBack  key.Binding
	GlobalKeyMap
}

func (k CommentsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Description, k.NavigateBack, k.OpenInBrowser, k.Refetch}
}

func (k CommentsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Description, k.NavigateBack, k.OpenInBrowser, k.Refetch}, // first column
		{k.Help, k.Quit}, // second column
	}
}

var CommentsKeys = CommentsKeyMap{
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

type PipelineKeyMap struct {
	Jobs          key.Binding
	Refetch       key.Binding
	OpenInBrowser key.Binding
	NavigateBack  key.Binding
	GlobalKeyMap
}

func (k PipelineKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Jobs, k.NavigateBack, k.OpenInBrowser, k.Refetch}
}

func (k PipelineKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Jobs, k.NavigateBack, k.OpenInBrowser, k.Refetch}, // first column
		{k.Help, k.Quit}, // second column
	}
}

var PipelineKeys = PipelineKeyMap{
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
	Jobs: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view jobs"),
	),
	GlobalKeyMap: GlobalKeys,
}

type JobsKeyMap struct {
	Refetch       key.Binding
	OpenInBrowser key.Binding
	NavigateBack  key.Binding
	GlobalKeyMap
}

func (k JobsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.NavigateBack, k.OpenInBrowser, k.Refetch}
}

func (k JobsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NavigateBack, k.OpenInBrowser, k.Refetch}, // first column
		{k.Help, k.Quit}, // second column
	}
}

var JobsKeys = JobsKeyMap{
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
	GlobalKeyMap: GlobalKeys,
}

type ProjectsKeyMap struct {
	ViewMRs key.Binding
	GlobalKeyMap
	// list.KeyMap
}

func (k ProjectsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ViewMRs}
}

func (k ProjectsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ViewMRs}, // first column
		// {k.ShowFullHelp, k.Quit, k.Filter, k.ReloadConfig}, // second column
	}
}

var ProjectsKeys = ProjectsKeyMap{
	ViewMRs: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view merge requests"),
	),
	GlobalKeyMap: GlobalKeys,
	// KeyMap: list.DefaultKeyMap(),
}

type MdKeyMap struct {
	NavigateBack key.Binding
	GlobalKeyMap
}

func (k MdKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NavigateBack}
}

func (k MdKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NavigateBack}, // first column
		// {k.ShowFullHelp, k.Quit, k.Filter, k.ReloadConfig}, // second column
	}
}

var MdKeys = MdKeyMap{
	NavigateBack: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "navigate back"),
	),
	GlobalKeyMap: GlobalKeys,
}
