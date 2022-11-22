package dag

import "sort"

type Node int

type NodeSorter []Node

func (ns NodeSorter) Len() int {
	return len(ns)
}

func (ns NodeSorter) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

func (ns NodeSorter) Less(i, j int) bool {
	return ns[i] < ns[j]
}

type Edge struct {
	Depender Node
	Dependee Node
}

type EdgeSorter []Edge

func (es EdgeSorter) Len() int {
	return len(es)
}

func (es EdgeSorter) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es EdgeSorter) Less(i, j int) bool {
	return es[i].Depender < es[j].Depender || (es[i].Depender == es[j].Depender && es[i].Dependee < es[j].Dependee)
}

type Graph struct {
	Nodes []Node
	Edges []Edge
}

func (g *Graph) HasNode(n Node) bool {
	for _, node := range g.Nodes {
		if node == n {
			return true
		}
	}
	return false
}

func (g *Graph) AddNode(n Node) {
	if !g.HasNode(n) {
		g.Nodes = append(g.Nodes, n)
		sort.Sort(NodeSorter(g.Nodes))
	}
}

func (g *Graph) HasEdge(e Edge) bool {
	for _, edge := range g.Edges {
		if e.Depender == edge.Depender && e.Dependee == edge.Dependee {
			return true
		}
	}
	return false
}

func (g *Graph) AddEdge(e Edge) {
	if !g.HasEdge(e) {
		g.Edges = append(g.Edges, e)
		sort.Sort(EdgeSorter(g.Edges))
		g.AddNode(e.Depender)
		g.AddNode(e.Dependee)
	}
}

func (g *Graph) ConnectedComponentRoots() []Node {
	nonRoot := make(map[Node]struct{}, 0)
	for _, edge := range g.Edges {
		nonRoot[edge.Depender] = struct{}{}
	}
	var root []Node
	for _, node := range g.Nodes {
		_, ok := nonRoot[node]
		if !ok {
			root = append(root, node)
		}
	}
	return root
}

func (g *Graph) DirectDependers(dependee Node) []Node {
	var dependers []Node
	for _, edge := range g.Edges {
		if edge.Dependee == dependee {
			appendIfUnique(&dependers, edge.Depender)
		}
	}
	sort.Sort(NodeSorter(dependers))
	return dependers
}

func (g *Graph) DirectDependees(depender Node) []Node {
	var dependees []Node
	for _, edge := range g.Edges {
		if edge.Depender == depender {
			appendIfUnique(&dependees, edge.Dependee)
		}
	}
	sort.Sort(NodeSorter(dependees))
	return dependees
}

func appendIfUnique(nodes *[]Node, n Node) {
	for _, node := range *nodes {
		if node == n {
			return
		}
	}
	*nodes = append(*nodes, n)
	return
}
