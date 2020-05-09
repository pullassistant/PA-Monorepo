package main

import (
	"github.com/domsu/pullassistant/analysis"
	. "github.com/domsu/pullassistant/template"
	"github.com/google/go-github/github"
	"strings"
)

const maxCommitNameLength = 70
const maxCommitNameWarningLength = 20
const maxCommitNameSquashedLength = 20

const maxFileNameLength = 30

func GetPRComment(pullRequest *github.PullRequest, commits []*github.RepositoryCommit, combinedAnalysis analysis.CombinedAnalysis) PRCommentData {
	var optimalCodeReview []PRCommentOptimalReviewCommit
	warningsCount := 0

	if combinedAnalysis.OptimalReviewType == analysis.CommitByCommit {
		optimalCodeReview = getOptimalCommitByCommitReview(combinedAnalysis, commits, pullRequest, optimalCodeReview)
		for _, v := range combinedAnalysis.FileInterferenceInfo {
			if v.FileInfos != nil {
				warningsCount++
			}
		}
	} else {
		optimalCodeReview = getOptimalWithAllChangesReview(commits, pullRequest, optimalCodeReview)
	}

	lastUpdateFromText := *commits[0].SHA
	lastUpdateToText := *commits[(len(commits) - 1)].SHA

	squashedCommitsCount := 0
	for _, v := range optimalCodeReview {
		squashedCommitsCount += len(v.SquashedCommits)
	}

	return PRCommentData{
		BestReviewedCommitByCommit:            combinedAnalysis.OptimalReviewType == analysis.CommitByCommit,
		CommitByCommitReviewScore:             combinedAnalysis.CommitByCommitReviewScore,
		CommitByCommitReviewScoreImgUrl:       getCommitByCommitReviewScoreUrl(combinedAnalysis.CommitByCommitReviewScore),
		OptimalCodeReviewWarningsCount:        warningsCount,
		OptimalCodeReviewSquashedCommitsCount: squashedCommitsCount,
		OptimalCodeReview:                     optimalCodeReview,
		LastUpdateFromText:                    lastUpdateFromText[:7],
		LastUpdateFromUrl:                     *commits[0].HTMLURL,
		LastUpdateToText:                      lastUpdateToText[:7],
		LastUpdateToUrl:                       *commits[(len(commits) - 1)].HTMLURL,
	}
}

func getOptimalWithAllChangesReview(commits []*github.RepositoryCommit, pullRequest *github.PullRequest, optimalCodeReview []PRCommentOptimalReviewCommit) []PRCommentOptimalReviewCommit {
	var squashedCommits []PRCommentOptimalReviewCommitSquashed
	for _, v := range commits {
		startCommitName := *v.Commit.Message
		squashedCommits = append(squashedCommits, PRCommentOptimalReviewCommitSquashed{
			CommitName: getFirstLine(truncateEnd(startCommitName, maxCommitNameLength)),
			CommitUrl:  getReviewCommitUrl(*pullRequest.HTMLURL, *v.SHA),
		})
	}
	optimalCodeReview = append(optimalCodeReview, PRCommentOptimalReviewCommit{
		CommitName:       getSquashedCommitName(commits, 0, len(commits)),
		CommitUrl:        getSquashedReviewCommitUrl(pullRequest, commits, 0, len(commits)),
		CommitSizeImgUrl: "https://pullassistant.com/assets/size3.svg",
		Warnings:         nil,
		SquashedCommits:  squashedCommits,
	})
	return optimalCodeReview
}

