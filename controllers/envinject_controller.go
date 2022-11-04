/*
Copyright 2022.

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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webappv1 "github.com/evolutica-software-solutions/go-custom-k8s-operator/api/v1"
)

// EnvInjectReconciler reconciles a EnvInject object
type EnvInjectReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.evolutica.co.za,resources=envinjects,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.evolutica.co.za,resources=envinjects/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.evolutica.co.za,resources=envinjects/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EnvInject object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *EnvInjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling EnvInject custom resource...")
	// TODO(user): your logic here

	// Get EnvInject resource that triggered the recociliation request
	envInjector := &webappv1.EnvInject{}

	if err := r.Client.Get(ctx, req.NamespacedName, envInjector); err != nil {
		log.Error(err, "Unable to fetch EnvInject")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deployments := &appsv1.DeploymentList{}
	if err := r.Client.List(ctx, deployments, &client.ListOptions{Namespace: "podinfo-test"}); err != nil {
		log.Error(err, "Unable to fetch deployments")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploys := []appsv1.Deployment{}
	for _, v := range deployments.Items {
		fmt.Println(v.Name)
		if v.Name == "frontend-podinfo" || v.Name == "backend-podinfo" {
			fmt.Println(v.Spec.Template.Spec.Containers[0].Env)
			deploys = append(deploys, v)
		}
	}

	for _, v := range deploys {
		println("Deployment to update: ", v.Name)
		newVars := GetNewEnvs(envInjector.Spec)
		v.Spec.Template.Spec.Containers[0].Env = append(v.Spec.Template.Spec.Containers[0].Env, newVars)
		r.Client.Update(ctx, &v, &client.UpdateOptions{})
	}
	return ctrl.Result{}, nil
}

func GetNewEnvs(spec webappv1.EnvInjectSpec) corev1.EnvVar {
	return corev1.EnvVar{
		Name: "Name Testing",
		Value: "Value Testing",
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *EnvInjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.EnvInject{}).
		Complete(r)
}
