package test

import "testing"

type (
	Case interface {
		Name() string
		Setup() error
		Teardown() error
		Expected() Result
		Result() Result
		TestFunc() func(t *testing.T)
	}

	Result interface {
		Value() interface{} // TODO set a concrete type
		Error() error
		// TODO: Define additional properties
	}

	Cases struct {
		cases []Case
	}
)

func (tcs *Cases) Add(tc ...Case) {
	tcs.cases = append(tcs.cases, tc...)
}

func (tcs *Cases) Get(i int) Case {
	return tcs.cases[i]
}

func (tcs *Cases) All() []Case {
	return tcs.cases
}
