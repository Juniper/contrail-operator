package randomstring

import (
	"regexp"
	"testing"
	"unicode"
        "math/rand"
        "time"
)

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsValid(s string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	return re.MatchString(s)
}
/*
type StubSize struct {
	size int
}
*/
type StubData interface {
	Generate() string
}

func TestRandomstring(t *testing.T) {

	checkGenerateString := func(t *testing.T, stubdata StubData) {
		t.Helper()
		got := stubdata.Generate()
		if !IsLetter(got) {
			t.Errorf("failed to verify Generate")
		}
	}

	t.Run("Random string Generate method verification1", func(t *testing.T) {
		stubsize := RandString{12}
		checkGenerateString(t, stubsize)
	})

	t.Run("Random string Generate method verification2", func(t *testing.T) {
                rand.Seed(time.Now().UnixNano())
                min := 10
                max := 30
		stubsize := RandString{rand.Intn(max - min + 1) + min}
		checkGenerateString(t, stubsize)
	})

	t.Run("Random string bytes function verification1", func(t *testing.T) {
		value2 := 3
		got := randStringBytes(value2)
		if !IsValid(got) {
			t.Errorf("failed to verify random bytes")
		}

	})

	t.Run("Random string bytes function verification2", func(t *testing.T) {
                rand.Seed(time.Now().UnixNano())
                min := 10
                max := 30
		got := randStringBytes(rand.Intn(max - min + 1) + min)
		if !IsValid(got) {
			t.Errorf("failed to verify random bytes")
		}

	})
}
