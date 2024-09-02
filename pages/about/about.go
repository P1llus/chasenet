package about

import (
	"bytes"
	"embed"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed posts/*.md
var assets embed.FS

type AboutManager struct {
	files          embed.FS
	aboutPage      AboutMe
	markdownParser goldmark.Markdown
}

type AboutMe struct {
	Title       string
	Description string
	Content     string
	Canonical   string
}

func NewAboutManager() AboutManager {
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	return AboutManager{
		files:          assets,
		aboutPage:      AboutMe{},
		markdownParser: md,
	}
}

func (m *AboutManager) ParseMarkdown(name string) (AboutMe, error) {
	source, err := m.files.ReadFile(name)
	if err != nil {
		return AboutMe{}, err
	}

	// Parse Markdown content
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := m.markdownParser.Convert(source, &buf, parser.WithContext(context)); err != nil {
		return AboutMe{}, err
	}
	metaData := meta.Get(context)
	aboutMe := AboutMe{Title: metaData["Title"].(string), Description: metaData["Description"].(string), Content: buf.String(), Canonical: "https://chasenet.org/about"}

	return aboutMe, nil
}

func (m *AboutManager) LoadAboutPage() error {
	aboutPage, err := m.ParseMarkdown("posts/about.md")
	if err != nil {
		return err
	}
	m.aboutPage = aboutPage

	return nil
}

func (m *AboutManager) GetAboutPage() *AboutMe {
	return &m.aboutPage
}
