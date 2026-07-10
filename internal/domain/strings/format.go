package strings

const previewTextLen = 32

func TruncatePreviewText(s string) string {
	runes := []rune(s)
	if len(runes) <= previewTextLen {
		return s
	}
	return string(runes[:previewTextLen]) + "..."
}
