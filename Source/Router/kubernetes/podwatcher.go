package kubernetes

import (
	"context"
	"github.com/dolittle/platform-router/config"
	"github.com/rs/zerolog/log"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

// PodHandler defines a system that reacts to when a pod is discovered, updated or deleted by adding, updating or deleting.
type PodHandler interface {
	Add(pod *coreV1.Pod)
	Update(pod *coreV1.Pod)
	Delete(pod *coreV1.Pod)
}

type PodWatcher struct {
	Client                  kubernetes.Interface
	Config                  *config.Config
	LabelSelectorConfigPath string
	Handler                 PodHandler
}

func (pw *PodWatcher) Run(resyncPeriod time.Duration, ctx context.Context) {
loop:
	for {
		change := pw.Config.Changed()

		factory := informers.NewSharedInformerFactoryWithOptions(pw.Client, resyncPeriod, informers.WithTweakListOptions(func(options *metaV1.ListOptions) {
			options.LabelSelector = pw.Config.String(pw.LabelSelectorConfigPath)
		}))

		informer := factory.Core().V1().Pods().Informer()
		informer.AddEventHandler(pw)

		go informer.Run(change)

		select {
		case <-ctx.Done():
			break loop
		case <-change:
		}
	}
}

func (e *PodWatcher) OnAdd(obj interface{}) {
	pod, ok := obj.(*coreV1.Pod)
	if !ok {
		log.Error().
			Str("component", "PodHandler").
			Str("method", "OnAdd").
			Msg("Object was not a Pod")
	}

	e.Handler.Add(pod)
}

func (e *PodWatcher) OnUpdate(oldObj, newObj interface{}) {
	pod, ok := newObj.(*coreV1.Pod)

	if !ok {
		log.Error().
			Str("component", "PodHandler").
			Str("method", "OnUpdate").
			Msg("Object was not a Pod")
	}

	e.Handler.Update(pod)
}

func (e *PodWatcher) OnDelete(obj interface{}) {
	pod, ok := obj.(*coreV1.Pod)
	if !ok {
		log.Error().
			Str("component", "PodHandler").
			Str("method", "OnDelete").
			Msg("Object was not a Pod")
	}

	e.Handler.Delete(pod)
}
