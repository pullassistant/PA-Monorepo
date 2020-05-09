package analysis

const minCommitOverlapToCombine = 0.50

// represents a set of combined commits [StartCommitNumber, EndCommitNumber)
type CombinedCommitsInterferenceStats struct {
	StartCommitNumber int
	EndCommitNumber   int
	InterferenceStats CommitChangesInterferenceStats
}

func GetOptimalReview(interferenceStats []CommitChangesInterferenceStats) []CombinedCommitsInterferenceStats {
	combinedCommits := getSingularCombinedCommits(interferenceStats)
	optimised := false

	for optimised != true {
		optimised = true

		for i := 0; i < len(combinedCommits)-1; i++ {
			commitsFrom := &combinedCommits[i]
			commitsTo := &combinedCommits[i+1]

			combinedCommitsInterferenceInfo := commitsFrom.InterferenceStats.GetInterferenceWithFollowingCommits(commitsTo.StartCommitNumber, commitsTo.EndCommitNumber)

			if combinedCommitsInterferenceInfo.hasAnyInterference() {
				if combinedCommitsInterferenceInfo.InterferenceValue > minCommitOverlapToCombine {
					combinedCommitsInterferenceStats := getCombinedCommitsInterferenceStats(commitsFrom.StartCommitNumber, commitsTo.EndCommitNumber, interferenceStats)
					commitsTo.StartCommitNumber = commitsFrom.StartCommitNumber
					commitsTo.InterferenceStats = combinedCommitsInterferenceStats
					optimised = false
					combinedCommits = append(combinedCommits[:i], combinedCommits[i+1:]...)
					break
				}
			}
		}
	}

	return combinedCommits
}

func getCombinedCommitsInterferenceStats(startCommitNumber, endCommitNumber int, interferenceStats []CommitChangesInterferenceStats) CommitChangesInterferenceStats {
	var newFileStats = make([]CommitFileChangesInterfereStats, len(interferenceStats[endCommitNumber-1].FileStats))
	copy(newFileStats, interferenceStats[endCommitNumber-1].FileStats)

	for i := range newFileStats {
		var newInterferenceStats = make([]FileInterfereStats, len(newFileStats[i].InterfereStats))
		copy(newInterferenceStats, newFileStats[i].InterfereStats)
		newFileStats[i].InterfereStats = newInterferenceStats
	}

	for i := startCommitNumber; i < endCommitNumber-1; i++ {
		for _, fileStat := range interferenceStats[i].FileStats {

			var previousFileStat []FileInterfereStats = nil
			for _, prevFileStat := range newFileStats {
				if prevFileStat.Origin == fileStat.Origin {
					previousFileStat = prevFileStat.InterfereStats
				}
			}

			indexDiff := endCommitNumber - i - 1
			var maxInterferenceInfo *InterferenceInfo = nil
			for m := 0; m < indexDiff; m++ {
				if fileStat.InterfereStats[m].InterferenceInfo.hasAnyInterference() {
					maxInterferenceInfo = &fileStat.InterfereStats[m].InterferenceInfo
				}
			}

			if maxInterferenceInfo != nil && maxInterferenceInfo.InterferenceValue == 1 {
				continue
			}

			if previousFileStat != nil {
				for k := 0; k < len(previousFileStat); k++ {
					var interferenceInfo = fileStat.InterfereStats[indexDiff+k].InterferenceInfo
					if maxInterferenceInfo != nil {
						interferenceInfo.InitialNumberOfChanges = maxInterferenceInfo.FinalNumberOfChanges
					}

					newInterferenceInfo := previousFileStat[k].InterferenceInfo
					if !newInterferenceInfo.hasAnyInterference() && !interferenceInfo.hasAnyInterference() {
						newInterferenceInfo.InitialNumberOfChanges += interferenceInfo.InitialNumberOfChanges
					} else if !newInterferenceInfo.hasAnyInterference() && interferenceInfo.hasAnyInterference() {
						newInterferenceInfo.FinalNumberOfChanges = newInterferenceInfo.InitialNumberOfChanges + interferenceInfo.FinalNumberOfChanges
						newInterferenceInfo.InitialNumberOfChanges += interferenceInfo.InitialNumberOfChanges
						newInterferenceInfo.calculateInterferenceValue()
					} else if newInterferenceInfo.hasAnyInterference() && !interferenceInfo.hasAnyInterference() {
						newInterferenceInfo.FinalNumberOfChanges += interferenceInfo.InitialNumberOfChanges
						newInterferenceInfo.InitialNumberOfChanges += interferenceInfo.InitialNumberOfChanges
						newInterferenceInfo.calculateInterferenceValue()
					} else {
						newInterferenceInfo.FinalNumberOfChanges += interferenceInfo.FinalNumberOfChanges
						newInterferenceInfo.InitialNumberOfChanges += interferenceInfo.InitialNumberOfChanges
						newInterferenceInfo.calculateInterferenceValue()
					}

					previousFileStat[k].InterferenceInfo = newInterferenceInfo
				}
			} else {
				var newInterferenceStats = make([]FileInterfereStats, len(fileStat.InterfereStats))
				copy(newInterferenceStats, fileStat.InterfereStats)

				newInterferenceStats = newInterferenceStats[indexDiff:]
				if maxInterferenceInfo != nil {
					for m := range newInterferenceStats {
						newInterferenceStats[m].InterferenceInfo.InitialNumberOfChanges = maxInterferenceInfo.FinalNumberOfChanges
						if newInterferenceStats[m].InterferenceInfo.hasAnyInterference() {
							newInterferenceStats[m].InterferenceInfo.calculateInterferenceValue()
						}
					}
				}

				var newFileStat = CommitFileChangesInterfereStats{
					Origin:         fileStat.Origin,
					InterfereStats: newInterferenceStats,
				}

				newFileStats = append(newFileStats, newFileStat)
			}
		}
	}

	return CommitChangesInterferenceStats{
		FileStats: newFileStats,
	}
}

func getSingularCombinedCommits(interferenceStats []CommitChangesInterferenceStats) []CombinedCommitsInterferenceStats {
	var result []CombinedCommitsInterferenceStats
	for i := 0; i < len(interferenceStats); i++ {
		result = append(result, CombinedCommitsInterferenceStats{
			StartCommitNumber: i,
			EndCommitNumber:   i + 1,
			InterferenceStats: interferenceStats[i],
		})
	}
	return result
}
