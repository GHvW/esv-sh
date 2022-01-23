package lib

import (
	"testing"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

func TestGivenAVerse(t *testing.T) {
	testStr := []rune("John 3:16")

	t.Run("WhenParsedWithItem", func(t *testing.T) {
		item := Item{}
		result, ok := item.Parse(testStr)
		if !ok {
			t.Errorf("Expected to be able to parse %s to (%v, %s) but nothing was returned", string(testStr), "J", "ohn 3:16")
		}

		if result.Data != 'J' {
			t.Errorf("Expected to be able to parse %s to (data: %v, rest: %s) but data was %s", string(testStr), "J", "ohn 3:16", result.Data)
		}

		if string(result.Rest) != "ohn 3:16" {
			t.Errorf("Expected to be able to parse %s to (data: %s, rest: %s) but rest was %s", string(testStr), "J", "ohn 3:16", string(result.Rest))
		}
	})

	t.Run("AndAZero_WhenParsed", func(t *testing.T) {
		zero := Zero{}
		s, _ := zero.Parse(testStr)

		if s != nil {
			t.Errorf("Expected zero to return nil, but got ParseSuccess with data %v", s.Data)
		}
	})

	t.Run("AndASuccess_WhenParsed", func(t *testing.T) {
		success := &Succeed{testStr}

		s, ok := success.Parse(testStr)

		if !ok {
			t.Errorf("Expected success to return ParseSuccess with data %v, but got nil", string(testStr))
		}

		if string(s.Data.([]rune)) != string(testStr) {
			t.Errorf("Expected success to return ParseSuccess with data %v, but got %v", string(testStr), string(s.Data.([]rune)))
		}
	})

	t.Run("AndAMappedParser_WhenParsed", func(t *testing.T) {
		mapper := &Map{func(rs interface{}) interface{} {
			upper := runes.Map(func(r rune) rune {
				return unicode.ToUpper(r)
			})

			s, _, _ := transform.String(upper, string(rs.([]rune)))
			return s
		}, &Succeed{testStr}}

		s, ok := mapper.Parse(testStr)

		if !ok {
			t.Errorf("Expected mapper to return ParseSuccess with data %v, but got nil", "JOHN 3:16")
		}

		if s.Data != "JOHN 3:16" {
			t.Errorf("Expected mapper to return ParseSuccess with data %v, but got %v", "JOHN 3:16", s.Data)
		}
	})
}
