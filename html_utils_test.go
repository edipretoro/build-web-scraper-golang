package main

import (
	"testing"
)

func TestGetH1FromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetH1FromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name:      "get_h1_from_html_basic",
			inputBody: "<html><body><h1>Test Title</h1></body></html>",
			expected:  "Test Title",
		},
		{
			name:      "get_h1_from_html_no_h1",
			inputBody: "<html><body><h2>Not H1</h2></body></html>",
			expected:  "",
		},
		{
			name:      "get_h1_from_html_with_whitespace",
			inputBody: "<html><body><h1>   Whitespace Title   </h1></body></html>",
			expected:  "Whitespace Title",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getH1FromHTML(tc.inputBody)
			if actual != tc.expected {
				t.Errorf("FAIL: expected: %q, actual: %q", tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name:      "get_first_paragraph_from_html_basic",
			inputBody: "<html><body><p>This is the first paragraph.</p></body></html>",
			expected:  "This is the first paragraph.",
		},
		{
			name: "get_first_paragraph_from_html_main_priority",
			inputBody: `<html><body>
            <p>Outside paragraph.</p>
            <main>
                <p>Main paragraph.</p>
            </main>
        </body></html>`,
			expected: "Main paragraph.",
		},
		{
			name:      "get_first_paragraph_from_html_no_paragraph",
			inputBody: "<html><body><h1>No paragraphs here</h1></body></html>",
			expected:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.inputBody)
			if actual != tc.expected {
				t.Errorf("FAIL: expected: %q, actual: %q", tc.expected, actual)
			}
		})
	}
}
