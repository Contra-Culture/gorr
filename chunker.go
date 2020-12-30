package gorr

import (
	"errors"
	"net/url"
	"strings"
)

type (
	Chunker struct {
		chunks []string
		params map[string]string
		idx    int
	}
)

var (
	emptyString          string
	paramAlreadySetError = errors.New("param already set")
)

func NewChunker(u *url.URL) *Chunker {
	chunks := strings.Split(u.String(), "/")
	if len(chunks) == 2 && chunks[0] == "" && chunks[0] == chunks[1] {
		chunks = []string{""}
	}
	return &Chunker{
		chunks: chunks,
		params: map[string]string{},
		idx:    0,
	}
}
func (ch *Chunker) Next() (ok bool) {
	ok = len(ch.chunks) > (ch.idx + 1)
	if !ok {
		return
	}
	ch.idx = ch.idx + 1
	return
}
func (ch *Chunker) Chunk() string {
	return ch.chunks[ch.idx]
}
func (ch *Chunker) Set(n, v string) (err error) {
	if ch.params[n] != emptyString {
		err = paramAlreadySetError
		return
	}
	ch.params[n] = v
	return
}
func (ch *Chunker) Params() map[string]string {
	return ch.params
}
