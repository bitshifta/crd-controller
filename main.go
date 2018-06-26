package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"flag"
	log "github.com/Sirupsen/logrus"
	podcounterapi_v1 "github.com/bitshifta/crd-controller/pkg/apis/podcounter/v1"
	podcounterclientset "github.com/bitshifta/crd-controller/pkg/client/clientset/versioned"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/client-go/rest"
)

const (
	counterName = "cluster-counter"
)

func main() {

	client, podcounterClientSet := getKubernetesClient()

	findOrCreatePodCounter(podcounterClientSet)

	informer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Pods(meta_v1.NamespaceAll).List(options)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Pods(meta_v1.NamespaceAll).Watch(options)
		},
	}, &api_v1.Pod{}, 0, cache.Indexers{})

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			log.Infof("Add pod: %s", key)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			log.Infof("Delete pod: %s", key)
			if err == nil {
				queue.Add(key)
			}
		},
	})

	controller := Controller{
		logger:           log.NewEntry(log.New()),
		clientset:        client,
		podcounterClient: podcounterClientSet,
		informer:         informer,
		queue:            queue,
		handler: &PodCrdHandler{
			podCounterInterface: podcounterClientSet.KhaliltV1().PodCounters(api_v1.NamespaceDefault),
		},
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	go controller.Run(stopCh)

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}

func getKubernetesClient() (kubernetes.Interface, podcounterclientset.Interface) {
	kubeconfig := ""
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "kubeconfig file")
	flag.Parse()

	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	var (
		config *rest.Config
		err    error
	)

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating client: %v", err)
		os.Exit(1)
	}

	client, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	log.Info("Successfully constructed k8s client")

	podcounterClientSet, err := podcounterclientset.NewForConfig(config)

	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	return client, podcounterClientSet
}

func findOrCreatePodCounter(podcounterClientSet podcounterclientset.Interface) {
	podCounterInterface := podcounterClientSet.KhaliltV1().PodCounters(api_v1.NamespaceDefault)

	_, err := podCounterInterface.Get(counterName, meta_v1.GetOptions{})

	if err == nil {
		fmt.Printf("Found CRD in cluster...\n")
		return
	}

	podcounter := &podcounterapi_v1.PodCounter{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: counterName,
		},
		Spec: podcounterapi_v1.PodCounterSpec{
			Current:    int2ptr(0),
			Historical: int2ptr(0),
		},
	}

	result, err := podCounterInterface.Create(podcounter)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Created pod-counter %q.\n", result.GetObjectMeta().GetName())
}

func int2ptr(integer int32) *int32 {
	return &integer
}
