package router

import "net/http"

const (
	MAX_TRIE_CHILD_NODES   = 26
)

type MethodTree struct {
	childTree 		map[rune]*MethodTree
	handlerChain 	[]http.HandlerFunc
}

func NewMethodTree() *MethodTree {
	root := new(MethodTree)
	root.childTree = make(map[rune]*MethodTree, MAX_TRIE_CHILD_NODES)
	return root
}

func (m *MethodTree) addRuleToMethodTree(path string, handlers []http.HandlerFunc) {
	for _, v := range path {
		if m.childTree[v] == nil {
			node := new(MethodTree)
			node.childTree = make(map[rune]*MethodTree, MAX_TRIE_CHILD_NODES)
			m.childTree[v] = node
		}
		m = m.childTree[v]
	}
	for _, v := range handlers {
		m.handlerChain = append(m.handlerChain, v)
	}
}

func (m *MethodTree) getHandlersByPath(path string) (handlers []http.HandlerFunc) {
	for _, v := range path {
		if m.childTree[v] == nil {
			return
		}
		m = m.childTree[v]
	}
	return m.handlerChain
}