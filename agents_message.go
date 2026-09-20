package openai

import "strings"

// OutputText joins this message's output_text blocks in order, without filtering
// by phase or fetching additional content. Input text and images are ignored.
func (r AgentSessionMessage) OutputText() string {
	var text strings.Builder
	for _, content := range r.Content {
		if content.Type == "output_text" {
			text.WriteString(content.Text)
		}
	}
	return text.String()
}
