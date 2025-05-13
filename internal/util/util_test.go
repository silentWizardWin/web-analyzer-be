package util

import (
	"net/http"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

const testHTML = `
<!DOCTYPE html>
<html>
<head><title>Test Page</title></head>
<body>
	<h1>Main</h1><h2>Sub</h2><h3>SubSub</h3>
	<form><input type="password"/></form>
	<a href="/internal">Internal</a>
	<a href="https://external.com">External</a>
	<a href="http://invalid.localhost.test">Broken</a>
</body>
</html>
`

func TestParseHTML(t *testing.T) {
	doc, err := ParseHTML(strings.NewReader(testHTML))
	if err != nil || doc == nil {
		t.Errorf("Failed to parse HTML: %v", err)
	}
}

func TestGetTitle(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(testHTML))
	title := GetTitle(doc)
	if title != "Test Page" {
		t.Errorf("Expected 'Test Page', got '%s'", title)
	}
}

func TestCountHeadings(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(testHTML))
	counts := CountHeadings(doc)
	if counts["h1"] != 1 || counts["h2"] != 1 || counts["h3"] != 1 {
		t.Errorf("Unexpected heading counts: %+v", counts)
	}
}

func TestContainsLoginForm(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(testHTML))
	if !ContainsLoginForm(doc) {
		t.Error("Expected login form to be detected")
	}
}

func TestGetHTMLVersion(t *testing.T) {
	resp := &http.Response{
		Header: map[string][]string{"Content-Type": {"text/html"}},
	}
	version := GetHTMLVersion(resp)
	if version != "HTML5" {
		t.Errorf("Expected HTML5, got %s", version)
	}
}

func TestResolveURL(t *testing.T) {
	result := resolveURL("/path", "https://example.com")
	if result != "https://example.com/path" {
		t.Errorf("Unexpected resolved URL: %s", result)
	}
}

func TestIsAccessible(t *testing.T) {
	ok := isAccessible("https://example.com") // should succeed
	if !ok {
		t.Error("Expected example.com to be accessible")
	}
}

func TestIsAccessibleInvalidURL(t *testing.T) {
	if isAccessible("http://invalid.##.com") {
		t.Error("Expected inaccessible, but returned true")
	}
}

func TestAnalyzeLinks_InternalOnly(t *testing.T) {
	htmlContent := `<html><body>
		<a href="/about">About</a>
		<a href="/contact">Contact</a>
	</body></html>`
	doc, _ := html.Parse(strings.NewReader(htmlContent))
	baseURL := "https://example.com"

	result := AnalyzeLinks(doc, baseURL)

	if result.InternalLinks != 2 {
		t.Errorf("Expected 2 internal links, got %d", result.InternalLinks)
	}
	if result.ExternalLinks != 0 {
		t.Errorf("Expected 0 external links, got %d", result.ExternalLinks)
	}
}

func TestAnalyzeLinks_ExternalAndInternal(t *testing.T) {
	htmlContent := `<html><body>
		<a href="/home">Home</a>
		<a href="https://google.com">Google</a>
	</body></html>`
	doc, _ := html.Parse(strings.NewReader(htmlContent))
	baseURL := "https://example.com"

	result := AnalyzeLinks(doc, baseURL)

	if result.InternalLinks != 1 {
		t.Errorf("Expected 1 internal link, got %d", result.InternalLinks)
	}
	if result.ExternalLinks != 1 {
		t.Errorf("Expected 1 external link, got %d", result.ExternalLinks)
	}
}

func TestAnalyzeLinks_InaccessibleLink(t *testing.T) {
	htmlContent := `<html><body>
		<a href="http://invalid.🦄.com">Broken</a>
	</body></html>`
	doc, _ := html.Parse(strings.NewReader(htmlContent))
	baseURL := "https://example.com"

	result := AnalyzeLinks(doc, baseURL)

	if result.InaccessibleLinks == 0 {
		t.Errorf("Expected at least 1 inaccessible link, got %d", result.InaccessibleLinks)
	}
}

func TestAnalyzeLinks_EmptyHTML(t *testing.T) {
	htmlContent := ``
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err == nil && doc != nil {
		result := AnalyzeLinks(doc, "https://example.com")
		if result.InternalLinks != 0 || result.ExternalLinks != 0 {
			t.Errorf("Expected 0 links in empty HTML, got internal: %d, external: %d",
				result.InternalLinks, result.ExternalLinks)
		}
	}
}
