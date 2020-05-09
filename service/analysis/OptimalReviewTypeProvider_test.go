package analysis

import "testing"

func TestGetOptimalReviewType(t *testing.T) {
	if GetOptimalReviewType(0.2) != AllChanges {
		t.Error("Invalid review type. Expecting AllChanges")
	}

	if GetOptimalReviewType(0.9) != CommitByCommit {
		t.Error("Invalid review type. Expecting CommitByCommit")
	}
}
