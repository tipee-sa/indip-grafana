package schema

import (
	"io"
	"io/fs"
	"path/filepath"
	"testing/fstest"

	"github.com/grafana/thema"
	"github.com/grafana/thema/load"
)

// LoadLineage
func LoadLineage(path string, cueFS fs.FS, lib thema.Library, opts ...thema.BindOption) (thema.Lineage, error) {
	prefix := filepath.FromSlash(path)
	fs, err := prefixWithGrafanaCUE(prefix, cueFS)
	if err != nil {
		return nil, err
	}
	inst, err := load.InstancesWithThema(fs, prefix)

	// Need to trick loading by creating the embedded file and
	// making it look like a module in the root dir.
	if err != nil {
		return nil, err
	}

	val := lib.Context().BuildInstance(inst)

	lin, err := thema.BindLineage(val, lib, opts...)
	if err != nil {
		return nil, err
	}

	return lin, nil
}

func prefixWithGrafanaCUE(prefix string, inputfs fs.FS) (fs.FS, error) {
	m := fstest.MapFS{
		filepath.Join("cue.mod", "module.cue"): &fstest.MapFile{Data: []byte(`module: "github.com/grafana/grafana"`)},
	}

	prefix = filepath.FromSlash(prefix)
	err := fs.WalkDir(inputfs, ".", (func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := inputfs.Open(path)
		if err != nil {
			return err
		}

		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		m[filepath.Join(prefix, path)] = &fstest.MapFile{Data: []byte(b)}
		return nil
	}))

	return m, err
}
