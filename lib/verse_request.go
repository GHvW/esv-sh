package lib

type Verse struct {
	Book    string
	Chapter int
	Verse   int
}

type VerseRequest struct {
	Verse Verse
	Count int
}

type VerseRange struct {
	VerseNumber int
	Count       int
}
