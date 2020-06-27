/*
Copyright 2020 disp.
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

package oilappservice

import (
	"context"
	"encoding/json"
	"reflect"

	oilappv1 "hmxq.top/oilapp-operator/pkg/apis/oilapp/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	//"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	//"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

var log = logf.Log.WithName("controller_oilappservice")

// Add creates a new OilappService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileOilappService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("oilappservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource OilappService
	err = c.Watch(&source.Kind{Type: &oilappv1.OilappService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner OilappService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &oilappv1.OilappService{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileOilappService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileOilappService{}

// ReconcileOilappService reconciles a OilappService object
type ReconcileOilappService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a OilappService object and makes changes based on the state read
// and what is in the OilappService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileOilappService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling OilappService")

	reqLogger.Info("获取OilappService对象实例")
	instance := &oilappv1.OilappService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("OilappService对象已被删除")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	/*
		// Define a new Pod object
		pod := newPodForCR(instance)

		// Set OilappService instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		// Check if this Pod already exists
		found := &corev1.Pod{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
			err = r.client.Create(context.TODO(), pod)
			if err != nil {
				return reconcile.Result{}, err
			}

			// Pod created successfully - don't requeue
			return reconcile.Result{}, nil
		} else if err != nil {
			return reconcile.Result{}, err
		}

		// Pod already exists - don't requeue
		reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
		return reconcile.Result{}, nil

	*/

	reqLogger.Info("获取Deployment对象进行检查")
	deploy := &appsv1.Deployment{}
	if err := r.client.Get(context.TODO(), request.NamespacedName, deploy); err != nil && errors.IsNotFound(err) {
		reqLogger.Info("对象Deployment不存在，开始创建关联资源")
		// 创建 Deploy
		deploy := CreateOilappDeploy(instance)
		if err := r.client.Create(context.TODO(), deploy); err != nil {
			return reconcile.Result{}, err
		}

		// 关联 Annotations
		data, _ := json.Marshal(instance.Spec)
		if instance.Annotations != nil {
			instance.Annotations["spec"] = string(data)
		} else {
			instance.Annotations = map[string]string{"spec": string(data)}
		}
		if err := r.client.Update(context.TODO(), instance); err != nil {
			return reconcile.Result{}, nil
		}

		// 创建 Service
		service := CreateOilappService(instance)
		if err := r.client.Create(context.TODO(), service); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	reqLogger.Info("对象Deployment已存在，开始检查是否需要更新")
	oldspec := oilappv1.OilappServiceSpec{}
	if err := json.Unmarshal([]byte(instance.Annotations["spec"]), &oldspec); err != nil {
		return reconcile.Result{}, err
	}
	//通过保存的spec，判断是否需要更新
	if !reflect.DeepEqual(instance.Spec, oldspec) {
		reqLogger.Info("数据已改变，需要更新相关资源")
		// 更新Deploy
		newDeploy := CreateOilappDeploy(instance)
		oldDeploy := &appsv1.Deployment{}
		if err := r.client.Get(context.TODO(), request.NamespacedName, oldDeploy); err != nil {
			return reconcile.Result{}, err
		}
		oldDeploy.Spec = newDeploy.Spec
		if err := r.client.Update(context.TODO(), oldDeploy); err != nil {
			return reconcile.Result{}, err
		}
		//更新Service
		newService := CreateOilappService(instance)
		oldService := &corev1.Service{}
		if err := r.client.Get(context.TODO(), request.NamespacedName, oldService); err != nil {
			return reconcile.Result{}, err
		}

		clusterip := oldService.Spec.ClusterIP
		oldService.Spec = newService.Spec
		oldService.Spec.ClusterIP = clusterip
		if err := r.client.Update(context.TODO(), oldService); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	reqLogger.Info("对象Deployment已存在，且无需更新")
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *oilappv1.OilappService) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

func CreateOilappDeploy(cr *oilappv1.OilappService) *appsv1.Deployment {
	/*
		apiVersion: apps/v1
		kind: Deployment
		metadata:
		  name: example-oilappservice
		  namespace: default
		  ownerReferences:
		  - apiVersion: oilapp.hmxq.top/v1
			blockOwnerDeletion: true
			controller: true
			kind: OilappService
			name: example-oilappservice
		spec:
		  progressDeadlineSeconds: 600
		  replicas: 2
		  revisionHistoryLimit: 10
		  selector:
			matchLabels:
			  app: example-oilappservice
		  strategy:
			rollingUpdate:
			  maxSurge: 25%
			  maxUnavailable: 25%
			type: RollingUpdate
		  template:
			metadata:
			  labels:
				app: example-oilappservice
			spec:
			  containers:
			  - image: tomcat:8.5-jdk11-openjdk
				imagePullPolicy: IfNotPresent
				name: example-oilappservice
				ports:
				- containerPort: 8080
				  protocol: TCP
			  dnsPolicy: ClusterFirst
			  restartPolicy: Always
	*/

	labels := map[string]string{"app": cr.Name}
	selector := &metav1.LabelSelector{MatchLabels: labels}
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,

			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   oilappv1.SchemeGroupVersion.Group,
					Version: oilappv1.SchemeGroupVersion.Version,
					Kind:    "OilappService",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: cr.Spec.Size,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: createOilappContainers(cr),
				},
			},

			Selector: selector,
		},
	}
}

func createOilappContainers(cr *oilappv1.OilappService) []corev1.Container {
	containerPorts := []corev1.ContainerPort{}
	for _, svcPort := range cr.Spec.Ports {
		cport := corev1.ContainerPort{}
		cport.ContainerPort = svcPort.TargetPort.IntVal
		containerPorts = append(containerPorts, cport)
	}
	return []corev1.Container{
		{
			Name:            cr.Name,
			Image:           cr.Spec.Image,
			Resources:       cr.Spec.Resources,
			Ports:           containerPorts,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Env:             cr.Spec.Envs,
		},
	}
}

func CreateOilappService(cr *oilappv1.OilappService) *corev1.Service {
	/*
		apiVersion: v1
		kind: Service
		metadata:
		  name: example-oilappservice
		  namespace: default
		  ownerReferences:
		  - apiVersion: oilapp.hmxq.top/v1
			blockOwnerDeletion: true
			controller: true
			kind: OilappService
			name: example-oilappservice
		spec:
		  externalTrafficPolicy: Cluster
		  ports:
		  - nodePort: 32001
			port: 8080
			protocol: TCP
			targetPort: 8080
		  selector:
			app: example-oilappservice
		  type: NodePort
	*/

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   oilappv1.SchemeGroupVersion.Group,
					Version: oilappv1.SchemeGroupVersion.Version,
					Kind:    "OilappService",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:  corev1.ServiceTypeNodePort,
			Ports: cr.Spec.Ports,
			Selector: map[string]string{
				"app": cr.Name,
			},
		},
	}
}