func getOptimalCommitByCommitReview(combinedAnalysis analysis.CombinedAnalysis, commits []*github.RepositoryCommit, pullRequest *github.PullRequest, optimalCodeReview []PRCommentOptimalReviewCommit) []PRCommentOptimalReviewCommit {
	for i, v := range combinedAnalysis.OptimalCommitByCommitReview {
		var warnings []PRCommentOptimalReviewCommitWarning
		for _, fileInfo := range combinedAnalysis.FileInterferenceInfo[i].FileInfos {
			commitName := *commits[fileInfo.CommitNumber].Commit.Message

			warnings = append(warnings, PRCommentOptimalReviewCommitWarning{
				FileName:              truncateStart(fileInfo.File, maxFileNameLength),
				Interference:          fileInfo.InterferenceInfo.InterferenceValue,
				InterferingCommitName: getFirstLine(truncateEnd(commitName, maxCommitNameWarningLength)),
				InterferingCommitUrl:  getReviewCommitUrl(*pullRequest.HTMLURL, *commits[fileInfo.CommitNumber].SHA),
			})
		}

		var truncatedCommitName = ""
		var commitUrl = ""
		var squashedCommits []PRCommentOptimalReviewCommitSquashed
		if v.EndCommitNumber-v.StartCommitNumber > 1 {
			for squashedCommitNumber := v.StartCommitNumber; squashedCommitNumber < v.EndCommitNumber; squashedCommitNumber++ {
				startCommitName := *commits[squashedCommitNumber].Commit.Message
				squashedCommits = append(squashedCommits, PRCommentOptimalReviewCommitSquashed{
					CommitName: getFirstLine(truncateEnd(startCommitName, maxCommitNameLength)),
					CommitUrl:  getReviewCommitUrl(*pullRequest.HTMLURL, *commits[squashedCommitNumber].SHA),
				})
			}

			truncatedCommitName = getSquashedCommitName(commits, v.StartCommitNumber, v.EndCommitNumber)
			commitUrl = getSquashedReviewCommitUrl(pullRequest, commits, v.StartCommitNumber, v.EndCommitNumber)
		} else {
			truncatedCommitName = getFirstLine(truncateEnd(*commits[v.StartCommitNumber].Commit.Message, maxCommitNameLength))
			commitUrl = getReviewCommitUrl(*pullRequest.HTMLURL, *commits[v.StartCommitNumber].SHA)
		}

		optimalCodeReview = append(optimalCodeReview, PRCommentOptimalReviewCommit{
			CommitName:       truncatedCommitName,
			CommitUrl:        commitUrl,
			CommitSizeImgUrl: "https://pullassistant.com/assets/size3.svg",
			Warnings:         warnings,
			SquashedCommits:  squashedCommits,
		})
	}
	return optimalCodeReview
}

func getSquashedCommitName(commits []*github.RepositoryCommit, startCommitNumber int, endCommitNumber int) string {
	return getFirstLine(truncateEnd(*commits[startCommitNumber].Commit.Message, maxCommitNameSquashedLength)) +
		" ... " + getFirstLine(truncateEnd(*commits[endCommitNumber-1].Commit.Message, maxCommitNameSquashedLength))
}

func truncateEnd(string string, max int) string {
	if len(string) > max {
		return string[:max] + "..."
	} else {
		return string
	}
}

func truncateStart(string string, max int) string {
	if len(string) > max {
		return "..." + string[len(string)-max:]
	} else {
		return string
	}
}

func getFirstLine(string string) string {
	index := strings.Index(string, "\n")
	if index == -1 {
		return string
	} else {
		return string[:index]
	}
}

func getReviewCommitUrl(pullRequestHtmlUrl string, commitSha string) string {
	return pullRequestHtmlUrl + "/commits/" + commitSha
}

func getSquashedReviewCommitUrl(pullRequest *github.PullRequest, commits []*github.RepositoryCommit, startCommitNumber int, endCommitNumber int) string {
	if startCommitNumber == 0 {
		return *pullRequest.HTMLURL + "/files/" + *commits[endCommitNumber-1].SHA
	} else {
		parentFromCommitSha := *commits[startCommitNumber].Parents[0].SHA
		return *pullRequest.HTMLURL + "/files/" + parentFromCommitSha + ".." + *commits[endCommitNumber-1].SHA
	}
}

func getCommitByCommitReviewScoreUrl(commitByCommitReviewScore float32) string {
	if commitByCommitReviewScore >= 0.95 {
		return "https://pullassistant.com/assets/score10.svg"
	} else if commitByCommitReviewScore >= 0.85 {
		return "https://pullassistant.com/assets/score9.svg"
	} else if commitByCommitReviewScore >= 0.75 {
		return "https://pullassistant.com/assets/score8.svg"
	} else if commitByCommitReviewScore >= 0.65 {
		return "https://pullassistant.com/assets/score7.svg"
	} else if commitByCommitReviewScore >= 0.55 {
		return "https://pullassistant.com/assets/score6.svg"
	} else if commitByCommitReviewScore >= 0.45 {
		return "https://pullassistant.com/assets/score5.svg"
	} else if commitByCommitReviewScore >= 0.35 {
		return "https://pullassistant.com/assets/score4.svg"
	} else if commitByCommitReviewScore >= 0.25 {
		return "https://pullassistant.com/assets/score3.svg"
	} else if commitByCommitReviewScore >= 0.15 {
		return "https://pullassistant.com/assets/score2.svg"
	} else if commitByCommitReviewScore >= 0.05 {
		return "https://pullassistant.com/assets/score1.svg"
	} else {
		return "https://pullassistant.com/assets/score0.svg"
	}
}
