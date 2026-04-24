package config

import "github.com/sileader/llama-run/builder"

type ApplicationBuilderVisitor interface {
	Visit(builder builder.ApplicationBuilder) error
}

func visitAll(builder builder.ApplicationBuilder, visitor ...ApplicationBuilderVisitor) error {
	for _, v := range visitor {
		if err := v.Visit(builder); err != nil {
			return err
		}
	}
	return nil
}
