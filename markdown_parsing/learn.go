package markdown_parsing

import (
	"errors"
	"fmt"
	"raiden_fumo/lang_notebook_core/markdown_parsing/extensions"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

	"github.com/yuin/goldmark/extension"
)

func makeWalker(source []byte, learning_lists *map[gast.Node][]string) gast.Walker {
	visited := make(map[gast.Node]bool)
	return func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if !visited[n] {
			if n.Kind() == extensions.KindExclamationMarkCommand {
				nextSibling := n.NextSibling()
				if nextSibling != nil {
					(*learning_lists)[nextSibling] = n.(*extensions.ExclamationMarkCommandNode).Tags
				}
				return gast.WalkSkipChildren, nil
			} else if strings.HasPrefix(string(n.Text(source)), "!learn") {
				nextSibling := n.NextSibling()
				if nextSibling != nil {
					fmt.Println("    ", nextSibling.Kind())
				}
			}
			visited[n] = !visited[n]
		}

		return gast.WalkContinue, nil
	}
}

type Pair struct {
	First string
	Second string
	Tags *[]string
}

func parsePair(line string) (string, string, error) {
	re := regexp.MustCompile(`^\s*(.+)\s*:\s*(.+)\s*$`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		return matches[1], matches[2], nil
	}
	return "", "", errors.New("Line doesn't contain a string pair")
}

func makeListWalker(source []byte, dest *[]Pair, tags []string) gast.Walker {
	visited := make(map[gast.Node]bool)
	return func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if !visited[n] && n.Kind() == gast.KindListItem {
			str := (string)(n.Text(source))
			a, b, err := parsePair(str)
			if err != nil {
				panic("Bad string " + str)
				// return gast.WalkStop, err
			}
			*dest = append(*dest, Pair{First: a, Second: b, Tags: &tags})
		}
		visited[n] = !visited[n]

		return gast.WalkContinue, nil
	}
}

func ExtractLearnPairs(source []byte) []Pair {
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.Table, extensions.ExclamationMarkCommand),
	)
	doc := markdown.Parser().Parse(text.NewReader(source))
	var buf = make(map[gast.Node][]string)
	gast.Walk(doc, makeWalker(source, &buf))
	var dest = make([]Pair, 0)
	for k, tags := range buf {
		walker := makeListWalker(source, &dest, tags)
		gast.Walk(k, walker)
	}
	return dest
}
