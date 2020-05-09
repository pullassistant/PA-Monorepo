package analysis

type CommitChanges struct {
	FileChanges []CommitFileChanges
}

type CommitFileChanges struct {
	File    string
	Origin  FileOrigin
	Changes []FileChange
}

type FileOrigin struct {
	OriginalFile       string
	OriginCommitNumber int
}

type FileChange struct {
	StartLine     int
	AffectedLines int
	CommitNumber  int
}

func GetAccumulatedCommitChanges(commitOperations []CommitOperations) []CommitChanges {
	var result []CommitChanges = nil

	for commitNumber, singleCommitFileOperations := range commitOperations {
		var singleCommitFileChanges []CommitFileChanges = nil

		for _, fileOperations := range singleCommitFileOperations.FileOperations {
			var fileChanges []FileChange
			var fileOrigin FileOrigin

			previousFileChanges := getPreviousFileChanges(result, fileOperations.File)
			if previousFileChanges == nil {
				fileOrigin.OriginCommitNumber = commitNumber
				fileOrigin.OriginalFile = fileOperations.File
				fileChanges = mergeFileChangesAndOperations(nil, fileOperations.Operations, commitNumber)
			} else {
				fileOrigin = previousFileChanges.Origin
				if len(previousFileChanges.Changes) == 0 {
					getPreviousFileChanges(nil, "")
				}
				fileChanges = mergeFileChangesAndOperations(previousFileChanges.Changes, fileOperations.Operations, commitNumber)
			}

			fileName := fileOperations.File
			if fileOperations.RenamedTo != "" {
				fileName = fileOperations.RenamedTo
			}

			singleCommitFileChanges = append(singleCommitFileChanges, CommitFileChanges{
				File:    fileName,
				Origin:  fileOrigin,
				Changes: fileChanges,
			})
		}

		result = append(result, CommitChanges{singleCommitFileChanges})
	}
	return result
}

func getPreviousFileChanges(commitChanges []CommitChanges, file string) *CommitFileChanges {
	for i := len(commitChanges) - 1; i >= 0; i-- {
		for j := range commitChanges[i].FileChanges {
			if commitChanges[i].FileChanges[j].File == file {
				return &commitChanges[i].FileChanges[j]
			}
		}
	}

	return nil
}

