package extensions

import (
	"go/token"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type ExclamationMarkCommandNode struct {
	gast.BaseBlock
	command string
	Tags []string
}

func parseCommand(line string) *ExclamationMarkCommandNode {
	node := &ExclamationMarkCommandNode{}
	tags := make([]string, 0)
	words := strings.Fields(line)
	for i, word := range words {
		if i == 0 {
			if strings.HasPrefix(word, "!") {
				node.command = strings.TrimPrefix(word, "!")
			}
		} else {
			if strings.HasPrefix(word, "#") {
				tags = append(tags, strings.TrimPrefix(word, "#"))
			}
		}
	}
	node.Tags = tags
	return node
}

var KindExclamationMarkCommand = gast.NewNodeKind("ExclamationMarkCommand")

func (*ExclamationMarkCommandNode) Kind() gast.NodeKind {
	return KindExclamationMarkCommand
}

func (c *ExclamationMarkCommandNode) Text(buf []byte) []byte {
	var ret string = ""
	for _, tag := range c.Tags {
		ret += tag + " "
	}
	return ([]byte)(ret)
}

func (c *ExclamationMarkCommandNode) Dump(source []byte, level int) {
	m := map[string]string{
		"Command": c.command,
		"Tags": strings.Join(c.Tags, ", "),
	}
	gast.DumpHelper(c, source, level, m, nil)
}

func (*ExclamationMarkCommandNode) Pos() token.Pos {
	panic("unimplemented")
}

type exclamationMarkCommandParser struct {
}

func (p *exclamationMarkCommandParser) Trigger() []byte {
	return []byte{'!'}
}

func (p *exclamationMarkCommandParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, segment := reader.PeekLine()
	if strings.HasPrefix(string(line), "!") {
		node := parseCommand(string(line))
		node.Lines().Append(text.NewSegment(segment.Start, segment.Stop - 1))
		reader.Advance(segment.Len() - 1)
		return node, parser.NoChildren
	}
	return nil, parser.NoChildren
}


func (p *exclamationMarkCommandParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {
}

func (p *exclamationMarkCommandParser) CanAcceptIndentedLine() bool {
	return true
}

func (p *exclamationMarkCommandParser) CanInterruptParagraph() bool {
	return true
}

func (b *exclamationMarkCommandParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

type exclamationMarkCommandExtension struct {
}

var ExclamationMarkCommand = &exclamationMarkCommandExtension{}

func NewExclamationMarkCommandParser() parser.BlockParser {
	return &exclamationMarkCommandParser{}
}

func (e *exclamationMarkCommandExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(NewExclamationMarkCommandParser(), 500),
	))
}
