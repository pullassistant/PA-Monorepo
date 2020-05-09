package analysis

import (
	"sourcegraph.com/sourcegraph/go-diff/diff"
	"strings"
)

const gitDiffFileNoName = "dev/null"

type CommitOperations struct {
	FileOperations []CommitFileOperations
}

type CommitFileOperations struct {
	File       string
	Operations []FileOperation
	RenamedTo  string
}

type FileOperation struct {
	Type          FileOperationType
	StartLine     int
	AffectedLines int
}

type FileOperationType int

const (
	Removed FileOperationType = iota
	Added
)

func GetCommitOperations(filesDiff []*diff.FileDiff) CommitOperations {
	commitFileOperations := make([]CommitFileOperations, len(filesDiff))

	for fileIndex := range filesDiff {
		origName := getDiffFileName(filesDiff[fileIndex].OrigName)
		newName := getDiffFileName(filesDiff[fileIndex].NewName)
		var fileOperations []FileOperation

		if origName == gitDiffFileNoName {
			commitFileOperations[fileIndex].File = newName
		} else {
			commitFileOperations[fileIndex].File = origName
			if origName != newName && newName != gitDiffFileNoName {
				commitFileOperations[fileIndex].RenamedTo = newName
			}
		}

		for hunkIndex := range filesDiff[fileIndex].Hunks {
			var hunk = string(filesDiff[fileIndex].Hunks[hunkIndex].Body)
			var hunkLines = strings.Split(hunk, "\n")
			var nonAdditionLinesCount = 0
			var currentOperation *FileOperation
			var linesToSkip = 0

			for hunkLineIndex, hunkLine := range hunkLines {
				if len(hunkLine) > 0 {
					operationType := getFileOperationType(hunkLine[0])

					// fix for "No newline at end of file"
					if operationType != nil && *operationType == Removed &&
						hunkLineIndex+1 < len(hunkLines) && len(hunkLines[hunkLineIndex+1]) > 0 {
						nextOperation := getFileOperationType(hunkLines[hunkLineIndex+1][0])
						if nextOperation != nil && *nextOperation == Added &&
							hunkLine[1:] == hunkLines[hunkLineIndex+1][1:] {
							linesToSkip = 2
							nonAdditionLinesCount++
						}
					}

					if linesToSkip == 0 {
						if currentOperation == nil && operationType != nil {
							startLine := max(1, int(filesDiff[fileIndex].Hunks[hunkIndex].OrigStartLine)) + nonAdditionLinesCount
							currentOperation = &FileOperation{Type: *operationType, StartLine: startLine, AffectedLines: 1}
						} else if currentOperation != nil && operationType != nil && currentOperation.Type == *operationType {
							currentOperation.AffectedLines++
						} else if currentOperation != nil && operationType == nil {
							fileOperations = append(fileOperations, *currentOperation)
							currentOperation = nil
						} else if currentOperation != nil && operationType != nil && currentOperation.Type != *operationType {
							fileOperations = append(fileOperations, *currentOperation)
							startLine := max(1, int(filesDiff[fileIndex].Hunks[hunkIndex].OrigStartLine)) + nonAdditionLinesCount
							currentOperation = &FileOperation{Type: *operationType, StartLine: startLine, AffectedLines: 1}
						}

						if currentOperation == nil || currentOperation.Type != Added {
							nonAdditionLinesCount++
						}
					} else {
						linesToSkip--
					}
				}
			}

			if currentOperation != nil {
				fileOperations = append(fileOperations, *currentOperation)
			}

			commitFileOperations[fileIndex].Operations = fileOperations
		}
	}

	return CommitOperations{FileOperations: commitFileOperations}
}

func getDiffFileName(fileName string) string {
	split := strings.Split(fileName, "/")
	filename := split[1:]
	return strings.Join(filename, "/")
}

func getFileOperationType(sign uint8) *FileOperationType {
	var operationType = new(FileOperationType)
	if sign == '-' {
		*operationType = Removed
	} else if sign == '+' {
		*operationType = Added
	} else {
		operationType = nil
	}
	return operationType
}
