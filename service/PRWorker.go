package main

import (
	"context"
	"fmt"
	"github.com/domsu/pullassistant/analysis"
	"github.com/domsu/pullassistant/api"
	"github.com/domsu/pullassistant/template"
	"github.com/google/go-github/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"sourcegraph.com/sourcegraph/go-diff/diff"
	"strconv"
	"time"
)

type PRWorker struct {
	jobChan <-chan PRWorkerData
	logger  zerolog.Logger
	tracker *Tracker
	config  *api.Config
}

type PRWorkerData struct {
	client githubapp.ClientCreator
	event  github.PullRequestEvent
}

func (w *PRWorker) start() {
	go func() {
		for data := range w.jobChan {
			err := w.process(data)
			if err != nil {
				w.logger.Error().Err(err).Msg("failed to process PR")
			}
		}
	}()
}

func (w PRWorker) process(data PRWorkerData) error {
	ctx, logger, cancel := prepareContext(w, data)

	logger.Info().Msg("Processing PR started")
	logger.Info().Msg(fmt.Sprintf("Repository: %s", data.event.GetRepo().GetHTMLURL()+"/pull/"+strconv.Itoa(*data.event.Number)))

	w.tracker.trackRequest(data.event.GetRepo().GetOwner().GetLogin(), *data.event.PullRequest.Commits)

	client, err := data.client.NewInstallationClient(githubapp.GetInstallationIDFromEvent(&data.event))
	if err != nil {
		return errors.Wrap(err, "failed to get installation client")
	}

	commits, _, err := client.PullRequests.ListCommits(ctx, data.event.GetRepo().GetOwner().GetLogin(), data.event.GetRepo().GetName(), *data.event.Number, &github.ListOptions{PerPage: 250})
	if err != nil {
		return errors.Wrap(err, "failed to get commits list")
	}

	logger.Info().Msg(fmt.Sprintf("Number of commits: %d", len(commits)))

	if len(commits) == 0 {
		logger.Info().Msg("No commits - skipping")
		return nil
	}

	var commitsDiff [][]*diff.FileDiff
	for _, commitDetail := range commits {
		rawCommit, _, err := client.Repositories.GetCommitRaw(ctx, data.event.GetRepo().GetOwner().GetLogin(), data.event.GetRepo().GetName(), *commitDetail.SHA, github.RawOptions{Type: github.Diff})
		if err != nil {
			return errors.Wrap(err, "failed to get raw commit")
		}

		fileDiff, err := diff.ParseMultiFileDiff([]byte(rawCommit))
		if err != nil {
			return errors.Wrap(err, "failed to parse multiple file diff")
		}

		commitsDiff = append(commitsDiff, fileDiff)
	}

	combinedAnalysis := analysis.GetCombinedAnalysis(commitsDiff)
	comment := GetPRComment(data.event.PullRequest, commits, combinedAnalysis)

	scoreText := fmt.Sprint("Score: ", comment.BestReviewedCommitByCommit, combinedAnalysis.CommitByCommitReviewScore)
	logger.Info().Msg(scoreText)

	content := template.GetPRCommentContent(comment)

	logger.Info().Msg(fmt.Sprintf("Comment content: %s", content))

	prComment := github.IssueComment{
		Body: &content,
	}

	existingComments, _, err := client.Issues.ListComments(ctx, data.event.GetRepo().GetOwner().GetLogin(), data.event.GetRepo().GetName(), *data.event.Number, &github.IssueListCommentsOptions{
		Sort:      "created",
		Direction: "asc",
	})
	if err != nil {
		return errors.Wrap(err, "failed to get existing comments")
	}

	var commentIdToEdit *int64 = nil
	for _, comment := range existingComments {
		if *comment.User.Login == w.config.AppConfig.CommentAppName {
			commentIdToEdit = comment.ID
			break
		}
	}

	//self healing
	if commentIdToEdit != nil {
		for _, comment := range existingComments {
			if *comment.ID != *commentIdToEdit && *comment.User.Login == w.config.AppConfig.CommentAppName {
				if _, err := client.Issues.DeleteComment(ctx, data.event.GetRepo().GetOwner().GetLogin(), data.event.GetRepo().GetName(), *comment.ID); err != nil {
					logger.Error().Err(err).Msg("failed to delete duplicated comment")
				}
			}
		}
	}

	if commentIdToEdit != nil {
		if _, _, err := client.Issues.EditComment(ctx, data.event.GetRepo().GetOwner().GetLogin(), data.event.GetRepo().GetName(), *commentIdToEdit, &prComment); err != nil {
			return errors.Wrap(err, "failed to edit comment on pull request")
		}
	} else {
		if _, _, err := client.Issues.CreateComment(ctx, data.event.GetRepo().GetOwner().GetLogin(), data.event.GetRepo().GetName(), *data.event.Number, &prComment); err != nil {
			return errors.Wrap(err, "failed to add comment on pull request")
		}
	}

	//todo defer?
	cancel()

	logger.Info().Msg("Processing PR finished successfully")

	return nil
}

func prepareContext(w PRWorker, data PRWorkerData) (context.Context, zerolog.Logger, context.CancelFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	ctx = w.logger.WithContext(ctx)
	ctx, logger := githubapp.PreparePRContext(ctx, githubapp.GetInstallationIDFromEvent(&data.event), data.event.GetRepo(), *data.event.Number)

	return ctx, logger, cancel
}
