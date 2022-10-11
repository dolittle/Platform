package microservices

import (
	"github.com/rs/zerolog/log"
	coreV1 "k8s.io/api/core/v1"
)

type Updater struct {
	Registry *Registry
}

func (u *Updater) Add(pod *coreV1.Pod) {
	logger := log.With().Str("component", "Updater").Str("method", "Add").Logger()

	microservice, err := convertPodToMicroservice(pod)
	if err != nil {
		logger.Error().Err(err).Interface("microservice", microservice).Msg("")
		return
	}

	logger.Trace().Interface("microservice", microservice).Msg("Updating in registry")
	u.Registry.Upsert(microservice)
}

func (u *Updater) Update(pod *coreV1.Pod) {
	logger := log.With().Str("component", "Updater").Str("method", "Update").Logger()

	microservice, err := convertPodToMicroservice(pod)
	if err != nil {
		logger.Error().Err(err).Interface("microservice", microservice).Msg("")
		return
	}

	logger.Trace().Interface("microservice", microservice).Msg("Updating in registry")
	u.Registry.Upsert(microservice)
}

func (u *Updater) Delete(pod *coreV1.Pod) {
	logger := log.With().Str("component", "Updater").Str("method", "Delete").Logger()

	microservice, err := convertPodToMicroservice(pod)
	if err != nil {
		logger.Error().Err(err).Interface("microservice", microservice).Msg("")
		return
	}

	logger.Trace().Interface("microservice", microservice).Msg("Deleting from registry")
	u.Registry.Delete(microservice)
}
