package models

type TermTrie struct {
	NodeId      int
	Letter      string
	ParentId    int
	IsEndOfWord bool
}
