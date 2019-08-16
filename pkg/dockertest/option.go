package dockertest

import "github.com/ory/dockertest"

type OptionModifier func(p *dockertest.RunOptions)

func applyOptionsModifiers(o *dockertest.RunOptions, modifiers ...OptionModifier) {
	for _, modifier := range modifiers {
		modifier(o)
	}
}

// NoOpModifier is a no operation modifier which does nothing.
func NoOpModifier(_ *dockertest.RunOptions) {
}
