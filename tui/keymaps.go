package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type GlobalKeyMap struct {
	Help         key.Binding
	Quit         key.Binding
	ThrowError   key.Binding
	NextTab      key.Binding
	PrevTab      key.Binding
	NextPage     key.Binding
	PrevPage     key.Binding
	NavigateBack key.Binding
}

func (k GlobalKeyMap) commonKeys() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.NavigateBack, k.NextTab, k.PrevTab, k.NextPage, k.PrevPage}
}

func (k GlobalKeyMap) ShortHelp() []key.Binding {
	return k.commonKeys()
}

func (k GlobalKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.commonKeys(),
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
	NextTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev tab"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("->", "next page"),
	),
	PrevPage: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("<-", "prev page"),
	),
	NavigateBack: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "navigate back"),
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
	Description   key.Binding
	Merge         key.Binding
	GlobalKeyMap
}

func (k MergeReqsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Comments, k.Pipelines, k.Description, k.OpenInBrowser, k.Merge, k.NextTab, k.PrevTab, k.NavigateBack, k.Help, k.Quit}
}

func (k MergeReqsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.commonKeys(), // first column
		{k.Comments, k.Pipelines, k.Description, k.OpenInBrowser, k.Merge, k.Refetch}, // second column
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
	GlobalKeyMap
}

func (k CommentsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Description, k.OpenInBrowser, k.Refetch, k.NextTab, k.PrevTab, k.NavigateBack, k.Help, k.Quit}
}

func (k CommentsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.commonKeys(),
		{k.Description, k.OpenInBrowser, k.Refetch},
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
	GlobalKeyMap
}

func (k PipelineKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Jobs, k.OpenInBrowser, k.Refetch, k.NextTab, k.PrevTab, k.NavigateBack, k.Help, k.Quit}
}

func (k PipelineKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.commonKeys(),
		{k.Jobs, k.OpenInBrowser, k.Refetch},
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
	Jobs: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view jobs"),
	),
	GlobalKeyMap: GlobalKeys,
}

type JobsKeyMap struct {
	Refetch       key.Binding
	OpenInBrowser key.Binding
	GlobalKeyMap
}

func (k JobsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.OpenInBrowser, k.Refetch, k.NextTab, k.PrevTab, k.NavigateBack, k.Help, k.Quit}
}

func (k JobsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.commonKeys(),
		{k.OpenInBrowser, k.Refetch},
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
	GlobalKeyMap: GlobalKeys,
}

type ProjectsKeyMap struct {
	ViewMRs key.Binding
	GlobalKeyMap
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
}

type MdKeyMap struct {
	NavigateBack key.Binding
	GlobalKeyMap
}

func (k MdKeyMap) ShortHelp() []key.Binding {
	return k.commonKeys()
}

func (k MdKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.commonKeys(),
	}
}

var MdKeys = MdKeyMap{
	GlobalKeyMap: GlobalKeys,
}
