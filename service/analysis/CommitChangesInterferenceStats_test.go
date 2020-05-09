package analysis

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

var commitFileChangesInterferenceStatsTestData = []struct {
	commitChanges []CommitChanges
	expected      []CommitChangesInterferenceStats
}{
	{
		commitChanges: []CommitChanges{
			{FileChanges: []CommitFileChanges{
				{
					File: "service/testdata.txt",
					Origin: FileOrigin{
						OriginalFile:       "service/testdata.txt",
						OriginCommitNumber: 0,
					},
					Changes: []FileChange{
						{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
					},
				}}},
		},
		expected: []CommitChangesInterferenceStats{
			{
				FileStats: []CommitFileChangesInterfereStats{
					{
						Origin: FileOrigin{
							OriginalFile:       "service/testdata.txt",
							OriginCommitNumber: 0,
						},
						InterfereStats: nil,
					},
				}},
		},
	},
	{
		commitChanges: []CommitChanges{
			{
				FileChanges: []CommitFileChanges{
					{
						File: "service/testdata.txt",
						Origin: FileOrigin{
							OriginalFile:       "service/testdata.txt",
							OriginCommitNumber: 0,
						},
						Changes: []FileChange{
							{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
						}},
				},
			},
			{
				FileChanges: []CommitFileChanges{
					{
						File: "service/testdata.txt",
						Origin: FileOrigin{
							OriginalFile:       "service/testdata.txt",
							OriginCommitNumber: 0,
						},
						Changes: []FileChange{
							{StartLine: 1, AffectedLines: 7, CommitNumber: 0},
						}},
				},
			},
		},
		expected: []CommitChangesInterferenceStats{
			{FileStats: []CommitFileChangesInterfereStats{
				{
					Origin: FileOrigin{
						OriginalFile:       "service/testdata.txt",
						OriginCommitNumber: 0,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 1,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 14,
								FinalNumberOfChanges:   7,
								InterferenceValue:      0.5,
							},
						},
					},
				},
			}},
			{
				FileStats: nil,
			},
		},
	},
	{
		commitChanges: []CommitChanges{
			{
				FileChanges: []CommitFileChanges{
					{
						File: "service/testdata.txt",
						Origin: FileOrigin{
							OriginalFile:       "service/testdata.txt",
							OriginCommitNumber: 0,
						},
						Changes: []FileChange{
							{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
						}},
				},
			},
			{
				FileChanges: []CommitFileChanges{
					{
						File: "service/testdata.txt",
						Origin: FileOrigin{
							OriginalFile:       "service/testdata.txt",
							OriginCommitNumber: 0,
						},
						Changes: []FileChange{
							{StartLine: 1, AffectedLines: 7, CommitNumber: 0},
							{StartLine: 10, AffectedLines: 1, CommitNumber: 1},
						}},
				},
			},
			{
				FileChanges: []CommitFileChanges{
					{
						File: "service/testdata.txt",
						Origin: FileOrigin{
							OriginalFile:       "service/testdata.txt",
							OriginCommitNumber: 0,
						},
						Changes: []FileChange{
							{StartLine: 1, AffectedLines: 7, CommitNumber: 0},
							{StartLine: 10, AffectedLines: 1, CommitNumber: 1},
							{StartLine: 20, AffectedLines: 1, CommitNumber: 2},
						}},
				},
			},
		},
		expected: []CommitChangesInterferenceStats{

			{FileStats: []CommitFileChangesInterfereStats{
				{
					Origin: FileOrigin{
						OriginalFile:       "service/testdata.txt",
						OriginCommitNumber: 0,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 1,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 14,
								FinalNumberOfChanges:   7,
								InterferenceValue:      0.5,
							},
						},
						{
							CommitNumber: 2,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 14,
								FinalNumberOfChanges:   7,
								InterferenceValue:      0.5,
							},
						},
					},
				}}},
			{[]CommitFileChangesInterfereStats{{
				Origin: FileOrigin{
					OriginalFile:       "service/testdata.txt",
					OriginCommitNumber: 0,
				},
				InterfereStats: []FileInterfereStats{
					{
						CommitNumber: 2,
						InterferenceInfo: InterferenceInfo{
							InitialNumberOfChanges: 1,
							FinalNumberOfChanges:   1,
							InterferenceValue:      0,
						},
					},
				},
			}},
			},
			{[]CommitFileChangesInterfereStats{{
				Origin: FileOrigin{
					OriginalFile:       "service/testdata.txt",
					OriginCommitNumber: 0,
				},
				InterfereStats: nil,
			},
			}},
		}},
	{
		commitChanges: []CommitChanges{
			{
				FileChanges: []CommitFileChanges{
					{
						File: "service/testdata.txt",
						Origin: FileOrigin{
							OriginalFile:       "service/testdata.txt",
							OriginCommitNumber: 0,
						},
						Changes: []FileChange{
							{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
						}},
				},
			},
			{
				FileChanges: []CommitFileChanges{
					{
						File: "service/testdata2.txt",
						Origin: FileOrigin{
							OriginalFile:       "service/testdata.txt",
							OriginCommitNumber: 0,
						},
						Changes: []FileChange{
							{StartLine: 1, AffectedLines: 7, CommitNumber: 0},
						}},
				},
			},
		},
		expected: []CommitChangesInterferenceStats{
			{FileStats: []CommitFileChangesInterfereStats{
				{
					Origin: FileOrigin{
						OriginalFile:       "service/testdata.txt",
						OriginCommitNumber: 0,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 1,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 14,
								FinalNumberOfChanges:   7,
								InterferenceValue:      0.5,
							},
						},
					},
				},
			}},
			{
				FileStats: nil,
			},
		},
	},
}

