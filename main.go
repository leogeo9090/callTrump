package main

import (
	"bytes"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"
)

func getLayout(title string) string {
	return `<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<title>` + title + `</title>
			<style>
				@import url(http://fonts.googleapis.com/css?family=Roboto:400,500,700);
				@import url(http://fonts.googleapis.com/css?family=Roboto+Mono:400,700);

				html {
					font-size: 1.125em;
				}

				body {
					background-color: rgb(255, 255, 255);
					color: rgba(0, 0, 0, 0.87);
					font-family: Roboto, sans-serif;
					font-weight: 400;
					line-height: 1.5;
					text-rendering: optimizeLegibility;
				}

				h1, h2, h3, h4 {
					font-weight: 500;
				}

				h1 {
					font-size: 1.602em;
				}

				h2 {
					font-size: 1.424em;
				}

				h3 {
					font-size: 1.266em;
				}

				h4 {
					font-size: 1.125em;
				}

				#page {
					margin: 2.5em auto;
					max-width: 35em;
					padding: 0 0.5rem;
				}

				a {
					text-decoration: none;
					color: rgb(21, 101, 192);
				}

				a:hover {
					text-decoration: underline;
				}

				a:visited {
					color: rgb(21, 101, 192);
				}

				nav ul {
					list-style-type: none;
					padding: 0;
				}

				p.date {
					color: rgba(0, 0, 0, 0.54);
				}

				pre {
					background-color: #f5f5f5;
				}

				code {
					font-family: 'Roboto Mono', monospace;
					font-size: 0.889em;
				}
			</style>
		</head>
		<body>
			<div id="page">`
}

func getFile(f string) []byte {
	b, err := ioutil.ReadFile(f)

	if err != nil {
		panic(err)
	}

	return b
}

func getDir(dir string) []os.FileInfo {
	p, err := ioutil.ReadDir(dir)

	if err != nil {
		panic(err)
	}

	return p
}

func writeFile(fileName string, b bytes.Buffer) {
	err := ioutil.WriteFile(fileName + ".html", b.Bytes(), 0644)

	if err != nil {
		panic(err)
	}
}

func getSiteTitle() string {
	return strings.Split(string(getFile("_sections/header.md")), "\n")[0][2:]
}

func getPostMeta(fi os.FileInfo) (string, string, string) {
	id := fi.Name()[:len(fi.Name()) - 3]
	date := fi.Name()[0:10]
	title := strings.Split(string(getFile("_posts/" + fi.Name())), "\n")[0][2:]

	return id, date, title
}

func getPageMeta(fi os.FileInfo) (string, string) {
	id := fi.Name()[:len(fi.Name()) - 3]
	title := strings.Split(string(getFile("_pages/" + fi.Name())), "\n")[0][2:]

	return id, title
}

func writeIndex() {
	var b bytes.Buffer
	b.WriteString(getLayout(getSiteTitle()))
	b.Write(blackfriday.MarkdownBasic(getFile("_sections/header.md")))
	writePostsSection(&b)
	writePagesSection(&b)
	b.WriteString("</div></body></html>")
	writeFile("index", b)
}

func writePostsSection(b *bytes.Buffer) {
	b.WriteString("<h2>Posts</h2><nav><ul>")

	posts := getDir("_posts")
	limit := int(math.Max(float64(len(posts)) - 5, 0))

	for i := len(posts) - 1; i >= limit; i-- {
		fileName, date, title := getPostMeta(posts[i])

		b.WriteString("<li><a href=\"posts/" +
			fileName + ".html\">" +
			date + " – " +
			title + "</a></li>\n")
	}

	b.WriteString("</ul></nav><p><a href=\"all-posts.html\">All posts</a></p>")
}

func writePagesSection(b *bytes.Buffer) {
	b.WriteString("<h2>Pages</h2><nav><ul>")

	pages := getDir("_pages")

	for i := 0; i < len(pages); i++ {
		id, title := getPageMeta(pages[i])

		b.WriteString("<li><a href=\"pages/" +
			id + ".html\">" +
			title + "</a></li>\n")
	}

	b.WriteString("</ul></nav>")
}

func writePosts() {
	posts := getDir("_posts")

	for i := 0; i < len(posts); i++ {
		id, date, title := getPostMeta(posts[i])

		var b bytes.Buffer
		b.WriteString(getLayout(title + " – " + getSiteTitle()))
		b.WriteString("<p><a href=\"../index.html\">←</a></p>")
		b.WriteString("<p class=\"date\">" + date + "</p>")
		b.Write(blackfriday.MarkdownBasic(getFile("_posts/" + posts[i].Name())))
		b.WriteString("<p><a href=\"../index.html\">←</a></p></div></body></html>")

		writeFile("posts/" + id, b)
	}
}

func writePostsPage() {
	posts := getDir("_posts")
	var b bytes.Buffer

	b.WriteString(getLayout("All posts – " + getSiteTitle()))
	b.WriteString("<p><a href=\"index.html\">←</a></p>")
	b.WriteString("<h1>All posts</h1>")
	b.WriteString("<nav><ul>")

	for i := len(posts) -1; i >= 0; i-- {
		id, date, title := getPostMeta(posts[i])

		b.WriteString("<li><a href=\"posts/" +
			id + ".html\">" +
			date + " – " +
			title + "</a></li>\n")

	}

	b.WriteString("</ul></nav><p><a href=\"index.html\">←</a></p>")
	b.WriteString("</div></body></html>")
	writeFile("all-posts", b)
}

func writePages() {
	pages := getDir("_pages")

	for i := 0; i < len(pages); i++ {
		fileName, title := getPageMeta(pages[i])

		var b bytes.Buffer
		b.WriteString(getLayout(title + " – " + getSiteTitle()))
		b.WriteString("<p><a href=\"../index.html\">←</a></p>")
		b.Write(blackfriday.MarkdownBasic(getFile("_pages/" + pages[i].Name())))
		b.WriteString("<p><a href=\"../index.html\">←</a></p></div></body></html>")

		writeFile("pages/" + fileName, b)
	}
}

func createFilesAndDirs() {
	os.MkdirAll("_sections", 0755)
	os.MkdirAll("_posts", 0755)
	os.MkdirAll("_pages", 0755)

	if _, err := os.Stat("_sections/header.md"); os.IsNotExist(err) {
		err := ioutil.WriteFile(
			"_sections/header.md",
			[]byte("# Title\n\nDescription"),
			0644)

		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat("posts"); os.IsNotExist(err) {
		err := ioutil.WriteFile(
			"_posts/" + time.Now().Format("2006-01-02") + "-initial-post.md",
			[]byte("# Initial post\n\nThis is the initial post."),
			0644)

		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat("pages"); os.IsNotExist(err) {
		err := ioutil.WriteFile(
			"_pages/about.md",
			[]byte("# About\n\nThis is the about page."),
			0644)

		if err != nil {
			panic(err)
		}
	}

	os.MkdirAll("posts", 0755)
	os.MkdirAll("pages", 0755)
}

func main() {
	createFilesAndDirs()
	writeIndex()
	writePosts()
	writePostsPage()
	writePages()
}
