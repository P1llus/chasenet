package blog

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed posts/*.md
var blogFs embed.FS

type BlogManager struct {
	blogPosts      BlogPosts
	markdownParser goldmark.Markdown
}

type BlogPost struct {
	Title       string
	Description string
	Date        string
	Slug        string
	Content     string
	Canonical   string
	Tags        []interface{}
}

type BlogPosts struct {
	Posts []BlogPost
}

func NewBlogManager() BlogManager {
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
		goldmark.WithExtensions(
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
	)

	return BlogManager{
		blogPosts:      BlogPosts{},
		markdownParser: md,
	}
}

func (m *BlogManager) GetBlogPostBySlug(slug string) *BlogPost {
	for _, post := range m.blogPosts.Posts {
		if post.Slug == slug {
			return &post
		}
	}
	return nil
}

func (m *BlogManager) ListBlogPosts() *BlogPosts {
	return &m.blogPosts
}

func (m *BlogManager) LoadBlogPosts() error {
	var files []string
	if err := fs.WalkDir(blogFs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	}); err != nil {
		return err
	}
	var blogPosts BlogPosts
	for _, file := range files {
		post, err := m.parseMarkdown(file)
		if err != nil {
			return err
		}
		blogPosts.Posts = append(blogPosts.Posts, post)
	}
	m.blogPosts = blogPosts

	return nil
}

func (m *BlogManager) parseMarkdown(name string) (BlogPost, error) {
	source, err := blogFs.ReadFile(name)
	if err != nil {
		return BlogPost{}, err
	}

	// Parse Markdown content
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := m.markdownParser.Convert(source, &buf, parser.WithContext(context)); err != nil {
		return BlogPost{}, err
	}
	metaData := meta.Get(context)
	slug := metaData["Slug"].(string)
	tags := metaData["Tags"].([]interface{})

	date, err := time.Parse("2/1/2006", metaData["Date"].(string))
	if err != nil {
		log.Fatal(err)
	}
	blogPost := BlogPost{Title: metaData["Title"].(string), Description: metaData["Description"].(string), Slug: slug, Content: buf.String(), Canonical: fmt.Sprintf("https://chasenet.org/posts/%s", slug), Date: date.Format("02-Jan-2006"), Tags: tags}

	return blogPost, nil
}
