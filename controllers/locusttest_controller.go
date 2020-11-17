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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	testsv1beta1 "github.com/luizbafilho/lokust/api/v1beta1"
)

const (
	ComponentMaster = "master"
	ComponentWorker = "worker"
)

// LocustTestReconciler reconciles a LocustTest object
type LocustTestReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tests.lokust.io,resources=locusttests,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tests.lokust.io,resources=locusttests/status,verbs=get;update;patch

func (r *LocustTestReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("locusttest", req.NamespacedName)

	var test testsv1beta1.LocustTest
	if err := r.Get(ctx, req.NamespacedName, &test); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner("locusttest-controller")}

	masterDeploy, err := r.desiredMasterDeployment(test)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Patch(ctx, &masterDeploy, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}
	masterSvc, err := r.desiredMasterService(test)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Patch(ctx, &masterSvc, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	workerDeploy, err := r.desiredWorkerDeployment(test)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Patch(ctx, &workerDeploy, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	// err = r.Status().Update(ctx, &book)
	// if err != nil {
	// 	return ctrl.Result{}, err
	// }

	log.Info("reconciled locust test")
	return ctrl.Result{}, nil
}

var (
	ownerKey = ".metadata.controller"
	apiGVStr = testsv1beta1.GroupVersion.String()
)

func (r *LocustTestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(&appsv1.Deployment{}, ownerKey, func(rawObj runtime.Object) []string {
		// grab the deployment object, extract the owner...
		dep := rawObj.(*appsv1.Deployment)
		owner := metav1.GetControllerOf(dep)
		if owner == nil {
			return nil
		}
		// ...make sure it's a SeldonDeployment...
		if owner.APIVersion != apiGVStr || owner.Kind != "LocustTest" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&testsv1beta1.LocustTest{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

func (r *LocustTestReconciler) desiredMasterService(test testsv1beta1.LocustTest) (corev1.Service, error) {
	svc := corev1.Service{
		TypeMeta: metav1.TypeMeta{APIVersion: corev1.SchemeGroupVersion.String(), Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      test.Name,
			Namespace: test.Namespace,
			// OwnerReferences: []metav1.OwnerReference{*controller.NewProxyOwnerRef(p)},
		},
		Spec: corev1.ServiceSpec{
			Selector: makeLabels(test, ComponentMaster),
			Type:     "ClusterIP",
			Ports: []corev1.ServicePort{
				{
					Name:     "loc-master-web",
					Protocol: "TCP",
					Port:     8089,
					TargetPort: intstr.IntOrString{
						IntVal: 8089,
					},
				},
				{
					Name:     "loc-master-p1",
					Protocol: "TCP",
					Port:     5557,
					TargetPort: intstr.IntOrString{
						IntVal: 5557,
					},
				},
				{
					Name:     "loc-master-p2",
					Protocol: "TCP",
					Port:     5558,
					TargetPort: intstr.IntOrString{
						IntVal: 5558,
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(&test, &svc, r.Scheme); err != nil {
		return svc, err
	}

	return svc, nil
}

func (r *LocustTestReconciler) desiredMasterDeployment(test testsv1beta1.LocustTest) (appsv1.Deployment, error) {
	replicas := int32(1)
	depl := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: appsv1.SchemeGroupVersion.String(), Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      test.Name + "-master",
			Namespace: test.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas, // won't be nil because defaulting
			Selector: &metav1.LabelSelector{
				MatchLabels: makeLabels(test, ComponentMaster),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: makeLabels(test, ComponentMaster),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "locust-master",
							Image: "locustio/locust:1.3.2",
							Env: []corev1.EnvVar{
								{Name: "HOST", Value: "https://www.google.com"},
							},
							Args: []string{"--master", "-f", "/mnt/locust/locustfile.py"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "locustfile",
									ReadOnly:  true,
									MountPath: "/mnt/locust",
								},
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8089, Name: "loc-master-web", Protocol: "TCP"},
								{ContainerPort: 5557, Name: "loc-master-p1", Protocol: "TCP"},
								{ContainerPort: 5558, Name: "loc-master-p2", Protocol: "TCP"},
							},
							Resources: *test.Spec.Resources.DeepCopy(),
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "locustfile",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "locustfile-configmap",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(&test, &depl, r.Scheme); err != nil {
		return depl, err
	}

	return depl, nil
}

func (r *LocustTestReconciler) desiredWorkerDeployment(test testsv1beta1.LocustTest) (appsv1.Deployment, error) {
	depl := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: appsv1.SchemeGroupVersion.String(), Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      test.Name + "-worker",
			Namespace: test.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: test.Spec.Replicas, // won't be nil because defaulting
			Selector: &metav1.LabelSelector{
				MatchLabels: makeLabels(test, ComponentWorker),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: makeLabels(test, ComponentWorker),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "locust-worker",
							Image: "locustio/locust:1.3.2",
							Env: []corev1.EnvVar{
								{Name: "HOST", Value: "https://www.google.com"},
							},
							Args: []string{
								"--worker",
								"--master-host", test.Name,
								"-f", "/mnt/locust/locustfile.py",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "locustfile",
									ReadOnly:  true,
									MountPath: "/mnt/locust",
								},
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 8089, Name: "loc-master-web", Protocol: "TCP"},
								{ContainerPort: 5557, Name: "loc-master-p1", Protocol: "TCP"},
								{ContainerPort: 5558, Name: "loc-master-p2", Protocol: "TCP"},
							},
							Resources: *test.Spec.Resources.DeepCopy(),
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "locustfile",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "locustfile-configmap",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(&test, &depl, r.Scheme); err != nil {
		return depl, err
	}

	return depl, nil
}

func makeLabels(test testsv1beta1.LocustTest, component string) map[string]string {
	return map[string]string{
		"locust-test":      test.Name,
		"locust-component": component,
	}
}
