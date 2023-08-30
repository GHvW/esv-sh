package lib

import (
	"strings"
	"unicode"
)

// completely overkill, but so so fun

type ParseSuccess[A any] struct {
	Data A
	Rest []rune
}

// Parser is a function that takes a string and returns a value and a bool
type Parser[A any] interface {
	// []rune is the input string
	Parse([]rune) (*ParseSuccess[A], bool)
}

// Zero is a parser that always fails
type Zero[A any] struct{}

func (z Zero[A]) Parse(s []rune) (*ParseSuccess[A], bool) {
	return nil, false
}

func Fail[A any]() *Zero[A] {
	return &Zero[A]{}
}

// Succeed is a parser that always succeeds and returns the given value.
type Succeed[A any] struct {
	item A
}

func Success[A any](item A) *Succeed[A] {
	return &Succeed[A]{item}
}

func (su *Succeed[A]) Parse(s []rune) (*ParseSuccess[A], bool) {
	return &ParseSuccess[A]{su.item, s}, true
}

// Map is a Parser that applies a function to the result of another Parser
type Map[A, B any] struct {
	// a -> b
	fn     func(A) B
	parser Parser[A]
}

func (m *Map[A, B]) Parse(s []rune) (*ParseSuccess[B], bool) {
	result, ok := m.parser.Parse(s)
	if !ok {
		return nil, false
	}
	return &ParseSuccess[B]{m.fn(result.Data), result.Rest}, true
}

// FlatMap is a Parser that applies a function that returns a Parser to the result of a Parser
type FlatMap[A, B any] struct {
	// a -> Parser[b]
	fn     func(A) Parser[B]
	parser Parser[A]
}

func (fm *FlatMap[A, B]) Parse(s []rune) (*ParseSuccess[B], bool) {
	item, ok := fm.parser.Parse(s)
	if !ok {
		return nil, false
	}
	return fm.fn(item.Data).Parse(item.Rest)
}

type Or[A any] struct {
	first  Parser[A]
	second Parser[A]
}

func (or *Or[A]) Parse(s []rune) (*ParseSuccess[A], bool) {
	item, ok := or.first.Parse(s)
	if ok {
		return item, true
	}
	return or.second.Parse(s)
}

type And[A, B any] struct {
	first  *Parser[A]
	second *Parser[B]
}

func (a *And[A, B]) Parse(s []rune) (*ParseSuccess[Pair[A, B]], bool) {
	item1, ok1 := a.first.Parse(s)

}

type Satisfy struct {
	predicate func(rune) bool
}

func Satisfies(predicate func(rune) bool) *Parser[rune] {
	return &Satisfy{predicate}
}

func (sat *Satisfy) Parse(s []rune) (*ParseSuccess[rune], bool) {
	if len(s) == 0 {
		return nil, false
	}

	data := s[0]
	if sat.predicate(data) {
		return &ParseSuccess[rune]{data, s[1:]}, true
	}
	return nil, false
}

// Alpha is a parser that matches any single character in the range a-z
type AlphaParser struct{}

func Alpha() *AlphaParser {
	return &AlphaParser{}
}

func (a AlphaParser) Parse(s []rune) (*ParseSuccess[rune], bool) {
	return Satisfies(unicode.IsLetter).Parse(s)
}

// Digit matches a string that is a digit character
type DigitParser struct{}

func Digit() *DigitParser {
	return &DigitParser{}
}

func (d DigitParser) Parse(s []rune) (*ParseSuccess[rune], bool) {
	return Satisfies(unicode.IsDigit).Parse(s)
}

type Many struct {
	parser *Combinator
}

func Multiple(parser *Combinator) *Combinator {
	return &Combinator{&Many{parser}}
}

func (many *Many) Parse(s []rune) (*ParseSuccess, bool) {
	return many.parser.FlatMap(func(it interface{}) Parser {
		return Multiple(many.parser).FlatMap(func(rest interface{}) Parser {
			// TODO - change this later to something more appropriate
			return Success(append([]interface{}{it}, rest.([]interface{})...))
		})
	}).Or(Success([]interface{}{})).Parse(s)
}

type Many1 struct {
	parser *Combinator
}

func AtLeastOne(parser *Combinator) *Combinator {
	return &Combinator{&Many1{parser}}
}

