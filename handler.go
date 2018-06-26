package main

import (
	log "github.com/Sirupsen/logrus"
	v1_podcounter "github.com/bitshifta/crd-controller/pkg/apis/podcounter/v1"
	podcounterclientset "github.com/bitshifta/crd-controller/pkg/client/clientset/versioned"
	"github.com/bitshifta/crd-controller/pkg/client/clientset/versioned/typed/podcounter/v1"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Handler interface {
	Init() error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
}

type PodCrdHandler struct {
	podCounterInterface v1.PodCounterInterface
	podcounterClientSet *podcounterclientset.Clientset
}

func (p *PodCrdHandler) Init() error {
	log.Info("PodCrdHandler.Init")
	p.podCounterInterface = p.podcounterClientSet.KhaliltV1().PodCounters(core_v1.NamespaceDefault)
	return nil
}

func (p *PodCrdHandler) ObjectCreated(obj interface{}) {
	log.Info("PodCrdHandler.ObjectCreated")

	counter, err := p.podCounterInterface.Get("cluster-counter", meta_v1.GetOptions{})

	if err != nil {
		log.Fatal("Could not find CRD, returning...")
		return
	}

	newHistoricCount := *counter.Spec.Historical + 1
	newCurrentCount := *counter.Spec.Current + 1

	err = p.updateCrdValues(&newHistoricCount, &newCurrentCount, counter)

	if err != nil {
		return
	}
}

func (p *PodCrdHandler) ObjectDeleted(obj interface{}) {
	log.Info("PodCrdHandler.ObjectDeleted")

	counter, err := p.podCounterInterface.Get("cluster-counter", meta_v1.GetOptions{})

	if err != nil {
		log.Fatal("Could not find CRD, returning...")
		return
	}

	newHistoricCount := *counter.Spec.Historical
	newCurrentCount := *counter.Spec.Current - 1

	err = p.updateCrdValues(&newHistoricCount, &newCurrentCount, counter)

	if err != nil {
		return
	}

	return
}

func (p *PodCrdHandler) updateCrdValues(historic *int32, current *int32, counter *v1_podcounter.PodCounter) error {
	counter.Spec.Historical = historic
	counter.Spec.Current = current

	_, err := p.podCounterInterface.Update(counter)

	if err != nil {
		return err
	}

	log.Infof("Counter.historical: %v", *counter.Spec.Historical)
	log.Infof("Counter.current: %v", *counter.Spec.Current)
	log.Info("Cluster count successfully updated...\n")

	return nil
}
