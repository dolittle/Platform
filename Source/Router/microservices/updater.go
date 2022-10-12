package microservices

import (
	"github.com/rs/zerolog/log"
	coreV1 "k8s.io/api/core/v1"
)

type updater struct {
	Registry *Registry
	Config   MicroserviceConfiguration
}

func NewUpdater(registry *Registry, config MicroserviceConfiguration) *updater {
	return &updater{registry, config}
}

func (u *updater) Add(pod *coreV1.Pod) {
	logger := log.With().Str("component", "Updater").Str("method", "Add").Logger()

	microservice, err := convertPodToMicroservice(pod, u.Config)
	if err != nil {
		logger.Error().Err(err).Interface("microservice", microservice.Identity).Msg("")
		return
	}

	logger.Trace().Interface("microservice", microservice.Identity).Msg("Updating in registry")
	u.Registry.Upsert(microservice)
}

func (u *updater) Update(pod *coreV1.Pod) {
	logger := log.With().Str("component", "Updater").Str("method", "Update").Logger()

	microservice, err := convertPodToMicroservice(pod, u.Config)
	if err != nil {
		logger.Error().Err(err).Interface("microservice", microservice.Identity).Msg("")
		return
	}

	logger.Trace().Interface("microservice", microservice.Identity).Msg("Updating in registry")
	u.Registry.Upsert(microservice)
}

func (u *updater) Delete(pod *coreV1.Pod) {
	logger := log.With().Str("component", "Updater").Str("method", "Delete").Logger()

	microservice, err := convertPodToMicroservice(pod, u.Config)
	if err != nil {
		logger.Error().Err(err).Interface("microservice", microservice.Identity).Msg("")
		return
	}

	logger.Trace().Interface("microservice", microservice.Identity).Msg("Deleting from registry")
	u.Registry.Delete(microservice)
}
