package lib

import (
	"strings"
	"testing"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

// TODO - refactor this and split out tests on strings, digits, spaces, etc.
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
		}, Success(testStr)}

		s, ok := mapper.Parse(testStr)

		if !ok {
			t.Errorf("Expected mapper to return ParseSuccess with data %v, but got nil", "JOHN 3:16")
		}

		if s.Data != "JOHN 3:16" {
			t.Errorf("Expected mapper to return ParseSuccess with data %v, but got %v", "JOHN 3:16", s.Data)
		}
	})

	t.Run("WhenParsedWithCombinatorMapIntegration", func(t *testing.T) {
		s, ok := TakeItem().Map(func(r interface{}) interface{} {
			return unicode.ToLower(r.(rune))
		}).Parse(testStr)

		// g, gok := TakeItem().Map(unicode.ToLower).Parse(testStr) // unfortunately, this doesn't work yet

		if !ok {
			t.Errorf("Expected to be able to parse %s and map to (data: %v, rest: %s) but nothing was returned", string(testStr), "j", "ohn 3:16")
		}

		if s.Data != 'j' {
			t.Errorf("Expected to be able to parse %s and map to (data: %v, rest: %s) but data was %s and rest was %s", string(testStr), "j", "ohn 3:16", string(s.Data.(rune)), string(s.Rest))
		}

		if string(s.Rest) != "ohn 3:16" {
			t.Errorf("Expected to be able to parse %s and map to (data: %v, rest: %s) but data was %s and rest was %s", string(testStr), "j", "ohn 3:16", string(s.Data.(rune)), string(s.Rest))
		}
	})

	t.Run("And A FlatMapped Parser, When Parsed", func(t *testing.T) {
		flatmapper := &FlatMap{func(rs interface{}) Parser {
			return Success(rs)
		}, TakeItem()}

		s, ok := flatmapper.Parse(testStr)

		if !ok {
			t.Errorf("Expected flatmapper to return ParseSuccess with data %v and rest %v, but got nil", "J", "ohn 3:16")
		}

		if string(s.Data.(rune)) != "J" {
			t.Errorf("Expected flatmapper to return ParseSuccess with data %v and rest %v, but got %v, %v", "J", "ohn 3:16", string(s.Data.(rune)), string(s.Rest))
		}

		if string(s.Rest) != "ohn 3:16" {
			t.Errorf("Expected flatmapper to return ParseSuccess with data %v and rest %v, but got %v, %v", "J", "ohn 3:16", string(s.Data.(rune)), string(s.Rest))
		}
	})

	t.Run("When Parsed With Combinator FlatMap Integration", func(t *testing.T) {

		s, ok := TakeItem().FlatMap(func(rs interface{}) Parser {
			return Success(rs)
		}).Parse(testStr)

		if !ok {
			t.Errorf("Expected flatmapper to return ParseSuccess with data %v and rest %v, but got nil", "J", "ohn 3:16")
		}

		if string(s.Data.(rune)) != "J" {
			t.Errorf("Expected flatmapper to return ParseSuccess with data %v and rest %v, but got %v, %v", "J", "ohn 3:16", string(s.Data.(rune)), string(s.Rest))
		}

		if string(s.Rest) != "ohn 3:16" {
			t.Errorf("Expected flatmapper to return ParseSuccess with data %v and rest %v, but got %v, %v", "J", "ohn 3:16", string(s.Data.(rune)), string(s.Rest))
		}
	})

	t.Run("When parsed with a satisfy test", func(t *testing.T) {
		sat := &Satisfy{unicode.IsLetter}

		s, ok := sat.Parse(testStr)

		if !ok {
			t.Errorf("Expected satisfy to return ParseSuccess with data %v and rest %v, but got nil", "J", "ohn 3:16")
		}

		if s.Data.(rune) != 'J' {
			t.Errorf("Expected satisfy to return ParseSuccess with data %v and rest %v, but got %v, %v", "J", "ohn 3:16", string(s.Data.(rune)), string(s.Rest))
		}
	})

	// TODO - find a better way to test this
	t.Run("When parsing multiple Alphabetic characters", func(t *testing.T) {
		many := Multiple(Alpha())

		s, ok := many.Parse(testStr)

		if !ok {
			t.Errorf("Expected multiple to return ParseSuccess with data %v and rest %v, but got nil", "John", "' 3:16'")
		}

		expected := []rune("John")
		for i, it := range s.Data.([]interface{}) {
			if expected[i] != it.(rune) {
				t.Errorf("Expected %v in position %v of word 'John', but was %v", expected[i], i, it)
			}
		}

		if string(s.Rest) != " 3:16" {
			t.Errorf("Expected multiple to return ParseSuccess with data %v and rest %v, but rest was %v", "John", " 3:16", string(s.Rest))
		}
	})

	t.Run("When parsing multiple Alphabetic characters", func(t *testing.T) {
		many := AtLeastOne(Alpha())

		s, ok := many.Parse(testStr)

		if !ok {
			t.Errorf("Expected multiple to return ParseSuccess with data %v and rest %v, but got nil", "John", "' 3:16'")
		}

		expected := []rune("John")
		for i, it := range s.Data.([]interface{}) {
			if expected[i] != it.(rune) {
				t.Errorf("Expected %v in position %v of word 'John', but was %v", expected[i], i, it)
			}
		}

		if string(s.Rest) != " 3:16" {
			t.Errorf("Expected multiple to return ParseSuccess with data %v and rest %v, but rest was %v", "John", " 3:16", string(s.Rest))
		}
	})

	t.Run("When parsing at least one Alphabetic character and the string is digits", func(t *testing.T) {
		many := AtLeastOne(Alpha())

		_, ok := many.Parse([]rune("123"))

		if ok {
			t.Errorf("Expected alphabetic parser to fail when given digits, but it didn't")
		}
	})

	t.Run("When parsing a natural number", func(t *testing.T) {
		many := NaturalNumber()

		result, ok := many.Parse([]rune("123a"))

		if !ok {
			t.Errorf("Expected success")
		}

		if result.Data.(int) != 123 {
			t.Errorf("Expected 123, but got %v", result.Data)
		}

		if string(result.Rest) != "a" {
			t.Errorf("Expected 'a', but got %v", string(result.Rest))
		}
	})

	t.Run("When parsing two numbers separated by a dash", func(t *testing.T) {
		parser := NaturalNumber().SeparatedBy(Rune('-'))

		result, ok := parser.Parse([]rune("12-25a"))

		if !ok {
			t.Errorf("Expected success")
		}

		if result.Data.([]interface{})[0].(int) != 12 {
			t.Errorf("Expected first item to be 12, but got %v", result.Data.([]interface{})[0].(int))
		}

		if result.Data.([]interface{})[1].(int) != 25 {
			t.Errorf("Expected second item to be 25, but got %v", result.Data.([]interface{})[1].(int))
		}

		if string(result.Rest) != "a" {
			t.Errorf("Expected 'a', but got %v", string(result.Rest))
		}
	})
}

