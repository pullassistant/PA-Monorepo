package analysis

import (
	"reflect"
	"strconv"
	"testing"
)

var commitByCommitReviewScoreTestData = []struct {
	interferenceStats []CommitChangesInterferenceStats
	expected          float32
}{
	{
		interferenceStats: []CommitChangesInterferenceStats{
			{
				FileStats: []CommitFileChangesInterfereStats{
					{
						Origin: FileOrigin{
							OriginalFile:       "service/testdata.txt",
							OriginCommitNumber: 0,
						},
						InterfereStats: []FileInterfereStats{
							{
								CommitNumber: 1,
								InterferenceInfo: InterferenceInfo{
									InitialNumberOfChanges: 10,
									FinalNumberOfChanges:   9,
									InterferenceValue:      0.1,
								},
							},
							{
								CommitNumber: 2,
								InterferenceInfo: InterferenceInfo{
									InitialNumberOfChanges: 10,
									FinalNumberOfChanges:   8,
									InterferenceValue:      0.2,
								},
							},
						},
					},
				},
			},
			{
				FileStats: []CommitFileChangesInterfereStats{
					{
						Origin: FileOrigin{
							OriginalFile:       "service/testdata2.txt",
							OriginCommitNumber: 1,
						},
						InterfereStats: []FileInterfereStats{
							{
								CommitNumber: 2,
								InterferenceInfo: InterferenceInfo{
									InitialNumberOfChanges: 20,
									FinalNumberOfChanges:   14,
									InterferenceValue:      0.3,
								},
							},
						},
					},
				},
			},
			{
				FileStats: []CommitFileChangesInterfereStats{
					{
						Origin: FileOrigin{
							OriginalFile:       "service/testdata2.txt",
							OriginCommitNumber: 1,
						},
						InterfereStats: nil,
					},
				},
			},
		},
		expected: 0.7333333333,
	},
	{
		interferenceStats: []CommitChangesInterferenceStats{
			{
				FileStats: []CommitFileChangesInterfereStats{
					{
						Origin: FileOrigin{
							OriginalFile:       "service/testdata2.txt",
							OriginCommitNumber: 0,
						},
						InterfereStats: nil,
					},
				},
			},
		},
		expected: 1,
	},
}

func TestGetCommitByCommitReviewScore(t *testing.T) {
	for i, testData := range commitByCommitReviewScoreTestData {
		score := GetCommitByCommitReviewScore(testData.interferenceStats)

		if !reflect.DeepEqual(score, testData.expected) {
			t.Log("Expected: " + strconv.FormatFloat(float64(testData.expected), 'f', 2, 64))
			t.Log("Received: ", strconv.FormatFloat(float64(score), 'f', 2, 64))
			t.Error("Scenario " + strconv.Itoa(i) + " failed.")
		}
	}
}
