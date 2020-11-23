/*


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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	loadtestsv1beta1 "github.com/luizbafilho/lokust/apis/loadtests/v1beta1"
)

// LocustTestReconciler reconciles a LocustTest object
type LocustTestReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=loadtests.lokust.io,resources=locusttests,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=loadtests.lokust.io,resources=locusttests/status,verbs=get;update;patch

func (r *LocustTestReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("locusttest", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *LocustTestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loadtestsv1beta1.LocustTest{}).
		Complete(r)
}
