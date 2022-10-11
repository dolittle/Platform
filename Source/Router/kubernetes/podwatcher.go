package kubernetes

import (
	"github.com/rs/zerolog/log"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"time"
)

type PodHandler interface {
	Add(pod *coreV1.Pod)
	Update(pod *coreV1.Pod)
	Delete(pod *coreV1.Pod)
}

func StartNewPodWatcher(client kubernetes.Interface, resyncPeriod time.Duration, labelSelector string, handler PodHandler, stop <-chan struct{}) {
	factory := informers.NewSharedInformerFactoryWithOptions(client, resyncPeriod, informers.WithTweakListOptions(func(options *metaV1.ListOptions) {
		options.LabelSelector = labelSelector
	}))

	informer := factory.Core().V1().Pods().Informer()
	informer.AddEventHandler(&eventHandler{
		handler: handler,
	})

	go informer.Run(stop)
}

type eventHandler struct {
	handler PodHandler
}

func (e *eventHandler) OnAdd(obj interface{}) {
	pod, ok := obj.(*coreV1.Pod)
	if !ok {
		log.Error().
			Str("component", "PodHandler").
			Str("method", "OnAdd").
			Msg("Object was not a Pod")
	}

	e.handler.Add(pod)
}

func (e *eventHandler) OnUpdate(oldObj, newObj interface{}) {
	oldPod, oldOk := oldObj.(*coreV1.Pod)
	newPod, newOk := newObj.(*coreV1.Pod)

	if !oldOk || !newOk {
		log.Error().
			Str("component", "PodHandler").
			Str("method", "OnUpdate").
			Msg("Object was not a Pod")
	}

	if oldPod.GetResourceVersion() == newPod.GetResourceVersion() {
		log.Trace().
			Str("component", "PodHandler").
			Str("method", "OnUpdate").
			Msg("Skipping update because resource version was unchanged")
		return
	}

	e.handler.Update(newPod)
}

func (e *eventHandler) OnDelete(obj interface{}) {
	pod, ok := obj.(*coreV1.Pod)
	if !ok {
		log.Error().
			Str("component", "PodHandler").
			Str("method", "OnDelete").
			Msg("Object was not a Pod")
	}

	e.handler.Delete(pod)
}
