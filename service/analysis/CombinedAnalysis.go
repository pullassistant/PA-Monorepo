package analysis

import "sourcegraph.com/sourcegraph/go-diff/diff"

type CombinedAnalysis struct {
	CommitChanges []CommitChanges

	CommitByCommitReviewScore float32

	OptimalReviewType ReviewType

	OptimalCommitByCommitReview []CombinedCommitsInterferenceStats

	FileInterferenceInfo []CommitReviewInterferenceInfo
}

func GetCombinedAnalysis(commitsDiff [][]*diff.FileDiff) CombinedAnalysis {
	var commitOperations []CommitOperations
	for _, v := range commitsDiff {
		commitOperations = append(commitOperations, GetCommitOperations(v))
	}

	commitChanges := GetAccumulatedCommitChanges(commitOperations)
	stats := GetCommitFileChangesInterferenceStats(commitChanges)
	commitByCommitReviewScore := GetCommitByCommitReviewScore(stats)
	optimalCommitByCommitReview := GetOptimalReview(stats)

	return CombinedAnalysis{
		CommitChanges:               commitChanges,
		CommitByCommitReviewScore:   commitByCommitReviewScore,
		OptimalReviewType:           GetOptimalReviewType(commitByCommitReviewScore),
		OptimalCommitByCommitReview: optimalCommitByCommitReview,
		FileInterferenceInfo:        GetOptimalReviewFileInterferenceInfo(optimalCommitByCommitReview, commitChanges),
	}
}
