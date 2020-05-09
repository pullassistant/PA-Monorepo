package analysis

type ReviewType int

const minScoreForCommitByCommitReviewType = 0.5

const (
	CommitByCommit ReviewType = iota
	AllChanges
)

func GetOptimalReviewType(commitByCommitReviewScore float32) ReviewType {
	if commitByCommitReviewScore >= minScoreForCommitByCommitReviewType {
		return CommitByCommit
	} else {
		return AllChanges
	}
}
