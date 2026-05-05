// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package openai

import (
	"net/url"
	"strings"
)

func pathSegment(value string) string {
	segments := strings.Split(value, "/")
	for i, segment := range segments {
		switch segment {
		case ".":
			segments[i] = "%2E"
		case "..":
			segments[i] = "%2E%2E"
		default:
			segments[i] = url.PathEscape(segment)
		}
	}
	return strings.Join(segments, "%2F")
}
