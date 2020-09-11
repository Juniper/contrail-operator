package randomstring

import (
	"testing"
	"unicode"
)

func hasLettersOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

type TestInterface interface {
	Generate() string
}

func TestRandomstring(t *testing.T) {

	checkGenerateString := func(t *testing.T, testinterface TestInterface) {
		t.Helper()
		got := testinterface.Generate()
		if !hasLettersOnly(got) {
			t.Errorf("failed to verify Generate")
		}
	}

	t.Run("verification", func(t *testing.T) {
		testvalue := RandString{Size: 12}
		checkGenerateString(t, testvalue)
	})
}
