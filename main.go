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
	return `<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<title>` + title + `</title>
			<style>
				@import url(http://fonts.googleapis.com/css?family=Open+Sans);

				#page {
					margin: 2.5em auto;
					max-width: 30em;
					padding: 0 0.5rem;
				}

				a {
					text-decoration: none;
					color: #03c;
				}

				a:visited {
					color: #639;
				}

				a:hover {
					text-decoration: underline;
				}

				nav ul {
					list-style-type: none;
					padding: 0;
				}

				html {
					font-size: 1em;
				}

				body {
				  background-color: #f9f9f9;
				  color: #333;
				  font-weight: 400;
				  line-height: 1.45;
				  font-family: 'Open Sans', sans-serif;
				  text-rendering: optimizeLegibility;
				}

				p {
					margin-bottom: 1.3em;
				}

				h1, h2, h3, h4 {
				  margin: 1.414em 0 0.5em;
				  font-weight: 600;
				  line-height: 1.2;
				}

				h1 {
				  margin-top: 0;
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

				pre {
					background-color: #eee;
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

func siteTitle() string {
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

func writeLayout(b *bytes.Buffer, title string) {
	b.WriteString(getLayout(title))
}

func writeIndex() {
	var b bytes.Buffer
	writeLayout(&b, siteTitle())
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

	b.WriteString("</nav></ul>")
}

func writePosts() {
	posts := getDir("_posts")

	for i := 0; i < len(posts); i++ {
		var b bytes.Buffer

		id, date, title := getPostMeta(posts[i])

		writeLayout(&b, title + " – " + siteTitle())
		b.WriteString("<p><a href=\"../index.html\">←</a></p>")
		b.WriteString("<p>" + date + "</p>")
		b.Write(blackfriday.MarkdownBasic(getFile("_posts/" + posts[i].Name())))
		b.WriteString("<p><a href=\"../index.html\">←</a></p></div></body></html>")

		writeFile("posts/" + id, b)
	}
}

func writePostsPage() {
	posts := getDir("_posts")
	var b bytes.Buffer

	writeLayout(&b, "All posts – " + siteTitle())
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
		var b bytes.Buffer

		fileName, title := getPageMeta(pages[i])

		writeLayout(&b, title + " – " + siteTitle())
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
