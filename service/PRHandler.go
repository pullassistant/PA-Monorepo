package main

import (
	"context"
	"encoding/json"
	"github.com/google/go-github/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type PRHandler struct {
	githubapp.ClientCreator
	channel chan<- PRWorkerData
}

func (h *PRHandler) Handles() []string {
	return []string{"pull_request"}
}

func (h *PRHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var pullRequestEvent github.PullRequestEvent
	if err := json.Unmarshal(payload, &pullRequestEvent); err != nil {
		return errors.Wrap(err, "failed to parse pull request event")
	}

	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("Received event: " + pullRequestEvent.GetAction())

	h.channel <- PRWorkerData{
		client: h.ClientCreator,
		event:  pullRequestEvent,
	}

	return nil
}
