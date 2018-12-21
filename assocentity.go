package assocentity

import (
	"bytes"
	"fmt"
)

type Filter struct {
	blacklist []string
	whitelist []string
	entity    string
	separator []byte
}

func NewFilter(blacklist, whitelist []string, entity string, separator []byte) *Filter {
	return &Filter{
		blacklist: blacklist,
		whitelist: whitelist,
		entity:    entity,
		separator: separator,
	}
}

func TraversableLatin(text []byte, f *Filter) [][]byte {
	var bl = make([]byte, 0)
	for _, ru := range text {
		ok := true
		for _, r := range f.blacklist {
			if string(ru) == r {
				ok = false
			}
		}

		if ok {
			bl = append(bl, byte(ru))
		}
	}

	spl := bytes.Split(bl, f.separator)

	var sw = make([][]byte, 0)
	for _, b := range spl {
		ok := true
		for _, str := range f.blacklist {
			if string(b) == str {
				ok = false
			}
		}

		if ok {
			sw = append(sw, b)
		}
	}

	return sw
}

func Graph(f *Filter, traversable [][]byte) bool {
	entityb := []byte(f.entity)

	for i, by := range traversable {
		ok, last := index(i, traversable[i:], entityb, f.separator)
		if ok {
			fmt.Println(i)
			fmt.Println(last)
		}

		if f.entity == string(by) {
			leftwalker(i-1, traversable)
			rightwalker(i+1, traversable)
		}
	}

	return true
}

func leftwalker(index int, traversable [][]byte) {
	for i := index; i >= 0; i-- {
		fmt.Println(string(traversable[i]))
		fmt.Println()
	}
}

func rightwalker(index int, traversable [][]byte) {
	for i := index; i < len(traversable); i++ {
		fmt.Println(string(traversable[i]))
		fmt.Println()
	}
}

func index(index int, traversable [][]byte, entityb []byte, separator []byte) (bool, int) {
	t := make([][]byte, len(traversable) + 1)
	copy(t, traversable)
	t = append(t, separator)

	flatt := reset(t, separator)

	ok := true
	i := len(entityb)
	for i = range entityb {
		if len(entityb) <= len(flatt) && entityb[i] != flatt[i] {
			ok = false
		}
	}

	return ok, i
}

func reset(traversable [][]byte, separator []byte) []byte {
	f := make([]byte, 0)
	for _, t := range traversable {
		f = append(f, t...)
		f = append(f, separator...)
	}

	return f
}

//func Graph(f *Filter, text []byte) bool {
//	t := measure(text, f)
//	//fmt.Println(t)
//
//	gr := make([]*simple.DirectedGraph, 0)
//	for i := range f.entities {
//		g := simple.NewDirectedGraph()
//		g.AddNode(simple.Node(i))
//
//		max := len(f.entities) + 1
//		for range t {
//			g.AddNode(simple.Node(max))
//			max++
//		}
//
//		gr = append(gr, g)
//	}
//
//	for _, g := range gr {
//		result, _ := dot.Marshal(g, "", "", "  ")
//		fmt.Print(string(result))
//	}
//
//	return true
//}

func printlnBytes(by [][]byte) {
	for _, b := range by {
		fmt.Println(string(b))
	}
}
