//go:embed datasource.cue

package datasource

import (
	"embed"

	"github.com/grafana/grafana/pkg/schema"
	"github.com/grafana/thema"
	"github.com/grafana/thema/kernel"
)

const (
	cuePath     = "internal/components/datasource"
	cueFilename = "datasource.cue"
)

var (
	cueFS embed.FS
)

var (
	_ thema.LineageFactory = NewLineage
)

func NewLineage(lib thema.Library, opts ...thema.BindOption) (thema.Lineage, error) {
	lin, err := schema.LoadLineage(cuePath, cueFS, lib)
	if err != nil {
		return nil, err
	}

	// Calling this ensures our program cannot start if the Go DataSource.Model type
	// is not aligned with the canonical schema version in our lineage
	if _, err := NewJSONKernel(lin); err != nil {
		return nil, err
	}

	zsch, err := lin.Schema(thema.SV(0, 0))
	if err != nil {
		return nil, err
	}

	if err := thema.AssignableTo(zsch, Model{}); err != nil {
		return nil, err
	}

	return lin, nil
}

func NewJSONKernel(lin thema.Lineage) (kernel.InputKernel, error) {
	return kernel.NewInputKernel(kernel.InputKernelConfig{
		Lineage:     lin,
		Loader:      kernel.NewJSONDecoder(cueFilename),
		To:          thema.SV(0, 0),
		TypeFactory: func() interface{} { return &Model{} },
	})
}
