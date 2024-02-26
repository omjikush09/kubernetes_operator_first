/*
Copyright 2024 Omji Kushwaha.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/omjikush09/operator_learn/api/v1alpha1"
)

// OmjiReconciler reconciles a Omji object
type OmjiReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.omjikushwaha.com,resources=omjis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.omjikushwaha.com,resources=omjis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.omjikushwaha.com,resources=omjis/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Omji object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *OmjiReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	log.Log.Info("Reconciling Omji")

	// TODO(user): your logic here
	omji := &apiv1alpha1.Omji{}
	err := r.Get(ctx, req.NamespacedName, omji)
	if err != nil {
		log.Log.Error(err, "unable to fetch omji")
		return ctrl.Result{}, err
	}
	startTime := omji.Spec.Start
	endTime := omji.Spec.End
	replicas := omji.Spec.Replicas

	currentHour := time.Now().UTC().Hour()
	if currentHour >= startTime && currentHour <= endTime {
		// scale up
		for _, deploy := range omji.Spec.Deployments {

			deployment := &v1.Deployment{}
			err := r.Get(ctx, types.NamespacedName{Namespace: omji.Namespace, Name: deploy.Name}, deployment)
			if err != nil {
				log.Log.Error(err, "unable to fetch deployment")
				return ctrl.Result{}, err
			}
			if deployment.Spec.Replicas != &replicas {
				deployment.Spec.Replicas = &replicas
				err = r.Update(ctx, deployment)
				if err != nil {
					log.Log.Error(err, "unable to update deployment")
					return ctrl.Result{}, err
				}
			}

		}
	}
	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OmjiReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Omji{}).
		Complete(r)
}
