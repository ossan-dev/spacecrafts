package domain

type contextKey struct {
	Key int
}

var (
	ClientKey contextKey = contextKey{Key: 1}
	ModelsKey contextKey = contextKey{Key: 2}
)
