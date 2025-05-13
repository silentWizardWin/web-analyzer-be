package util

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"web-analyzer-be/internal/model"
)

func ParseHTML(r io.Reader) (*html.Node, error) {
	return html.Parse(r)
}

func GetHTMLVersion(resp *http.Response) string {
	for _, c := range resp.Header["Content-Type"] {
		if strings.Contains(c, "html") {
			return "HTML5"
		}
	}
	return "Unknown"
}

func GetTitle(doc *html.Node) string {
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return title
}

func CountHeadings(doc *html.Node) map[string]int {
	counts := map[string]int{}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1", "h2", "h3", "h4", "h5", "h6":
				counts[n.Data]++
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return counts
}

func ContainsLoginForm(doc *html.Node) bool {
	var found bool
	var f func(*html.Node)
	keywords := []string{"login", "log in", "sign in", "password", "email", "username"}

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "form" {
				if containsKeyword(n, keywords) {
					found = true
					return
				}
			} else if n.Data == "input" || n.Data == "button" || n.Data == "label" {
				if containsKeyword(n, keywords) {
					found = true
					return
				}
			}
		}
		if n.FirstChild != nil && !found {
			f(n.FirstChild)
		}
		if n.NextSibling != nil && !found {
			f(n.NextSibling)
		}
	}

	f(doc)
	return found
}

func containsKeyword(n *html.Node, keywords []string) bool {
	for _, attr := range n.Attr {
		val := strings.ToLower(attr.Val)
		for _, kw := range keywords {
			if strings.Contains(val, kw) {
				return true
			}
		}
	}

	// check text content
	if n.Type == html.TextNode {
		text := strings.ToLower(n.Data)
		for _, kw := range keywords {
			if strings.Contains(text, kw) {
				return true
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if containsKeyword(c, keywords) {
			return true
		}
	}

	return false
}


func AnalyzeLinks(doc *html.Node, baseURL string) model.LinkAnalysis {
	internal, external, broken := 0, 0, 0
	visited := make(map[string]bool)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					if href == "" || strings.HasPrefix(href, "#") {
						continue
					}
					absURL := resolveURL(href, baseURL)
					if visited[absURL] {
						continue
					}
					visited[absURL] = true

					if strings.HasPrefix(absURL, baseURL) {
						internal++
					} else {
						external++
					}
					if !isAccessible(absURL) {
						broken++
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return model.LinkAnalysis{
		InternalLinks:     internal,
		ExternalLinks:     external,
		InaccessibleLinks: broken,
	}
}

func resolveURL(href, base string) string {
	baseParsed, err := url.Parse(base)
	if err != nil {
		return href
	}
	hrefParsed, err := url.Parse(href)
	if err != nil {
		return href
	}
	return baseParsed.ResolveReference(hrefParsed).String()
}

func isAccessible(link string) bool {
	resp, err := http.Head(link)
	if err != nil || resp.StatusCode >= 400 {
		return false
	}
	return true
}
