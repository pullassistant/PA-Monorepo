package analysis

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
)

var accumulatedCommitFileChangesTestData = []struct {
	commitOperations []CommitOperations
	expected         []CommitChanges
}{
	//Adding
	{
		commitOperations: []CommitOperations{
			{FileOperations: []CommitFileOperations{{
				File: "service/testdata.txt",
				Operations: []FileOperation{
					{Type: Added, StartLine: 1, AffectedLines: 14},
				},
			}}},
		},
		expected: []CommitChanges{{[]CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
			},
		}}},
		}},
	{
		commitOperations: []CommitOperations{
			{FileOperations: []CommitFileOperations{{
				File: "service/testdata.txt",
				Operations: []FileOperation{
					{Type: Added, StartLine: 1, AffectedLines: 14},
				},
			}}}, {FileOperations: []CommitFileOperations{{
				File: "service/testdata2.txt",
				Operations: []FileOperation{
					{Type: Added, StartLine: 1, AffectedLines: 1},
				},
			}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata2.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata2.txt",
				OriginCommitNumber: 1,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 1, CommitNumber: 1},
			},
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 14},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 5, AffectedLines: 1},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 4, CommitNumber: 0},
				{StartLine: 5, AffectedLines: 1, CommitNumber: 1},
				{StartLine: 6, AffectedLines: 10, CommitNumber: 0},
			},
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 14},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 1},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 1, CommitNumber: 1},
				{StartLine: 2, AffectedLines: 14, CommitNumber: 0},
			},
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 14},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 5, AffectedLines: 1},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 4, CommitNumber: 0},
				{StartLine: 5, AffectedLines: 1, CommitNumber: 1},
				{StartLine: 6, AffectedLines: 10, CommitNumber: 0},
			},
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 14},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 14, AffectedLines: 1},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 13, CommitNumber: 0},
				{StartLine: 14, AffectedLines: 1, CommitNumber: 1},
				{StartLine: 15, AffectedLines: 1, CommitNumber: 0},
			},
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 14},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 15, AffectedLines: 5},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
				{StartLine: 15, AffectedLines: 5, CommitNumber: 1},
			},
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 14},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 20, AffectedLines: 5},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 14, CommitNumber: 0},
				{StartLine: 20, AffectedLines: 5, CommitNumber: 1},
			},
		}}}},
	},
	{ // fix test
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 2},
				{Type: Added, StartLine: 4, AffectedLines: 2},
				{Type: Added, StartLine: 10, AffectedLines: 2},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 2, AffectedLines: 5},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 2, CommitNumber: 0},
				{StartLine: 4, AffectedLines: 2, CommitNumber: 0},
				{StartLine: 10, AffectedLines: 2, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 1, CommitNumber: 0},
				{StartLine: 2, AffectedLines: 5, CommitNumber: 1},
				{StartLine: 7, AffectedLines: 1, CommitNumber: 0},
				{StartLine: 9, AffectedLines: 2, CommitNumber: 0},
				{StartLine: 15, AffectedLines: 2, CommitNumber: 0},
			},
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 2},
				{Type: Added, StartLine: 4, AffectedLines: 2},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 1},
				{Type: Added, StartLine: 3, AffectedLines: 1},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 2, CommitNumber: 0},
				{StartLine: 4, AffectedLines: 2, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 1, AffectedLines: 1, CommitNumber: 1},
				{StartLine: 2, AffectedLines: 2, CommitNumber: 0},
				{StartLine: 4, AffectedLines: 1, CommitNumber: 1},
				{StartLine: 6, AffectedLines: 2, CommitNumber: 0},
			},
		}}}},
	},
	//Removing
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Removed, StartLine: 1, AffectedLines: 14},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: nil,
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 5, AffectedLines: 5},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Removed, StartLine: 1, AffectedLines: 1},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 5, AffectedLines: 5, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 4, AffectedLines: 5, CommitNumber: 0},
			},
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 5, AffectedLines: 5},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Removed, StartLine: 5, AffectedLines: 1},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 5, AffectedLines: 5, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 5, AffectedLines: 4, CommitNumber: 0},
			},
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 5, AffectedLines: 5},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Removed, StartLine: 4, AffectedLines: 6},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 5, AffectedLines: 5, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: nil,
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 5, AffectedLines: 5},
				{Type: Added, StartLine: 11, AffectedLines: 1},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Removed, StartLine: 4, AffectedLines: 8},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 5, AffectedLines: 5, CommitNumber: 0},
				{StartLine: 11, AffectedLines: 1, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: nil,
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 5, AffectedLines: 5},
				{Type: Added, StartLine: 11, AffectedLines: 1},
				{Type: Added, StartLine: 12, AffectedLines: 1},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Removed, StartLine: 4, AffectedLines: 8},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 5, AffectedLines: 5, CommitNumber: 0},
				{StartLine: 11, AffectedLines: 1, CommitNumber: 0},
				{StartLine: 12, AffectedLines: 1, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 4, AffectedLines: 1, CommitNumber: 0},
			},
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 5, AffectedLines: 5},
				{Type: Added, StartLine: 11, AffectedLines: 3},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Removed, StartLine: 4, AffectedLines: 8},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 5, AffectedLines: 5, CommitNumber: 0},
				{StartLine: 11, AffectedLines: 3, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 4, AffectedLines: 2, CommitNumber: 0},
			},
		}}}},
	},
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 5, AffectedLines: 5},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Removed, StartLine: 6, AffectedLines: 2},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 5, AffectedLines: 5, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 5, AffectedLines: 1, CommitNumber: 0},
				{StartLine: 6, AffectedLines: 2, CommitNumber: 0},
			},
		}}}},
	},
	//// Mix
	{
		commitOperations: []CommitOperations{{FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Added, StartLine: 5, AffectedLines: 5},
			},
		}}}, {FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt",
			Operations: []FileOperation{
				{Type: Removed, StartLine: 5, AffectedLines: 1},
				{Type: Added, StartLine: 6, AffectedLines: 1},
			},
		}}}},
		expected: []CommitChanges{{FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 5, AffectedLines: 5, CommitNumber: 0},
			},
		}}}, {FileChanges: []CommitFileChanges{{
			File: "service/testdata.txt",
			Origin: FileOrigin{
				OriginalFile:       "service/testdata.txt",
				OriginCommitNumber: 0,
			},
			Changes: []FileChange{
				{StartLine: 5, AffectedLines: 1, CommitNumber: 1},
				{StartLine: 6, AffectedLines: 4, CommitNumber: 0},
			},
		}}}},
	},
}

func TestGetAccumulatedCommitFileOperations(t *testing.T) {
	for i, testData := range accumulatedCommitFileChangesTestData {
		operations := GetAccumulatedCommitChanges(testData.commitOperations)

		if !reflect.DeepEqual(operations, testData.expected) {
			bytesExpected, _ := json.Marshal(testData.expected)
			t.Log("Expected: " + string(bytesExpected))
			bytesOperations, _ := json.Marshal(operations)
			t.Log("Received: " + string(bytesOperations))
			t.Error("Scenario " + strconv.Itoa(i) + " failed.")
		}
	}
}
