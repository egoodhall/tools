package sync

import (
	"fmt"
	"strings"
)

func GetUrls(users []string) []string {
	urls := make([]string, len(users))
	for i, user := range users {
		if strings.HasPrefix(user, "gitlab") {
			urls[i] = fmt.Sprintf("https://gitlab.com/%s.keys", strings.TrimPrefix("gitlab:", user))
		} else {
			urls[i] = fmt.Sprintf("https://github.com/%s.keys", user)
		}
	}
	return urls
}
