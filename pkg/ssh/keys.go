package ssh

import (
	"fmt"
	"strings"
)

func ParseAuthorizedKey(from string) (AuthorizedKey, error) {
	var authorizedKey AuthorizedKey

	typ, rest, found := strings.Cut(strings.Trim(from, " "), " ")
	if !found {
		return authorizedKey, fmt.Errorf("couldn't parse type")
	}
	authorizedKey.Type = strings.Trim(typ, " ")

	key, comment, found := strings.Cut(strings.Trim(rest, " "), " ")
	if !found {
		authorizedKey.Key = strings.Trim(rest, " ")
		return authorizedKey, nil
	}
	authorizedKey.Key = strings.Trim(key, " ")
	authorizedKey.Comment = strings.Trim(comment, " ")

	return authorizedKey, nil
}

type AuthorizedKey struct {
	Type    string
	Key     string
	Comment string
}

func (key *AuthorizedKey) WithoutComment() AuthorizedKey {
	return AuthorizedKey{Type: key.Type, Key: key.Key}
}

func (key *AuthorizedKey) String() string {
	return strings.Trim(fmt.Sprintf("%s %s %s", key.Type, key.Key, key.Comment), " ")
}
