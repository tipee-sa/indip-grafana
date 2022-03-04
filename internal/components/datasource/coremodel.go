package datasource

import (
	"context"
	"errors"
	"fmt"
	"time"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/grafana/grafana/internal/cuectx"
	"github.com/grafana/grafana/internal/k8sbridge"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/schema"
	"github.com/grafana/thema"
)

// Store
//
// TODO: I think we should define a generic store interface similar to k8s rest.Interface
// and have storeset around (similar to clientset) from which we can grab specific store implementation for schema.
type Store interface {
	Get(ctx context.Context, uid string) (CR, error)
	Insert(ctx context.Context, ds CR) error
	Update(ctx context.Context, ds CR) error
	Delete(ctx context.Context, uid string) error
}

// Coremodel
type Coremodel struct {
	schema schema.ObjectSchema
	client rest.Interface
	store  Store
}

// ProvideDatasourceCoreModel
func ProvideDatasourceCoreModel(bridge k8sbridge.Service, store Store) (*Coremodel, error) {
	schema, ok := schema.LoadCoreSchema(schemaName)
	if !ok {
		return nil, fmt.Errorf("no schema registered for %s", schemaName)
	}

	client, err := bridge.Client().ClientForSchema(schema)
	if err != nil {
		return nil, err
	}

	return &Coremodel{
		client: client,
		store:  store,
		schema: schema,
	}, nil
}

// Schema
func (d *Coremodel) Schema() schema.ObjectSchema {
	return d.schema
}

// Reconcile
func (d *Coremodel) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	ds := CR{}

	err := d.client.Get().Namespace(req.Namespace).Resource("datasources").Name(req.Name).Do(ctx).Into(&ds)

	// TODO: check ACTUAL error
	if errors.Is(err, errNotFound) {
		return reconcile.Result{}, d.store.Delete(ctx, req.Name)
	}

	if err != nil {
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: 1 * time.Minute,
		}, err
	}

	_, err = d.store.Get(ctx, string(ds.UID))
	if err != nil {
		if !errors.Is(err, models.ErrDataSourceNotFound) {
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: 1 * time.Minute,
			}, err
		}

		if err := d.store.Insert(ctx, ds); err != nil {
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: 1 * time.Minute,
			}, err
		}
	}

	if err := d.store.Update(ctx, ds); err != nil {
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: 1 * time.Minute,
		}, err
	}

	return reconcile.Result{}, nil
}

func init() {
	lib := cuectx.ProvideThemaLibrary()
	lin, err := NewLineage(lib)
	if err != nil {
		panic(err)
	}

	// Calling this ensures our program cannot start if the Go DataSource.Model type
	// is not aligned with the canonical schema version in our lineage
	if _, err := NewJSONKernel(lin); err != nil {
		panic(err)
	}

	zsch, _ := lin.Schema(thema.SV(0, 0))
	if err = thema.AssignableTo(zsch, Model{}); err != nil {
		panic(err)
	}

	schema.RegisterCoreSchema(
		schema.NewThemaSchema(
			lin,
			groupName, groupVersion,
			openapiSchema,
			&CR{}, &CRList{},
		),
	)
}
