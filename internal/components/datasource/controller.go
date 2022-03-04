package datasource

import (
	"context"
	"errors"
	"time"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/grafana/grafana/internal/k8sbridge"
	"github.com/grafana/grafana/pkg/models"
)

type DatasourceReconciler struct {
	cli rest.Interface
	sto Store
}

// TODO: Looks like this should be the other way around,
// Otherwise the reconciler will never be registered,
// since there are no components that depend on DatasourceReconciler
// I think we need some kind of reconciler registry (maybe just as part of the k8s bridge),
// similar to background services registry.
func ProvideDatasourceReconciler(bridge *k8sbridge.Service, store Store) (*DatasourceReconciler, error) {
	// d := &DatasourceReconciler{
	// 	cli: bridge.Client().GrafanaCoreV1(),
	// }

	// // TODO should Thema-based approaches differ from pure k8s here? (research!)
	// if err := ctrl.NewControllerManagedBy(bridge.ControllerManager()).
	// 	Named("datasource-controller").
	// 	For(&CR{}).
	// 	Complete(reconcile.Func(d.Reconcile)); err != nil {
	// 	return nil, err
	// }

	// return d, nil
	return nil, nil
}

var errNotFound = errors.New("k8s obj not found")

func (d *DatasourceReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	ds := CR{}

	err := d.cli.Get().Namespace(req.Namespace).Resource("datasources").Name(req.Name).Do(ctx).Into(&ds)

	// TODO: check ACTUAL error
	if errors.Is(err, errNotFound) {
		return reconcile.Result{}, d.sto.Delete(ctx, req.Name)
	}

	if err != nil {
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: 1 * time.Minute,
		}, err
	}

	_, err = d.sto.Get(ctx, string(ds.UID))
	if err != nil {
		if !errors.Is(err, models.ErrDataSourceNotFound) {
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: 1 * time.Minute,
			}, err
		}

		if err := d.sto.Insert(ctx, ds); err != nil {
			return reconcile.Result{
				Requeue:      true,
				RequeueAfter: 1 * time.Minute,
			}, err
		}
	}

	if err := d.sto.Update(ctx, ds); err != nil {
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: 1 * time.Minute,
		}, err
	}

	return reconcile.Result{}, nil
}
