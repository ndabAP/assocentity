package assocentity

import (
	"bytes"
	"fmt"
	"strings"
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
	entityb := []byte(strings.Replace(f.entity, string(f.separator), "", -1))

	runes, flatt := 0, flatten(traversable)
	for i, by := range traversable {
		if len(flatt) == runes {
			break
		}

		if len(flatt[runes:]) >= len(entityb) &&
			len(flatt[runes:runes+len(entityb)]) >= len(entityb)-1 &&
			bytes.Equal(flatt[runes:runes+len(entityb)], entityb) {
			leftwalker(i-1, traversable)
			//rightwalker(i+1, traversable)
		}

		runes += len(by)
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

func flatten(traversable [][]byte) []byte {
	f := make([]byte, 0)
	for _, t := range traversable {
		f = append(f, t...)
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
