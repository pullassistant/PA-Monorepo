package analysis

// if 1 - set of commit should be reviewed commit by commit
func GetCommitByCommitReviewScore(interferenceStats []CommitChangesInterferenceStats) float32 {
	initialNumberOfChanges := 0
	finalNumberOfChanges := 0

	// single commit
	if len(interferenceStats) == 1 {
		return 1
	}

	for _, commitStats := range interferenceStats {
		for _, fileStats := range commitStats.FileStats {

			var maxInterference *InterferenceInfo = nil
			for _, interferenceStat := range fileStats.InterfereStats {

				if maxInterference == nil || maxInterference.InterferenceValue < interferenceStat.InterferenceInfo.InterferenceValue {
					interference := interferenceStat.InterferenceInfo
					maxInterference = &interference
				}
			}

			if maxInterference != nil {
				initialNumberOfChanges += maxInterference.InitialNumberOfChanges
				if maxInterference.InterferenceValue == 0 {
					finalNumberOfChanges += maxInterference.InitialNumberOfChanges
				} else {
					finalNumberOfChanges += maxInterference.FinalNumberOfChanges
				}
			}
		}
	}

	return float32(finalNumberOfChanges) / float32(initialNumberOfChanges)
}
