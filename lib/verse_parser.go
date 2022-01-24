package lib

import "unicode"

// completely overkill, but so so fun

type ParseSuccess struct {
	Data interface{}
	Rest []rune
}

// Parser is a function that takes a string and returns a value and a bool
type Parser interface {
	// []rune is the input string
	Parse([]rune) (*ParseSuccess, bool)
}

// func satsify(predicate func(string) bool) bool {
// 	if predicate(c) {
// }
type Item struct{}

func TakeItem() *Combinator {
	return &Combinator{Item{}}
}

func (item Item) Parse(s []rune) (*ParseSuccess, bool) {
	if len(s) == 0 {
		return nil, false
	}
	return &ParseSuccess{s[0], s[1:]}, true
}

// Zero is a parser that always fails
type Zero struct{}

func Fail() *Combinator {
	return &Combinator{Zero{}}
}

func (z Zero) Parse(s []rune) (*ParseSuccess, bool) {
	return nil, false
}

// Succeed is a parser that always succeeds and returns the given value.
type Succeed struct {
	item interface{}
}

func Success(item interface{}) *Combinator {
	return &Combinator{&Succeed{item}}
}

func (su *Succeed) Parse(s []rune) (*ParseSuccess, bool) {
	return &ParseSuccess{su.item, s}, true
}

// Map is a Parser that applies a function to the result of another Parser
type Map struct {
	// a -> b
	fn     func(interface{}) interface{}
	parser Parser
}

func (m *Map) Parse(s []rune) (*ParseSuccess, bool) {
	result, ok := m.parser.Parse(s)
	if !ok {
		return nil, false
	}
	return &ParseSuccess{m.fn(result.Data), result.Rest}, true
}

// FlatMap is a Parser that applies a function that returns a Parser to the result of a Parser
type FlatMap struct {
	// a -> Parser[b]
	fn     func(interface{}) Parser
	parser Parser
}

func (fm *FlatMap) Parse(s []rune) (*ParseSuccess, bool) {
	item, ok := fm.parser.Parse(s)
	if !ok {
		return nil, false
	}
	return fm.fn(item.Data).Parse(item.Rest)
}

type Or struct {
	first  Parser
	second Parser
}

func (or *Or) Parse(s []rune) (*ParseSuccess, bool) {
	item, ok := or.first.Parse(s)
	if ok {
		return item, true
	}
	return or.second.Parse(s)
}

type Combinator struct {
	parser Parser
}

func (pc *Combinator) Map(fn func(interface{}) interface{}) *Combinator {
	return &Combinator{&Map{fn, pc.parser}}
}

func (pc *Combinator) FlatMap(fn func(interface{}) Parser) *Combinator {
	return &Combinator{&FlatMap{fn, pc.parser}}
}

func (pc *Combinator) Or(other Parser) *Combinator {
	return &Combinator{&Or{pc.parser, other}}
}

func (pc *Combinator) Parse(s []rune) (*ParseSuccess, bool) {
	return pc.parser.Parse(s)
}

type Satisfy struct {
	predicate func(rune) bool
}

func Satisfies(predicate func(rune) bool) *Combinator {
	return &Combinator{&Satisfy{predicate}}
}

func (sat *Satisfy) Parse(s []rune) (*ParseSuccess, bool) {
	return TakeItem().FlatMap(func(it interface{}) Parser {
		if sat.predicate(it.(rune)) {
			return Success(it)
		}
		return Fail()
	}).Parse(s)
}

// Alpha is a parser that matches any single character in the range a-z
type Alpha struct{}

func IsAlpha() *Combinator {
	return &Combinator{Alpha{}}
}

func (a Alpha) Parse(s []rune) (*ParseSuccess, bool) {
	return Satisfies(unicode.IsLetter).Parse(s)
}

// Digit matches a string that is a digit character
type Digit struct{}

func IsDigit() *Combinator {
	return &Combinator{Digit{}}
}

func (d Digit) Parse(s []rune) (*ParseSuccess, bool) {
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