func mergeFileChangesAndOperations(previousChanges []FileChange, operations []FileOperation, commitNumber int) []FileChange {
	var newFileChanges []FileChange
	var linesShift int
	var remainingOperation *FileOperation
	var remainingChange *FileChange
	var operationRange rangeInfo
	var changeRange rangeInfo

	// todo support empty commit

	if len(previousChanges) > 0 {
		remainingChange, previousChanges = &previousChanges[0], previousChanges[1:]
	}

	if len(operations) > 0 { //operations can be empty if the file is renamed without any changes
		remainingOperation, operations = &operations[0], operations[1:]
	}

	for remainingChange != nil || remainingOperation != nil {

		if remainingChange != nil {
			changeRange = rangeInfo{
				start: remainingChange.StartLine + linesShift,
				end:   remainingChange.StartLine + remainingChange.AffectedLines + linesShift,
			}
		}

		if remainingOperation != nil {
			operationRange = rangeInfo{
				start: remainingOperation.StartLine + linesShift,
				end:   remainingOperation.StartLine + remainingOperation.AffectedLines + linesShift,
			}
		}

		if remainingChange == nil && remainingOperation != nil {
			// Operations below all changes
			// *
			// *
			// -+/

			if remainingOperation.Type != Removed {
				newFileChanges = append(newFileChanges, FileChange{StartLine: operationRange.start, AffectedLines: operationRange.length(), CommitNumber: commitNumber})
				linesShift += operationRange.length()
			} else {
				linesShift -= operationRange.length()
			}
			remainingOperation = nil
		}

		if remainingOperation == nil && remainingChange != nil {
			// Changes below operations
			// -+/
			// -+/
			// *

			newFileChanges = append(newFileChanges, FileChange{StartLine: changeRange.start, AffectedLines: changeRange.length(), CommitNumber: remainingChange.CommitNumber})
			remainingChange = nil
		}

		if remainingChange != nil && remainingOperation != nil {

			if remainingOperation.Type == Added {
				if operationRange.isStartBefore(changeRange) {
					// Add operation started before change
					// +
					// *+
					// *+

					newFileChanges = append(newFileChanges, FileChange{StartLine: operationRange.start, AffectedLines: operationRange.length(), CommitNumber: commitNumber})
					linesShift += operationRange.length()
					remainingOperation = nil
				} else if operationRange.isAfter(changeRange) {
					// Operations below change
					// *
					// *
					// +

					newFileChanges = append(newFileChanges, FileChange{StartLine: changeRange.start, AffectedLines: changeRange.length(), CommitNumber: remainingChange.CommitNumber})
					remainingChange = nil
				} else {
					// Add operation started before change ends
					// *
					// *+
					// *+
					// +
					// ---
					// *
					// +
					// *

					preservedTopChange := changeRange.startToStartDistance(operationRange)
					if preservedTopChange > 0 {
						newFileChanges = append(newFileChanges, FileChange{StartLine: changeRange.start, AffectedLines: preservedTopChange, CommitNumber: remainingChange.CommitNumber})
					}

					newFileChanges = append(newFileChanges, FileChange{StartLine: operationRange.start, AffectedLines: operationRange.length(), CommitNumber: commitNumber})

					changesAffected := changeRange.endToStartDistance(operationRange)
					if changesAffected > 0 {
						remainingChange = &FileChange{StartLine: remainingChange.StartLine + preservedTopChange, AffectedLines: changesAffected, CommitNumber: remainingChange.CommitNumber}
					} else {
						remainingChange = nil
					}

					linesShift += operationRange.length()
					remainingOperation = nil
				}
			} else if remainingOperation.Type == Removed {

				if operationRange.isBefore(changeRange) {
					// Remove operation started and finished before change
					// -
					// *
					// *

					linesShift -= operationRange.length()
					remainingOperation = nil
				} else if operationRange.isAfter(changeRange) {
					// Operations below change
					// *
					// *
					// -

					newFileChanges = append(newFileChanges, FileChange{StartLine: changeRange.start, AffectedLines: changeRange.length(), CommitNumber: remainingChange.CommitNumber})
					remainingChange = nil
				} else {
					// Remove operation started before change ends
					// *
					// *-
					// *-
					// -
					// ---
					// *
					// -
					// *

					preservedTopChange := 0
					if changeRange.isStartBefore(operationRange) {
						preservedTopChange = changeRange.startToStartDistance(operationRange)
						newFileChanges = append(newFileChanges, FileChange{StartLine: changeRange.start, AffectedLines: preservedTopChange, CommitNumber: remainingChange.CommitNumber})
					}

					preservedOperation := 0
					if changeRange.isEndBeforeEnd(operationRange) {
						preservedOperation = changeRange.endToEndDistance(operationRange)
						remainingOperation = &FileOperation{Type: remainingOperation.Type, StartLine: changeRange.end, AffectedLines: preservedOperation}
					} else {
						remainingOperation = nil
					}

					if operationRange.isEndBeforeEnd(changeRange) {
						// todo czy - line shift jest dobre?
						remainingChange = &FileChange{StartLine: operationRange.end - linesShift, AffectedLines: changeRange.endToEndDistance(operationRange), CommitNumber: remainingChange.CommitNumber}
					} else {
						remainingChange = nil
					}

					linesShift -= operationRange.length() - preservedOperation
				}
			}
		}

		if len(previousChanges) > 0 && remainingChange == nil {
			remainingChange, previousChanges = &previousChanges[0], previousChanges[1:]
		}

		if len(operations) > 0 && remainingOperation == nil {
			remainingOperation, operations = &operations[0], operations[1:]

		}
	}

	return newFileChanges
}
