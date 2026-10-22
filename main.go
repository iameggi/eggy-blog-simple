package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
)

const dateLayout = "02-01-2006"

var filenamePattern = regexp.MustCompile(`^\((\d{2}-\d{2}-\d{4})\)([A-Za-z0-9]+(?:-[A-Za-z0-9]+)*)\.md$`)

type Post struct {
	Title   string
	Slug    string
	Date    time.Time
	URL     string
	Content template.HTML
}

type Page struct {
	Title    string
	Root     string
	IsHome   bool
	Blogs    []Post
	Projects []Post
	Post     *Post
}

func main() {
	log.SetFlags(0)
	if err := build(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Selesai: hasil situs tersedia di public/")
}

func build() error {
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		return fmt.Errorf("baca template: %w", err)
	}

	blogs, err := readPosts("blog")
	if err != nil {
		return err
	}
	projects, err := readPosts("project")
	if err != nil {
		return err
	}
	if err := copyStatic(); err != nil {
		return fmt.Errorf("salin aset statis: %w", err)
	}

	for _, section := range []struct {
		name  string
		posts []Post
	}{
		{"blog", blogs},
		{"project", projects},
	} {
		dir := filepath.Join("public", section.name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("buat folder %s: %w", dir, err)
		}
		for _, post := range section.posts {
			page := Page{Title: post.Title, Root: "..", Post: &post}
			if err := writePage(tmpl, filepath.Join(dir, post.Slug+".html"), page); err != nil {
				return err
			}
		}
	}

	return writePage(tmpl, filepath.Join("public", "index.html"), Page{
		Title: "eggy.me", Root: ".", IsHome: true,
		Blogs: blogs, Projects: projects,
	})
}

func readPosts(section string) ([]Post, error) {
	dir := filepath.Join("content", section)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("baca folder %s: %w", dir, err)
	}

	var posts []Post
	seen := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		parts := filenamePattern.FindStringSubmatch(entry.Name())
		if parts == nil {
			return nil, fmt.Errorf("%s: nama file harus (DD-MM-YYYY)slug-judul.md", path)
		}
		date, err := time.Parse(dateLayout, parts[1])
		if err != nil {
			return nil, fmt.Errorf("tanggal tidak valid pada %s: %w", path, err)
		}
		slug := parts[2]
		key := strings.ToLower(slug)
		if seen[key] {
			return nil, fmt.Errorf("%s: slug duplikat %q", dir, slug)
		}
		seen[key] = true

		source, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("baca %s: %w", path, err)
		}
		var body bytes.Buffer
		if err := goldmark.Convert(source, &body); err != nil {
			return nil, fmt.Errorf("render Markdown %s: %w", path, err)
		}
		posts = append(posts, Post{
			Title: slug, Slug: slug, Date: date,
			URL:     section + "/" + slug + ".html",
			Content: template.HTML(body.String()),
		})
	}

	sort.Slice(posts, func(i, j int) bool {
		if posts[i].Date.Equal(posts[j].Date) {
			return posts[i].Slug < posts[j].Slug
		}
		return posts[i].Date.After(posts[j].Date)
	})
	return posts, nil
}

func copyStatic() error {
	return filepath.WalkDir("static", func(path string, entry fs.DirEntry, err error) error {
		if path == "static" && os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		relative, err := filepath.Rel("static", path)
		if err != nil {
			return err
		}
		target := filepath.Join("public", relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("%s: aset harus berupa file biasa", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}

func writePage(tmpl *template.Template, path string, page Page) error {
	var output bytes.Buffer
	if err := tmpl.Execute(&output, page); err != nil {
		return fmt.Errorf("render %s: %w", path, err)
	}
	if err := os.WriteFile(path, output.Bytes(), 0644); err != nil {
		return fmt.Errorf("tulis %s: %w", path, err)
	}
	return nil
}