func (many1 *Many1) Parse(s []rune) (*ParseSuccess, bool) {
	return many1.parser.FlatMap(func(it interface{}) Parser {
		return Multiple(many1.parser).FlatMap(func(rest interface{}) Parser {
			// TODO - change this later to something more appropriate
			return Success(append([]interface{}{it}, rest.([]interface{})...))
		})
	}).Parse(s)
}

type Natural struct{}

func NaturalNumber() *Combinator {
	return &Combinator{&Natural{}}
}

func (nat *Natural) Parse(s []rune) (*ParseSuccess, bool) {
	return AtLeastOne(Digit()).Map(func(it interface{}) interface{} {
		result := 0
		for _, next := range it.([]interface{}) {
			result = (10 * result) + int(next.(rune)-'0')
		}
		return result
	}).Parse(s)
}

type RuneParser struct {
	r rune
}

func Rune(r rune) *Combinator {
	return &Combinator{RuneParser{r}}
}

func (rp RuneParser) Parse(s []rune) (*ParseSuccess, bool) {
	return Satisfies(func(r rune) bool {
		return r == rp.r
	}).Parse(s)
}

type WordParser struct{}

func Word() *Combinator {
	return &Combinator{WordParser{}}
}

func (wp WordParser) Parse(s []rune) (*ParseSuccess, bool) {
	return AtLeastOne(Alpha()).Parse(s)
}

type Space struct{}

func WhiteSpace() *Combinator {
	return &Combinator{Space{}}
}

func (space Space) Parse(s []rune) (*ParseSuccess, bool) {
	return Satisfies(unicode.IsSpace).Parse(s)
}

// Token
type TokenParser struct {
	parser *Combinator
}

func Token(parser *Combinator) *Combinator {
	return &Combinator{&TokenParser{parser}}
}

func (tok *TokenParser) Parse(s []rune) (*ParseSuccess, bool) {
	return tok.parser.FlatMap(func(token interface{}) Parser {
		return Multiple(WhiteSpace()).FlatMap(func(_spaces interface{}) Parser {
			return Success(token)
		})
	}).Parse(s)
}

type BookParser struct{}

func Book() *Combinator {
	return &Combinator{BookParser{}}
}

func (bp BookParser) Parse(s []rune) (*ParseSuccess, bool) {
	return Token(Rune('1').Or(Rune('2')).FlatMap(func(n interface{}) Parser {
		return WhiteSpace().FlatMap(func(_space interface{}) Parser {
			return Word().FlatMap(func(word interface{}) Parser {
				return Success(append([]interface{}{n, ' '}, word.([]interface{})...))
			})
		})
	})).Or(Token(Word())).Parse(s)
}

type VerseRangeParser struct{}

func VerseRange() *Combinator {
	return &Combinator{VerseRangeParser{}}
}

func (vrp VerseRangeParser) Parse(s []rune) (*ParseSuccess, bool) {
	return Token(NaturalNumber().FlatMap(func(verse interface{}) Parser {
		return Token(Rune('-')).FlatMap(func(_dash interface{}) Parser {
			return NaturalNumber().FlatMap(func(to interface{}) Parser {
				return Success([]interface{}{verse, to})
			})
		})
	})).Map(func(ns interface{}) interface{} {
		numbers := ns.([]interface{})
		return Verses{numbers[0].(int), numbers[1].(int) - numbers[0].(int) + 1}
	}).Or(Token(NaturalNumber()).Map(func(n interface{}) interface{} {
		return Verses{n.(int), 0}
	})).Parse(s)
}

type VerseRequestParser struct{}

func VerseReq() *Combinator {
	return &Combinator{VerseRequestParser{}}
}

func runeToStr(them []interface{}) string {
	sb := strings.Builder{}
	for _, it := range them {
		sb.WriteRune(it.(rune))
	}
	return sb.String()
}

func (vrp VerseRequestParser) Parse(s []rune) (*ParseSuccess, bool) {
	return Token(Book()).And(Token(NaturalNumber()).IgnoreNext(Rune(':')).And(VerseRange())).Map(func(data interface{}) interface{} {
		book := runeToStr(data.(*Pair).First.([]interface{}))
		chapterAndVerseRange := data.(*Pair).Second.(*Pair)
		chapter := chapterAndVerseRange.First.(int)
		verseRange := chapterAndVerseRange.Second.(Verses)
		return &VerseRequest{Verse{book, chapter, verseRange.VerseNumber}, verseRange.Count}
	}).Parse(s)
}
