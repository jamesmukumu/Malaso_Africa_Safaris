package slug

import (
	"regexp"
	"strings"
)

func SlugMaker(title string)string {
slug := strings.TrimSpace(title)
slug = strings.ToLower(slug)
reg := regexp.MustCompile(`[^a-z0-9\s]+`)
slug = reg.ReplaceAllString(slug,"")

finalSlug := strings.ReplaceAll(slug," ","-")
return finalSlug
}
