package analysis

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetOptimalReview(t *testing.T) {
	interferenceStats := []CommitChangesInterferenceStats{
		{
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
								InitialNumberOfChanges: 50,
								FinalNumberOfChanges:   5,
								InterferenceValue:      0.9,
							},
						},
						{
							CommitNumber: 2,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 50,
								FinalNumberOfChanges:   0,
								InterferenceValue:      0,
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
						OriginalFile:       "a.txt",
						OriginCommitNumber: 0,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 2,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 10,
								FinalNumberOfChanges:   9,
								InterferenceValue:      0.1,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "b.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 2,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 50,
								FinalNumberOfChanges:   0,
								InterferenceValue:      0,
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
						OriginalFile:       "c.txt",
						OriginCommitNumber: 0,
					},
					InterfereStats: nil,
				},
			},
		},
	}

	expected := []CombinedCommitsInterferenceStats{
		{
			StartCommitNumber: 0,
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
									InitialNumberOfChanges: 15,
									FinalNumberOfChanges:   14,
									InterferenceValue:      0.06666666,
								},
							},
						},
					},
					{
						Origin: FileOrigin{
							OriginalFile:       "b.txt",
							OriginCommitNumber: 1,
						},
						InterfereStats: []FileInterfereStats{
							{
								CommitNumber: 2,
								InterferenceInfo: InterferenceInfo{
									InitialNumberOfChanges: 50,
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
			StartCommitNumber: 2,
			EndCommitNumber:   3,
			InterferenceStats: CommitChangesInterferenceStats{
				FileStats: []CommitFileChangesInterfereStats{
					{
						Origin: FileOrigin{
							OriginalFile:       "c.txt",
							OriginCommitNumber: 0,
						},
						InterfereStats: nil,
					},
				},
			},
		},
	}

	received := GetOptimalReview(interferenceStats)

	if !reflect.DeepEqual(expected, received) {
		t.Error(fmt.Sprint("Test failed.", expected, received))
	}

}

func TestGetCombinedCommitsInterferenceStats(t *testing.T) {
	interferenceStats := []CommitChangesInterferenceStats{
		{
			FileStats: []CommitFileChangesInterfereStats{
				{
					Origin: FileOrigin{
						OriginalFile:       "a.txt",
						OriginCommitNumber: 0,
					},
					InterfereStats: nil,
				},
			},
		},
		{
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
								InitialNumberOfChanges: 10,
								FinalNumberOfChanges:   9,
								InterferenceValue:      0.1,
							},
						},
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 10,
								FinalNumberOfChanges:   8,
								InterferenceValue:      0.2,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "b.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 2,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 20,
								FinalNumberOfChanges:   18,
								InterferenceValue:      0.1,
							},
						},
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 20,
								FinalNumberOfChanges:   16,
								InterferenceValue:      0.2,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "d.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 2,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 20,
								FinalNumberOfChanges:   10,
								InterferenceValue:      0.5,
							},
						},
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 20,
								FinalNumberOfChanges:   0,
								InterferenceValue:      0,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "e.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 2,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 30,
								FinalNumberOfChanges:   0,
								InterferenceValue:      0,
							},
						},
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 30,
								FinalNumberOfChanges:   0,
								InterferenceValue:      0,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "f.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 2,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 40,
								FinalNumberOfChanges:   0,
								InterferenceValue:      0,
							},
						},
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 40,
								FinalNumberOfChanges:   0,
								InterferenceValue:      0,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "g.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 2,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 40,
								FinalNumberOfChanges:   10,
								InterferenceValue:      0.75,
							},
						},
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 40,
								FinalNumberOfChanges:   0,
								InterferenceValue:      0,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "h.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 2,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 10,
								FinalNumberOfChanges:   0,
								InterferenceValue:      1,
							},
						},
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 10,
								FinalNumberOfChanges:   0,
								InterferenceValue:      0,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "i.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 2,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 10,
								FinalNumberOfChanges:   0,
								InterferenceValue:      1,
							},
						},
						{
							CommitNumber: 3,
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
		{
			FileStats: []CommitFileChangesInterfereStats{
				{
					Origin: FileOrigin{
						OriginalFile:       "c.txt",
						OriginCommitNumber: 2,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 30,
								FinalNumberOfChanges:   26,
								InterferenceValue:      0.2,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "b.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 40,
								FinalNumberOfChanges:   32,
								InterferenceValue:      0.2,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "d.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 20,
								FinalNumberOfChanges:   0,
								InterferenceValue:      0,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "e.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 30,
								FinalNumberOfChanges:   0,
								InterferenceValue:      0,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "f.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 3,
							InterferenceInfo: InterferenceInfo{
								InitialNumberOfChanges: 30,
								FinalNumberOfChanges:   27,
								InterferenceValue:      0.1,
							},
						},
					},
				},
				{
					Origin: FileOrigin{
						OriginalFile:       "i.txt",
						OriginCommitNumber: 1,
					},
					InterfereStats: []FileInterfereStats{
						{
							CommitNumber: 3,
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
	}

	expected := CommitChangesInterferenceStats{
		FileStats: []CommitFileChangesInterfereStats{
			{
				Origin: FileOrigin{
					OriginalFile:       "c.txt",
					OriginCommitNumber: 2,
				},
				InterfereStats: []FileInterfereStats{
					{
						CommitNumber: 3,
						InterferenceInfo: InterferenceInfo{
							InitialNumberOfChanges: 30,
							FinalNumberOfChanges:   26,
							InterferenceValue:      0.2,
						},
					},
				},
			},
			{
				Origin: FileOrigin{
					OriginalFile:       "b.txt",
					OriginCommitNumber: 1,
				},
				InterfereStats: []FileInterfereStats{
					{
						CommitNumber: 3,
						InterferenceInfo: InterferenceInfo{
							InitialNumberOfChanges: 58,
							FinalNumberOfChanges:   48,
							InterferenceValue:      0.17241377,
						},
					},
				},
			},
			{
				Origin: FileOrigin{
					OriginalFile:       "d.txt",
					OriginCommitNumber: 1,
				},
				InterfereStats: []FileInterfereStats{
					{
						CommitNumber: 3,
						InterferenceInfo: InterferenceInfo{
							InitialNumberOfChanges: 30,
							FinalNumberOfChanges:   0,
							InterferenceValue:      0,
						},
					},
				},
			},
			{
				Origin: FileOrigin{
					OriginalFile:       "e.txt",
					OriginCommitNumber: 1,
				},
				InterfereStats: []FileInterfereStats{
					{
						CommitNumber: 3,
						InterferenceInfo: InterferenceInfo{
							InitialNumberOfChanges: 60,
							FinalNumberOfChanges:   0,
							InterferenceValue:      0,
						},
					},
				},
			},
			{
				Origin: FileOrigin{
					OriginalFile:       "f.txt",
					OriginCommitNumber: 1,
				},
				InterfereStats: []FileInterfereStats{
					{
						CommitNumber: 3,
						InterferenceInfo: InterferenceInfo{
							InitialNumberOfChanges: 70,
							FinalNumberOfChanges:   67,
							InterferenceValue:      0.04285717,
						},
					},
				},
			},
			{
				Origin: FileOrigin{
					OriginalFile:       "i.txt",
					OriginCommitNumber: 1,
				},
				InterfereStats: []FileInterfereStats{
					{
						CommitNumber: 3,
						InterferenceInfo: InterferenceInfo{
							InitialNumberOfChanges: 10,
							FinalNumberOfChanges:   0,
							InterferenceValue:      0,
						},
					},
				},
			},
			{
				Origin: FileOrigin{
					OriginalFile:       "a.txt",
					OriginCommitNumber: 0,
				},
				InterfereStats: []FileInterfereStats{
					{
						CommitNumber: 3,
						InterferenceInfo: InterferenceInfo{
							InitialNumberOfChanges: 9,
							FinalNumberOfChanges:   8,
							InterferenceValue:      0.111111104,
						},
					},
				},
			},
			{
				Origin: FileOrigin{
					OriginalFile:       "g.txt",
					OriginCommitNumber: 1,
				},
				InterfereStats: []FileInterfereStats{
					{
						CommitNumber: 3,
						InterferenceInfo: InterferenceInfo{
							InitialNumberOfChanges: 10,
							FinalNumberOfChanges:   0,
							InterferenceValue:      0,
						},
					},
				},
			},
		},
	}

	received := getCombinedCommitsInterferenceStats(1, 3, interferenceStats)

	if !reflect.DeepEqual(expected, received) {
		t.Error(fmt.Sprint("Test failed.", expected, received))
	}
}

func TestGetSingularCombinedCommits(t *testing.T) {
	interferenceStats := []CommitChangesInterferenceStats{
		{
			FileStats: []CommitFileChangesInterfereStats{
				{
					Origin: FileOrigin{
						OriginalFile:       "a.txt",
						OriginCommitNumber: 0,
					},
					InterfereStats: nil,
				},
			},
		},
		{
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
						OriginCommitNumber: 1,
					},
					InterfereStats: nil,
				},
			},
		},
	}

	expected := []CombinedCommitsInterferenceStats{
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
						InterfereStats: nil,
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
						InterfereStats: nil,
					},
					{
						Origin: FileOrigin{
							OriginalFile:       "b.txt",
							OriginCommitNumber: 1,
						},
						InterfereStats: nil,
					},
				},
			},
		},
	}

	received := getSingularCombinedCommits(interferenceStats)

	if !reflect.DeepEqual(expected, received) {
		t.Error(fmt.Sprint("Test failed.", expected, received))
	}
}
