package sipparser

type URI interface {
	Scheme() string
	Parse(src []byte, pos int) (newPos int, err error)
}
