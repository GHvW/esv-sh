package lib

// completely overkill, but so so fun

type ParseSuccess struct {
	Data interface{}
	Rest []rune
}

// Parser is a function that takes a string and returns a value and a bool
type Parser interface {
	// []rune is the input string
	Parse([]rune) (ParseSuccess, bool)
}

// func satsify(predicate func(string) bool) bool {
// 	if predicate(c) {
// }
type Item struct {
	// component *Combinator
}

func (item Item) Parse(s []rune) (*ParseSuccess, bool) {
	if len(s) == 0 {
		return nil, false
	}
	return &ParseSuccess{s[0], s[1:]}, true
}

// Zero is a parser that always fails
type Zero struct {
	// component *Combinator
}

func (z Zero) Parse(s []rune) (*ParseSuccess, bool) {
	return nil, false
}

// Succeed is a parser that always succeeds and returns the given value.
type Succeed struct {
	item interface{}
	// component *Combinator
}

func SucceedWith(item interface{}) *Succeed {
	return &Succeed{item}
}

func (su *Succeed) Parse(s []rune) (*ParseSuccess, bool) {
	return &ParseSuccess{su.item, s}, true
}

// Map is a Parser that applies a function to the result of another Parser
// type Map struct {
// 	// a -> b
// 	fn        func(interface{}) interface{}
// 	component *Combinator
// }

// func (m *Map) Parse(s []rune) (interface{}, bool) {
// 	item, ok := m.component.Parse(s)
// 	if ok {
// 		return m.fn(item), true
// 	}
// 	return nil, false
// }

// // FlatMap is a Parser that applies a function that returns a Parser to the result of a Parser
// type FlatMap struct {
// 	// a -> Parser[b]
// 	fn        func(interface{}) Parser
// 	component *Combinator
// }

// func (fm *FlatMap) Parse(s []rune) (interface{}, bool) {
// 	item, ok := fm.component.Parse(s)
// 	if ok {
// 		return fm.fn(item).Parse(s)
// 	}
// 	return nil, false
// }

// type Combinator struct {
// 	parser Parser
// }

// func (pc *Combinator) Map(fn func(interface{}) interface{}) Parser {
// 	return &Map{fn, &Combinator{pc.parser}}
// }

// func (pc *Combinator) FlatMap(fn func(interface{}) Parser) Parser {
// 	return &FlatMap{fn, &Combinator{pc.parser}}
// }

// func (pc *Combinator) Or(other Parser) Parser {
// 	return &Or{pc.parser, other}
// }

// func (pc *Combinator) Parse(s []rune) (interface{}, bool) {
// 	return pc.parser.Parse(s)
// }

// type Satisfy struct {
// 	predicate func(rune) bool
// 	component Combinator
// }

// func (sat *Satisfy) Parse(a rune) (interface{}, bool) {
// 	if sat.predicate(a) {
// 		return &Succeed{a}, true
// 	}

// 	return &Zero{}, false
// }

// // Alpha is a parser that matches any single character in the range a-z
// type Alpha struct {
// 	component Combinator
// }

// func (a *Alpha) Parse(s rune) (interface{}, bool) {
// 	return &Satisfy{func(rune) bool {
// 		return s >= 'a' && s <= 'z'
// 	}}, true

// }

// // Digit matches a string that is a digit character
// type Digit struct {
// 	component Combinator
// }

// func (d *Digit) Parse(s rune) (interface{}, bool) {
// 	return &Satisfy{func(rune) bool {
// 		return s >= '0' && s <= '9'
// 	}}, true
// }

// type Or struct {
// 	first     Parser
// 	second    Parser
// 	component Combinator
// }

// func (or *Or) Parse(s []rune) (interface{}, bool) {
// 	item, ok := or.first.Parse(s)
// 	if ok {
// 		return item, true
// 	}
// 	return or.second.Parse(s)
// }

// type Many struct {
// 	parser    Parser
// 	component Combinator
// }

// func (many *Many) Parse(s []rune) (interface{}, bool) {
// 	// results := []interface{}{}
// 	// next := s
// 	// for {
// 	// 	item, ok := many.parser.Parse(s)
// 	// 	if ok {
// 	// 		results = append(results, item)
// 	// 		next = next[1:]
// 	// 	} else {
// 	// 		return results, true
// 	// 	}
// 	// }
// 	return many.component.FlatMap(func(item rune) Combinator {
// 		return Many{many.parser}.FlatMap(func(rest []rune) Combinator {
// 			return append(rest, item)
// 		})
// 	}).Or(&Succeed{[]interface{}{}}).Parse(s)
// }
