package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 60,
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

const baseUrl string = "https://inoveltranslation.com"

type novelData struct {
	url string
	img string
}

func main() {
	request, err := http.NewRequest(http.MethodGet, baseUrl+"/novels", nil)
	handleError(err)
	response, err := client.Do(request)
	handleError(err)

	body, err := io.ReadAll(response.Body)
	handleError(err)
	err = response.Body.Close()
	handleError(err)

	document, err := html.Parse(strings.NewReader(string(body)))
	handleError(err)

	novels := make(map[string]novelData)

	for bodyChild := range document.LastChild.LastChild.ChildNodes() {
		if bodyChild.Data == "main" {
			for novelLink := range bodyChild.FirstChild.FirstChild.LastChild.ChildNodes() {
				if novelLink.Type == html.ElementNode {
					novels[novelLink.LastChild.FirstChild.Data] = novelData{
						url: novelLink.Attr[4].Val,
						img: novelLink.FirstChild.Attr[1].Val,
					}
				}
			}
		}
	}

	novelName := os.Args[1]
	request, err = http.NewRequest(http.MethodGet, baseUrl+novels[novelName].url, nil)
	handleError(err)
	response, err = client.Do(request)
	handleError(err)

	body, err = io.ReadAll(response.Body)
	handleError(err)
	err = response.Body.Close()
	handleError(err)

	document, err = html.Parse(strings.NewReader(string(body)))
	handleError(err)

	chapterLinks := make([]string, 0)

essa:
	for bodyChild := range document.LastChild.LastChild.ChildNodes() {
		if bodyChild.Data == "main" {
			for mainChild := range bodyChild.FirstChild.FirstChild.ChildNodes() {
				if mainChild.Data == "section" {
					for chapter := range mainChild.LastChild.ChildNodes() {
						if chapter.FirstChild != nil && chapter.LastChild.Data == "a" {
							chapterLinks = append(chapterLinks, chapter.FirstChild.Attr[3].Val)
						}
					}
					break essa
				}
			}
		}
	}

	var file *os.File
	items := make([]string, 0)
	itemref := make([]string, 0)
	nav := make([]string, 0)

	for _, link := range chapterLinks {
		if link != "/chapters/965daeb4-f291-475f-9925-b394f0bb0914" {
			continue
		}
		request, err = http.NewRequest(http.MethodGet, baseUrl+link, nil)
		handleError(err)
		response, err = client.Do(request)
		handleError(err)

		body, err = io.ReadAll(response.Body)
		handleError(err)
		err = response.Body.Close()
		handleError(err)
		document, err = html.Parse(strings.NewReader(string(body)))
		handleError(err)
		fmt.Println(string(body))

		for bodyChild := range document.LastChild.LastChild.ChildNodes() {
			if bodyChild.Data == "main" {
				var chapterNumber string
				for mainChild := range bodyChild.ChildNodes() {
					if mainChild.Data == "header" {
						name := mainChild.FirstChild.FirstChild.Data
						num := strings.Index(name, "-")
						if num == -1 {
							num = 6
						}
						chapterNumber = name[4 : num-1]
						items = append(items, "<item id=\"c"+chapterNumber+"\" href=\"xhtml/"+chapterNumber+".xhtml\" media-type=\"application/xhtml+xml\"/>")
						itemref = append(itemref, "<itemref idref=\"c"+chapterNumber+"\"/>")
						nav = append(nav, "<li>\n                    <a href=\"xhtml/"+chapterNumber+".xhtml\">\n                        Chapter "+chapterNumber+"\n                    </a>\n                </li>")
						file, err = os.Create("./I Became A Flashing Genius At The Magic Academy/OEBPS/xhtml/" + chapterNumber + ".xhtml")
						file.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<html xmlns=\"http://www.w3.org/1999/xhtml\">\n    <head>\n        <meta charset=\"utf-8\" />\n        <title>I Became A Flashing Genius At The Magic Academy</title>\n<link rel=\"stylesheet\" href=\"../styles/styles.css\" type=\"text/css\"/>\n    </head>\n    <body>")
						handleError(err)
					}
					if mainChild.Data == "div" {
						if mainChild.Attr[0].Key == "style" {
							_, err = file.WriteString("<div style=\"border: solid;height: 2px; width: 100%; background-color: hsl(var(--border)); margin: 1rem 0px;\"></div>\n")
							// fmt.Println("kurwa?")
							handleError(err)
						}
					}
					if mainChild.Data == "p" {
						if mainChild.FirstChild != nil {
							_, err = file.WriteString("<p>")
							i := 0
							for text := range mainChild.ChildNodes() {
								if text.Data == "span" {
									_, err = file.WriteString("<span style=\"" + text.Attr[0].Val + "\">" + strings.ReplaceAll(strings.ReplaceAll(text.FirstChild.Data, "&", "&amp;"), "<", "&lt;") + "</span>")
									handleError(err)
									// if i != 0 {
									// 	fmt.Printf("Chapter: %s\n", chapterNumber)
									// }
									continue
								}
								_, err = file.WriteString(strings.ReplaceAll(strings.ReplaceAll(text.Data, "&", "&amp;"), "<", "&lt;"))
								// fmt.Println(text.Data)
								handleError(err)
								i++
							}
							_, err = file.WriteString("</p>\n")
						}
					}
				}
				file.WriteString("</body>\n</html>")
				file.Close()
				break
			}
		}
	}

	file, err = os.Create("./I Became A Flashing Genius At The Magic Academy/OEBPS/package.opf")
	handleError(err)
	file.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<package version=\"3.0\" xml:lang=\"en\" xmlns=\"http://www.idpf.org/2007/opf\" unique-identifier=\"pub-id\" dir=\"ltr\">\n    <metadata xmlns:dc=\"http://purl.org/dc/elements/1.1/\">\n        <dc:identifier id=\"pub-id\">__KURWA__</dc:identifier>\n        <dc:language>en</dc:language>\n        <dc:title>I Became A Flashing Genius At The Magic Academy</dc:title>\n        <dc:creator id=\"creator\">은밀히</dc:creator>\n        <meta property=\"dcterms:modified\">2024-11-23T15:12:00Z</meta>\n        <meta content=\"cover-image\" name=\"cover\"/>\n    </metadata>\n    <manifest>\n        <item id=\"font\" href=\"fonts/MonoLisa.ttf\" media-type=\"application/font-sfnt\"/>\n        <item id=\"styles\" href=\"styles/styles.css\" media-type=\"text/css\"/>\n        <item properties=\"cover-image\" id=\"cover-image\" href=\"images/cover.webp\" media-type=\"image/webp\"/>\n        <item properties=\"nav\" id=\"nav\" href=\"nav.xhtml\" media-type=\"application/xhtml+xml\"/>\n <item id=\"cover\" href=\"xhtml/cover.xhtml\" media-type=\"application/xhtml+xml\"/>")
	for _, item := range slices.Backward(items) {
		file.WriteString(item + "\n")
	}
	file.WriteString("</manifest>\n    <spine page-progression-direction=\"ltr\">\n<itemref idref=\"cover\"/>\n        <itemref idref=\"nav\"/>")
	for _, item := range slices.Backward(itemref) {
		file.WriteString(item + "\n")
	}
	file.WriteString("    </spine>\n</package>")
	file.Close()
	file, err = os.Create("./I Became A Flashing Genius At The Magic Academy/OEBPS/nav.xhtml")
	handleError(err)
	file.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<html xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:epub=\"http://www.idpf.org/2007/ops\">\n    <head>\n        <title>EPUB 3 Navigation Document</title>\n    </head>\n    <body>\n        <nav epub:type=\"toc\">\n            <ol>")
	for _, item := range slices.Backward(nav) {
		file.WriteString(item + "\n")
	}
	file.WriteString(" </ol>\n        </nav>\n  <nav epub:type=\"landmarks\">\n            <ol>\n                <li><a href=\"xhtml/cover.xhtml\" epub:type=\"cover\">Cover</a></li>\n                <li><a href=\"xhtml/0.xhtml\" epub:type=\"bodymatter\">Start</a></li>\n            </ol>\n        </nav>  </body>\n</html>")
	file.Close()
	file, err = os.Create("./I Became A Flashing Genius At The Magic Academy/OEBPS/xhtml/cover.xhtml")
	handleError(err)
	file.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<html xmlns=\"http://www.w3.org/1999/xhtml\">\n<head>\n    <meta charset=\"utf-8\" />\n    <title>Chapter 1</title>\n    <link rel=\"stylesheet\" href=\"../styles/styles.css\" type=\"text/css\"/>\n</head>\n<body>\n    <img src=\"../images/cover.webp\" alt=\"kurwa\"/>\n</body>\n</html>")
	file.Close()
}
