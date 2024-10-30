package order

type service struct {
}

type Opt struct{}

func New(opt Opt) *service {
	return &service{}
}
