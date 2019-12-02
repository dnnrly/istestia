package markdown

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

type CodeBlock struct {
	Type     string
	Contents string
}

func Extract(doc string) ([]CodeBlock, error) {
	r := text.NewReader([]byte(doc))

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
	node := md.Parser().Parse(r)

	blocks := []CodeBlock{}
	walker := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		b, ok := n.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}

		//contents := string(b.Text(r.Source()))
		//contents := fmt.Sprintf("%v", b.Lines())
		contents := ""
		for i := 0; i < b.Lines().Len(); i++ {
			line := b.Lines().At(i)
			contents += string(line.Value(r.Source()))
		}

		block := CodeBlock{
			Type:     string(b.Language(r.Source())),
			Contents: contents,
		}

		blocks = append(blocks, block)

		return ast.WalkContinue, nil
	}

	err := ast.Walk(node, walker)

	if err != nil {
		return []CodeBlock{}, err
	}

	return blocks, nil
}
