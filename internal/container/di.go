package container

import (
	"fmt"

	"github.com/sarulabs/di/v2"
)

var (
	container di.Container
	defs      []di.Def
)

const App = di.App

func Add(def di.Def) {
	defs = append(defs, def)
}

func Build(scopes ...string) error {
	var (
		builder *di.Builder
		err     error
	)
	if builder, err = di.NewBuilder(scopes...); err != nil {
		return fmt.Errorf("can't create builder: %w", err)
	}

	if err = builder.Add(defs...); err != nil {
		return fmt.Errorf("can't add definitions: %w", err)
	}

	container = builder.Build()

	return nil
}

func Fill(name string, dst interface{}) error {
	return container.Fill(name, dst)
}
