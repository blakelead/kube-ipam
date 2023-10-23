package informer

import (
	"context"
	"log"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/blakelead/kube-ipam/internal/ipam"
)

func StartInformer(clientset *kubernetes.Clientset, ipPool *ipam.IPPool) {
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	informer := informerFactory.Core().V1().Services().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { addHandler(obj, clientset, ipPool) },
		DeleteFunc: func(obj interface{}) { deleteHandler(obj, ipPool) },
	})

	stopCh := make(chan struct{})
	defer close(stopCh)
	log.Println("kube-ipam started")
	informerFactory.Start(stopCh)
	<-stopCh
}

func addHandler(obj interface{}, clientset *kubernetes.Clientset, ipPool *ipam.IPPool) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		log.Fatalf("Expected *corev1.Service but got %T", obj)
		return
	}

	if svc.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return
	}

	ip, err := ipPool.AllocateIP()
	if err != nil {
		log.Println("Failed to allocate IP: ", err)
		return
	}

	// allocate IP
	svc.Spec.LoadBalancerIP = ip.String()
	_, err = clientset.CoreV1().Services(svc.Namespace).Update(context.TODO(), svc, metav1.UpdateOptions{})
	if err != nil {
		log.Println("Failed to allocate IP:", err)
		return
	}

	// update status and retry if error
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		newSvc, err := clientset.CoreV1().Services(svc.Namespace).Get(context.TODO(), svc.Name, metav1.GetOptions{})
		if err != nil {
			log.Println("error getting service, retrying:", err)
		}
		newSvc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{IP: ip.String()}}
		_, err = clientset.CoreV1().Services(svc.Namespace).UpdateStatus(context.TODO(), newSvc, metav1.UpdateOptions{})
		if err != nil {
			log.Println("failed to allocate IP, retrying:", err)
		}
		return err
	})
	if retryErr != nil {
		log.Println(retryErr)
		return
	}
	log.Printf("Allocated IP %s to service %s/%s", ip.String(), svc.Namespace, svc.Name)
}

func deleteHandler(obj interface{}, ipPool *ipam.IPPool) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		log.Println("couldn't get service")
	}
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		ipPool.DeallocateIP(svc.Spec.LoadBalancerIP)
		log.Printf("Deallocated IP %s for service %s/%s", svc.Spec.LoadBalancerIP, svc.Namespace, svc.Name)
	}
}
