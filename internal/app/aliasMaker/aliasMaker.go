package aliasmaker

import (
	"math/rand"
)

const (
	aliasKeyLen        int = 9
	TrysToMakeAliasKey int = 5
)

// --------------------------------------------------
//
//	Make short alias from URL
//	Output:
//		alias string - short alias to "longURL"
func CreateAliasKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	alliasKey := make([]byte, aliasKeyLen)

	for i := range alliasKey {
		alliasKey[i] = charset[rand.Intn(len(charset))]
	}

	return string(alliasKey)
}