func TestWhiteSpace(t *testing.T) {
	testStr := []rune("\tab")

	t.Run("When parsing a string of whitespace then other characters", func(t *testing.T) {
		parser := WhiteSpace()

		result, ok := parser.Parse(testStr)

		if !ok {
			t.Errorf("Expected success")
		}

		if string(result.Data.(rune)) != "\t" {
			t.Errorf("Expected to parse all whitespace, but did not")
		}

		if string(result.Rest) != "ab" {
			t.Errorf("Expected 'ab', but got %v", string(result.Rest))
		}
	})
}

func TestAnd(t *testing.T) {
	testStr := []rune("J3")

	t.Run("When parsing a string for two items using And", func(t *testing.T) {
		parser := Alpha().And(Digit())

		result, ok := parser.Parse(testStr)

		if !ok {
			t.Errorf("Expected success")
		}

		if result.Data.(*Pair).First.(rune) != 'J' {
			t.Errorf("Expected first item to be 'J', but got %v", result.Data.(*Pair).First.(rune))
		}

		if result.Data.(*Pair).Second.(rune) != '3' {
			t.Errorf("Expected second item to be '3', but got %v", result.Data.(*Pair).Second.(rune))
		}
	})
}

func TestParseBook(t *testing.T) {
	john := []rune("John 1:1")
	firstJohn := []rune("1 John 1:1")

	parser := Book()

	runeToStr := func(them []interface{}) string {
		sb := strings.Builder{}
		for _, it := range them {
			sb.WriteRune(it.(rune))
		}
		return sb.String()
	}

	t.Run("When parsing a regular book", func(t *testing.T) {
		result, ok := parser.Parse(john)

		if !ok {
			t.Errorf("Expected success")
		}

		if runeToStr(result.Data.([]interface{})) != "John" {
			t.Errorf("Expected 'John', but got %v", result.Data)
		}

		if string(result.Rest) != "1:1" {
			t.Errorf("Expected ' 1:1', but got %v", string(result.Rest))
		}
	})

	t.Run("When parsing a book that starts with a digit", func(t *testing.T) {
		result, ok := parser.Parse(firstJohn)

		if !ok {
			t.Errorf("Expected success")
		}

		if runeToStr(result.Data.([]interface{})) != "1 John" {
			t.Errorf("Expected '1 John', but got %v", result.Data)
		}

		if string(result.Rest) != "1:1" {
			t.Errorf("Expected1:1', but got %v", string(result.Rest))
		}
	})
}

func TestVerseRange(t *testing.T) {
	f := []rune("16-20 a")
	s := []rune("16 a")

	vp := VerseRange()

	t.Run("When parsing a verse range", func(t *testing.T) {

		result, ok := vp.Parse(f)

		if !ok {
			t.Errorf("Expected success")
		}

		if result.Data.(Verses).VerseNumber != 16 {
			t.Errorf("Expected Verse # to be 16, but got %v", result.Data.(Verses).VerseNumber)
		}

		if result.Data.(Verses).Count != 4 {
			t.Errorf("Expected Verse Count to be 16, but got %v", result.Data.(Verses).Count)
		}

		if string(result.Rest) != " a" {
			t.Errorf("Expected ' a', but got %v", string(result.Rest))
		}
	})

	t.Run("When parsing a verse", func(t *testing.T) {

		result, ok := vp.Parse(s)

		if !ok {
			t.Errorf("Expected success")
		}

		if result.Data.(Verses).VerseNumber != 16 {
			t.Errorf("Expected Verse # to be 3, but got %v", result.Data)
		}

		if string(result.Rest) != "a" {
			t.Errorf("Expected 'a', but got %v", string(result.Rest))
		}
	})
}
