package main

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	podcounterclientset "github.com/bitshifta/crd-controller/pkg/client/clientset/versioned"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	//api_v1 "k8s.io/api/core/v1"
	//meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	numOfRetries = 5
)

type Controller struct {
	logger           *logrus.Entry
	clientset        kubernetes.Interface
	podcounterClient podcounterclientset.Interface
	queue            workqueue.RateLimitingInterface
	informer         cache.SharedIndexInformer
	handler          Handler
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Info("Controller initialising...")

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Error syncing cache"))
		return
	}

	c.logger.Info("Controller cache sync complete")
	c.logger.Infof("The cluster now has [%v] running pods...", len(c.informer.GetStore().ListKeys()))

	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}
func (c *Controller) runWorker() {
	c.logger.Info("Controller.runworker starting....")
	for c.processNextItem() {
		c.logger.Debug("Controller.runWorker: processing next item....")
	}
	c.logger.Info("Controller.runWorker: completed!")
}

func (c *Controller) processNextItem() bool {
	c.logger.Debug("Controller.processNextItem: starting....")

	key, quit := c.queue.Get()

	if quit {
		return false
	}

	defer c.queue.Done(key)

	rawKey := key.(string)

	item, exists, err := c.informer.GetIndexer().GetByKey(rawKey)

	if err != nil {
		if c.queue.NumRequeues(key) < numOfRetries {
			c.logger.Errorf("Controller.processNextItem: Failed processing item with key %s with error %v, retrying", key, err)
			c.queue.AddRateLimited(key)
		} else {
			c.logger.Errorf("Controller.processNextItem: Failed processing item with key %s with error %v, no more retries", key, err)
			c.queue.Forget(key)
			utilruntime.HandleError(err)
		}
	}

	if !exists {
		c.deleted(key, rawKey, item)

	} else {
		c.created(key, rawKey, item)
	}

	return true
}

func (c *Controller) deleted(key interface{}, rawKey string, item interface{}) {
	defer c.queue.Forget(key)

	c.logger.Infof("Controller.processNextItem: object deleted: %s", rawKey)

	c.handler.ObjectDeleted(item)
}

func (c *Controller) created(key interface{}, rawKey string, item interface{}) {
	defer c.queue.Forget(key)

	c.logger.Infof("Controller.processNextItem: object created detected: %s", rawKey)

	c.handler.ObjectCreated(item)

}
