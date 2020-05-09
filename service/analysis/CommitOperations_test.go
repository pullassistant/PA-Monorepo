package analysis

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	diff2 "sourcegraph.com/sourcegraph/go-diff/diff"
	"testing"
)

var commitFileOperationsTestData = []struct {
	testFile string
	expected CommitOperations
}{
	{"diff_new_file_14_lines.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 14},
			}}},
	}},
	{"diff_added_2_lines.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 8, AffectedLines: 2},
			}}},
	}},
	{"diff_removed_1_line.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Removed, StartLine: 9, AffectedLines: 1},
			}}},
	}},
	{"diff_replaced_1_line.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Removed, StartLine: 8, AffectedLines: 1},
				{Type: Added, StartLine: 9, AffectedLines: 1},
			}}},
	}},
	{"diff_removed_1_line_replaced_1_line.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Removed, StartLine: 1, AffectedLines: 1},
				{Type: Added, StartLine: 2, AffectedLines: 1},
				{Type: Removed, StartLine: 14, AffectedLines: 1},
			}}},
	}},
	{"diff_replaced_1_line_added_30.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Removed, StartLine: 2, AffectedLines: 1},
				{Type: Added, StartLine: 3, AffectedLines: 1},
				{Type: Added, StartLine: 15, AffectedLines: 30},
			}}},
	}},
	{"diff_replaced_1_line_replaced_1_line_added_4_lines.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Removed, StartLine: 3, AffectedLines: 1},
				{Type: Added, StartLine: 4, AffectedLines: 1},
				{Type: Removed, StartLine: 15, AffectedLines: 1},
				{Type: Added, StartLine: 16, AffectedLines: 1},
				{Type: Added, StartLine: 31, AffectedLines: 4},
			}}},
	}},
	{"diff_removed_1_line_added_2_lines.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 4, AffectedLines: 2},
				{Type: Removed, StartLine: 42, AffectedLines: 1},
			}}},
	}},
	{"diff_new_file_7_lines.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata2.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 7},
			}}},
	}},
	{"diff_added_3_lines_added_1_line_multiple_files.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 15, AffectedLines: 3},
			}}, {
			File: "service/testdata2.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 8, AffectedLines: 1},
			}}},
	}},
	{"diff_removed_14_lines_removed_1_line_added_2_lines_multiple_files.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 8, AffectedLines: 2},
				{Type: Removed, StartLine: 29, AffectedLines: 14},
			}}, {
			File: "service/testdata2.txt", Operations: []FileOperation{
				{Type: Removed, StartLine: 7, AffectedLines: 1},
			}}},
	}},
	{"diff_removed_1_replaced_1_added_1.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: []FileOperation{
				{Type: Removed, StartLine: 5, AffectedLines: 2},
				{Type: Added, StartLine: 7, AffectedLines: 1},
				{Type: Added, StartLine: 9, AffectedLines: 1},
			}}},
	}},
	{"diff_renamed.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata.txt", Operations: nil, RenamedTo: "service/testdata_renamed.txt",
		}},
	}},
	{"diff_renamed_added_2.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata2.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 8, AffectedLines: 2},
			}, RenamedTo: "service/testdata3_renamed.txt",
		}},
	}},
	{"diff_renamed_removed_6_replaced_1.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata3_renamed.txt", Operations: []FileOperation{
				{Type: Removed, StartLine: 1, AffectedLines: 9},
			}}, {
			File: "service/testdata4_renamed.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 3},
			}}},
	}},
	{"diff_renamed_replaced_3_added_14.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata4_renamed.txt", Operations: []FileOperation{
				{Type: Removed, StartLine: 1, AffectedLines: 3},
			}}, {
			File: "service/testdata5_renamed.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 17},
			}}},
	}},
	{"diff_removed_1_added_2_lines.txt", CommitOperations{
		FileOperations: []CommitFileOperations{{
			File: "service/testdata_renamed.txt", Operations: []FileOperation{
				{Type: Added, StartLine: 1, AffectedLines: 1},
				{Type: Added, StartLine: 3, AffectedLines: 1},
				{Type: Removed, StartLine: 33, AffectedLines: 1},
				{Type: Added, StartLine: 34, AffectedLines: 1},
			}}},
	}},
}

func TestGetCommitFileOperations(t *testing.T) {
	for _, testData := range commitFileOperationsTestData {
		bytes := helperLoadBytes(t, testData.testFile)
		diff, err := diff2.ParseMultiFileDiff(bytes)
		if err != nil {
			t.Fatal("Unable to parse file: " + testData.testFile)
		}
		operations := GetCommitOperations(diff)
		if !reflect.DeepEqual(operations, testData.expected) {
			bytesExpected, _ := json.Marshal(testData.expected)
			t.Log("Expected: " + string(bytesExpected))
			bytesReceived, _ := json.Marshal(operations)
			t.Log("Received: " + string(bytesReceived))
			t.Error("Scenario " + testData.testFile + " failed")
		}
	}
}

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
