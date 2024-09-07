package blog

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
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

// The argument is the string part of the URL after /blog/someslug,
// which will retrieve the matching blogpost based on slug value in the markdown file.
func (m *BlogManager) GetBlogPostBySlug(slug string) *BlogPost {
	for _, post := range m.blogPosts.Posts {
		if post.Slug == slug {
			return &post
		}
	}
	return nil
}

// Gets all blogposts with a specific tag
func (m *BlogManager) GetBlogPostsByTag(tag string) *BlogPosts {
	var posts BlogPosts
	for _, post := range m.blogPosts.Posts {
		for _, t := range post.Tags {
			if t == tag {
				posts.Posts = append(posts.Posts, post)
			}
		}
	}
	return &posts
}

// Returns an already sorted list of posts
func (m *BlogManager) ListBlogPosts() *BlogPosts {
	return &m.blogPosts
}

// Read all blog markdown files from the embedded fs which includes all files in ./posts
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
	// Sort posts by date, newest first
	sort.Slice(blogPosts.Posts, func(i, j int) bool {
		dateI, errI := time.Parse("02-Jan-2006", blogPosts.Posts[i].Date)
		dateJ, errJ := time.Parse("02-Jan-2006", blogPosts.Posts[j].Date)

		// If there's an error parsing either date, consider the post with the error to be older
		if errI != nil {
			return false
		}
		if errJ != nil {
			return true
		}

		return dateI.After(dateJ)
	})
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
