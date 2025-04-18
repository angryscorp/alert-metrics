package zipper

const minSize = 1024

func isContentZippable(contentType string) bool {
	return contentType == "application/json" || contentType == "text/html"
}
