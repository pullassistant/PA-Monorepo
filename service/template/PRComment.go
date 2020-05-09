package template

type PRCommentData struct {
	BestReviewedCommitByCommit bool

	CommitByCommitReviewScore       float32
	CommitByCommitReviewScoreImgUrl string

	OptimalCodeReviewWarningsCount        int
	OptimalCodeReviewSquashedCommitsCount int

	OptimalCodeReview []PRCommentOptimalReviewCommit

	LastUpdateFromText string
	LastUpdateFromUrl  string

	LastUpdateToText string
	LastUpdateToUrl  string
}

type PRCommentOptimalReviewCommit struct {
	CommitName       string
	CommitUrl        string
	CommitSizeImgUrl string
	Warnings         []PRCommentOptimalReviewCommitWarning
	SquashedCommits  []PRCommentOptimalReviewCommitSquashed
}

type PRCommentOptimalReviewCommitWarning struct {
	FileName              string
	Interference          float32
	InterferingCommitName string
	InterferingCommitUrl  string
}
type PRCommentOptimalReviewCommitSquashed struct {
	CommitName string
	CommitUrl  string
}
