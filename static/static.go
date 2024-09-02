package static

import "embed"

//go:embed *.css
var styles embed.FS

func GetStyles() embed.FS {
	return styles
}
