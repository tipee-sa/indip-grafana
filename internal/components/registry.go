package components

import (
	"context"

	"github.com/grafana/grafana/pkg/schema"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Coremodel interface {
	Schema() schema.ObjectSchema
}

type ReconcilableCoremodel interface {
	Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error)
}

type Registry struct{}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(m Coremodel) bool {
	return false
}

func (r *Registry) Coremodels() []Coremodel {
	return nil
}

var coremodelRegistry = NewRegistry()

func ProvideReadOnlyCoremodelRegistry() *Registry {
	return coremodelRegistry
}

func RegisterCoremodel() {

}
