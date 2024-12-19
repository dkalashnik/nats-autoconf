package operator

import (
	"context"
	"fmt"
	"strings"

	"github.com/dkalashnik/nats-autoconf/pkg/apis/config/v1alpha1"
	"github.com/dkalashnik/nats-autoconf/pkg/nats"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	client.Client
	natsClient *nats.Client
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.Log.WithName("reconcile").WithValues("request", req.NamespacedName)

	config := &v1alpha1.ServiceConfiguration{}
	err := r.Get(ctx, req.NamespacedName, config)

	if err != nil {
		if errors.IsNotFound(err) {
			keys, err := r.natsClient.ListKeys()
			if err != nil {
				log.Error(err, "Failed to list NATS KV keys")
				return ctrl.Result{}, err
			}

			prefix := fmt.Sprintf("%s.%s.%s", req.Namespace, req.Name, "")
			for _, key := range keys {
				if strings.HasPrefix(key, prefix) {
					if err := r.natsClient.DeleteKey(key); err != nil {
						log.Error(err, "Failed to delete NATS KV key", "key", key)
						return ctrl.Result{}, err
					}
					log.Info("Deleted NATS KV key", "key", key)
				}
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	key := nats.GetConfigPath(
		config.Spec.Organization,
		config.Spec.Product,
		config.Spec.Version,
		config.Spec.ServiceName,
		config.Spec.ConfigName,
	)

	if err := r.natsClient.PutConfig(key, &config.Spec.Config); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to store config: %v", err)
	}

	log.Info("Successfully reconciled ServiceConfiguration", "key", key)
	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ServiceConfiguration{}).
		Complete(r)
}

func Run(natsClient *nats.Client) error {
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		return fmt.Errorf("unable to create manager: %v", err)
	}

	if err = (&Reconciler{
		Client:     mgr.GetClient(),
		natsClient: natsClient,
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("unable to create controller: %v", err)
	}

	return mgr.Start(ctrl.SetupSignalHandler())
}