func TestGetInterferenceWithFollowingCommits(t *testing.T) {
	testData := CommitChangesInterferenceStats{FileStats: []CommitFileChangesInterfereStats{
		{
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			InterfereStats: []FileInterfereStats{
				{
					CommitNumber: 1,
					InterferenceInfo: InterferenceInfo{
						InitialNumberOfChanges: 14,
						FinalNumberOfChanges:   7,
						InterferenceValue:      0.5,
					},
				},
				{
					CommitNumber: 2,
					InterferenceInfo: InterferenceInfo{
						InitialNumberOfChanges: 14,
						FinalNumberOfChanges:   0,
						InterferenceValue:      0,
					},
				},
			},
		},
		{
			Origin: FileOrigin{
				OriginalFile:       "service/testdata2.txt",
				OriginCommitNumber: 0,
			},
			InterfereStats: []FileInterfereStats{
				{
					CommitNumber: 1,
					InterferenceInfo: InterferenceInfo{
						InitialNumberOfChanges: 10,
						FinalNumberOfChanges:   7,
						InterferenceValue:      0.3,
					},
				},
				{
					CommitNumber: 2,
					InterferenceInfo: InterferenceInfo{
						InitialNumberOfChanges: 10,
						FinalNumberOfChanges:   2,
						InterferenceValue:      0.8,
					},
				},
			},
		},
		{
			Origin: FileOrigin{
				OriginalFile:       "service/testdata3.txt",
				OriginCommitNumber: 0,
			},
			InterfereStats: []FileInterfereStats{
				{
					CommitNumber: 1,
					InterferenceInfo: InterferenceInfo{
						InitialNumberOfChanges: 14,
						FinalNumberOfChanges:   0,
						InterferenceValue:      0,
					},
				},
				{
					CommitNumber: 2,
					InterferenceInfo: InterferenceInfo{
						InitialNumberOfChanges: 14,
						FinalNumberOfChanges:   7,
						InterferenceValue:      0.5,
					},
				},
			},
		}}}

	expected := InterferenceInfo{
		InitialNumberOfChanges: 38,
		FinalNumberOfChanges:   16,
		InterferenceValue:      0.57894736842,
	}

	received := testData.GetInterferenceWithFollowingCommits(1, 3)

	if !reflect.DeepEqual(received, expected) {
		t.Error(fmt.Sprint("Test failed.", expected, received))
	}
}

func TestGetInterferenceWithFollowingCommitsWhenNoInterference(t *testing.T) {
	testData := CommitChangesInterferenceStats{FileStats: []CommitFileChangesInterfereStats{
		{
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			InterfereStats: []FileInterfereStats{
				{
					CommitNumber: 1,
					InterferenceInfo: InterferenceInfo{
						InitialNumberOfChanges: 14,
						FinalNumberOfChanges:   0,
						InterferenceValue:      0,
					},
				},
				{
					CommitNumber: 2,
					InterferenceInfo: InterferenceInfo{
						InitialNumberOfChanges: 14,
						FinalNumberOfChanges:   0,
						InterferenceValue:      0,
					},
				},
			},
		}}}

	expected := InterferenceInfo{
		InitialNumberOfChanges: 14,
		FinalNumberOfChanges:   0,
		InterferenceValue:      0,
	}

	received := testData.GetInterferenceWithFollowingCommits(1, 3)

	if !reflect.DeepEqual(received, expected) {
		t.Error(fmt.Sprint("Test failed.", expected, received))
	}
}

func TestHasAnyInterference(t *testing.T) {
	if (InterferenceInfo{
		InitialNumberOfChanges: 1,
		FinalNumberOfChanges:   0,
		InterferenceValue:      0,
	}.hasAnyInterference()) {
		t.Error("Test failed")
	}

	if !(InterferenceInfo{
		InitialNumberOfChanges: 1,
		FinalNumberOfChanges:   0,
		InterferenceValue:      1,
	}.hasAnyInterference()) {
		t.Error("Test failed")
	}
}

func TestGetCommitFileChangesInterferenceStats(t *testing.T) {
	for i, testData := range commitFileChangesInterferenceStatsTestData {
		stats := GetCommitFileChangesInterferenceStats(testData.commitChanges)

		if !reflect.DeepEqual(stats, testData.expected) {
			bytesExpected, _ := json.Marshal(testData.expected)
			t.Log("Expected: " + string(bytesExpected))
			bytesOperations, _ := json.Marshal(stats)
			t.Log("Received: " + string(bytesOperations))
			t.Error("Scenario " + strconv.Itoa(i) + " failed.")
		}
	}
}
