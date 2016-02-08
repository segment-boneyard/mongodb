package schema

type HandlerFunc func()

type Service struct {
	Name    string `validate:"nonzero"`
	Version string `validate:"nonzero"`
	Config  Config
	Handler HandlerFunc
}
