/*
Copyright 2017 The Kubernetes Authors.

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

package main

import (
	"fmt"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	samplev1beta1 "github.com/fanfengqiang/cert-controller/pkg/apis/certcontroller/v1beta1"
	clientset "github.com/fanfengqiang/cert-controller/pkg/generated/clientset/versioned"
	samplescheme "github.com/fanfengqiang/cert-controller/pkg/generated/clientset/versioned/scheme"
	informers "github.com/fanfengqiang/cert-controller/pkg/generated/informers/externalversions/certcontroller/v1beta1"
	listers "github.com/fanfengqiang/cert-controller/pkg/generated/listers/certcontroller/v1beta1"
)

const controllerAgentName = "cert-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Cert is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Cert fails
	// to sync due to a Secret of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Secret already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Cert"
	// MessageResourceSynced is the message used for an Event fired when a Cert
	// is synced successfully
	MessageResourceSynced = "Cert synced successfully"
)

// Controller is the controller implementation for Cert resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	sampleclientset clientset.Interface

	secretsLister corelisters.SecretLister
	secretsSynced cache.InformerSynced
	certsLister        listers.CertLister
	certsSynced        cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new cert controller
func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	secretInformer coreinformers.SecretInformer,
	certInformer informers.CertInformer) *Controller {

	// Create event broadcaster
	// Add cert-controller types to the default Kubernetes Scheme so Events can be
	// logged for cert-controller types.
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:     kubeclientset,
		sampleclientset:   sampleclientset,
		secretsLister: secretInformer.Lister(),
		secretsSynced: secretInformer.Informer().HasSynced,
		certsLister:        certInformer.Lister(),
		certsSynced:        certInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Certs"),
		recorder:          recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Cert resources change
	certInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueCert,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueCert(new)
		},
	})
	// Set up an event handler for when Secret resources change. This
	// handler will lookup the owner of the given Secret, and if it is
	// owned by a Cert resource will enqueue that Cert resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Secret resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	secretInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*corev1.Secret)
			oldDepl := old.(*corev1.Secret)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Secrets.
				// Two different versions of the same Secret will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Cert controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.secretsSynced, c.certsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Cert resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Cert resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Cert resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Cert resource with this namespace/name
	cert, err := c.certsLister.Certs(namespace).Get(name)
	if err != nil {
		// The Cert resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("cert '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	secretName := cert.Spec.SecretName
	if secretName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: secret name must be specified", key))
		return nil
	}

	// Get the secret with the name specified in Cert.spec
	secret, err := c.secretsLister.Secrets(cert.Namespace).Get(secretName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		secret, err = c.kubeclientset.CoreV1().Secrets(cert.Namespace).Create(newSecret(cert))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the Secret is not controlled by this Cert resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(secret, cert) {
		msg := fmt.Sprintf(MessageResourceExists, secret.Name)
		c.recorder.Event(cert, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If this number of the replicas on the Cert resource is specified, and the
	// number does not equal the current desired replicas on the Secret, we
	// should update the Secret resource.
	now := time.Now()

	loc, _ := time.LoadLocation("Asia/Shanghai")
	formatTime,_:=time.ParseInLocation("2006-01-02-15-04-05",secret.Annotations["updateTime"],loc)

	s :=now.Sub(formatTime).Hours()/24
	t := float64(cert.Spec.ValidityPeriod)
	remain := int(t-s)
	if remain > 0 {
		klog.V(4).Infof("Cert %s ;updateTime: %s", name, secret.GetCreationTimestamp().Format("2006-01-02 15:04:05"))
		secret, err = c.kubeclientset.CoreV1().Secrets(cert.Namespace).Update(newSecret(cert))
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// Finally, we update the status block of the Cert resource to reflect the
	// current state of the world
	err = c.updateCertStatus(cert, secret,remain)
	if err != nil {
		return err
	}

	c.recorder.Event(cert, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) updateCertStatus(cert *samplev1beta1.Cert, secret *corev1.Secret,remain int) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	certCopy := cert.DeepCopy()
	certCopy.Status.SecretCreateTime = secret.Annotations["updateTime"]
	certCopy.Status.RemainingValidDays = remain
	//certCopy.Status.RemainingValidDays = "This feature is temporarily unavailable"
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the Cert resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.sampleclientset.CertcontrollerV1beta1().Certs(cert.Namespace).Update(certCopy)
	return err
}

// enqueueCert takes a Cert resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Cert.
func (c *Controller) enqueueCert(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Cert resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Cert resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Cert, we should not do anything more
		// with it.
		if ownerRef.Kind != "Cert" {
			return
		}

		cert, err := c.certsLister.Certs(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of cert '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueCert(cert)
		return
	}
}

// newSecret creates a new Secret for a Cert resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Cert resource that 'owns' it.
func newSecret(cert *samplev1beta1.Cert) *corev1.Secret {
	labels := map[string]string{
		"certificateAuthority":		"letsencrypt",
		"controller": 				cert.Name,
		"updateTime":				time.Unix(time.Now().Unix(), 0).Format("2006-01-02-15-04-05"),
	}
	log.Println("beigin to create a secret.........")
	var newCert = CreateCert(cert.Spec.Domain, cert.Spec.Type, cert.Spec.Env)
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cert.Spec.SecretName,
			Namespace: cert.Namespace,
			Annotations: labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cert, schema.GroupVersionKind{
					Group:   samplev1beta1.SchemeGroupVersion.Group,
					Version: samplev1beta1.SchemeGroupVersion.Version,
					Kind:    "Cert",
				}),
			},
		},

		Data: map[string][]byte{
			"tls.crt": []byte(newCert.cert),
			"tls.key": []byte(newCert.key),
		},
		Type: "kubernetes.io/tls",
	}
}
