package slogutils

import (
	"fmt"
	"log/slog"
	"strings"
)

// CleanSourceAttr makes source attribute more readable
func CleanSourceAttr(groups []string, attr slog.Attr) slog.Attr {
	//fmt.Printf("groups=%#v, attr=%#v\n", groups, attr)
	if attr.Key == "source" {
		source, ok := attr.Value.Any().(*slog.Source)
		if ok {
			// truncate before github.com
			githubIdx := strings.Index(source.File, "github.com")
			if githubIdx > 0 {
				source.File = source.File[githubIdx:]
			}

			//// truncate file path
			//keepParts := 6
			//parts := strings.Split(source.File, "/")
			//nbParts := len(parts)
			//if nbParts > keepParts {
			//	parts = parts[nbParts-keepParts:]
			//	source.File = strings.Join(parts, "/")
			//}

			sourceStr := fmt.Sprintf("%s:%d (%s)", source.File, source.Line, source.Function)
			return slog.String("source", sourceStr)
		}
	}
	return attr
}
