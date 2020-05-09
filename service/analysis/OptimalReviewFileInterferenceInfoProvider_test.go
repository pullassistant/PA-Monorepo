package analysis

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetOptimalReviewFileInterferenceInfo(t *testing.T) {
	combinedCommitStats := []CombinedCommitsInterferenceStats{
		{
			StartCommitNumber: 0,
			EndCommitNumber:   1,
			InterferenceStats: CommitChangesInterferenceStats{
				FileStats: []CommitFileChangesInterfereStats{
					{
						Origin: FileOrigin{
							OriginalFile:       "a.txt",
							OriginCommitNumber: 0,
						},
						InterfereStats: []FileInterfereStats{
							{
								CommitNumber: 1,
								InterferenceInfo: InterferenceInfo{
									InitialNumberOfChanges: 100,
									FinalNumberOfChanges:   2,
									InterferenceValue:      0.98,
								},
							},
							{
								CommitNumber: 2,
								InterferenceInfo: InterferenceInfo{
									InitialNumberOfChanges: 100,
									FinalNumberOfChanges:   1,
									InterferenceValue:      0.99,
								},
							},
						},
					},
					{
						Origin: FileOrigin{
							OriginalFile:       "b.txt",
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
									FinalNumberOfChanges:   0,
									InterferenceValue:      0,
								},
							},
						},
					},
				},
			},
		},
		{
			StartCommitNumber: 1,
			EndCommitNumber:   2,
			InterferenceStats: CommitChangesInterferenceStats{
				FileStats: []CommitFileChangesInterfereStats{
					{
						Origin: FileOrigin{
							OriginalFile:       "a.txt",
							OriginCommitNumber: 0,
						},
						InterfereStats: []FileInterfereStats{
							{
								CommitNumber: 2,
								InterferenceInfo: InterferenceInfo{
									InitialNumberOfChanges: 5,
									FinalNumberOfChanges:   5,
									InterferenceValue:      0,
								},
							},
						},
					},
					{
						Origin: FileOrigin{
							OriginalFile:       "b.txt",
							OriginCommitNumber: 0,
						},
						InterfereStats: []FileInterfereStats{
							{
								CommitNumber: 2,
								InterferenceInfo: InterferenceInfo{
									InitialNumberOfChanges: 50,
									FinalNumberOfChanges:   1,
									InterferenceValue:      0.98,
								},
							},
						},
					},
				},
			},
		},
		{
			StartCommitNumber: 2,
			EndCommitNumber:   3,
			InterferenceStats: CommitChangesInterferenceStats{
				FileStats: []CommitFileChangesInterfereStats{
					{
						Origin: FileOrigin{
							OriginalFile:       "a.txt",
							OriginCommitNumber: 0,
						},
						InterfereStats: nil,
					},
					{
						Origin: FileOrigin{
							OriginalFile:       "b.txt",
							OriginCommitNumber: 0,
						},
						InterfereStats: nil,
					},
				},
			},
		},
	}

	commitChanges := []CommitChanges{
		{
			FileChanges: []CommitFileChanges{
				{
					File: "a.txt",
					Origin: FileOrigin{
						OriginalFile:       "a.txt",
						OriginCommitNumber: 0,
					},
					Changes: nil,
				},
				{
					File: "b.txt",
					Origin: FileOrigin{
						OriginalFile:       "b.txt",
						OriginCommitNumber: 0,
					},
					Changes: nil,
				},
			},
		},
		{
			FileChanges: []CommitFileChanges{
				{
					File: "a.txt",
					Origin: FileOrigin{
						OriginalFile:       "a.txt",
						OriginCommitNumber: 0,
					},
					Changes: nil,
				},
				{
					File: "b_renamed.txt",
					Origin: FileOrigin{
						OriginalFile:       "b.txt",
						OriginCommitNumber: 0,
					},
					Changes: nil,
				},
			},
		},
		{
			FileChanges: []CommitFileChanges{
				{
					File: "a.txt",
					Origin: FileOrigin{
						OriginalFile:       "a.txt",
						OriginCommitNumber: 0,
					},
					Changes: nil,
				},
				{
					File: "b_renamed.txt",
					Origin: FileOrigin{
						OriginalFile:       "b.txt",
						OriginCommitNumber: 0,
					},
					Changes: nil,
				},
			},
		},
	}

	expected := []CommitReviewInterferenceInfo{
		{
			FileInfos: []ReviewFileInterferenceInfo{
				{
					File: "a.txt",
					Origin: FileOrigin{
						OriginalFile:       "a.txt",
						OriginCommitNumber: 0,
					},
					InterferenceInfo: InterferenceInfo{
						InitialNumberOfChanges: 100,
						FinalNumberOfChanges:   1,
						InterferenceValue:      0.99,
					},
					CommitNumber: 2,
				},
			},
		},
		{
			FileInfos: []ReviewFileInterferenceInfo{
				{
					File: "b_renamed.txt",
					Origin: FileOrigin{
						OriginalFile:       "b.txt",
						OriginCommitNumber: 0,
					},
					InterferenceInfo: InterferenceInfo{
						InitialNumberOfChanges: 50,
						FinalNumberOfChanges:   1,
						InterferenceValue:      0.98,
					},
					CommitNumber: 2,
				},
			},
		},
		{},
	}

	received := GetOptimalReviewFileInterferenceInfo(combinedCommitStats, commitChanges)

	if !reflect.DeepEqual(expected, received) {
		t.Error(fmt.Sprint("Test failed.", expected, received))
	}
}
