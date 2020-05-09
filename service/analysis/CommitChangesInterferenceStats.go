package analysis

type CommitChangesInterferenceStats struct {
	FileStats []CommitFileChangesInterfereStats
}

type CommitFileChangesInterfereStats struct {
	Origin         FileOrigin
	InterfereStats []FileInterfereStats
}

type FileInterfereStats struct {
	CommitNumber     int
	InterferenceInfo InterferenceInfo
}

type InterferenceInfo struct {
	InitialNumberOfChanges int
	FinalNumberOfChanges   int
	InterferenceValue      float32
}

func (info *InterferenceInfo) calculateInterferenceValue() {
	info.InterferenceValue = 1 - float32(info.FinalNumberOfChanges)/float32(info.InitialNumberOfChanges)
}

func (info InterferenceInfo) hasAnyInterference() bool {
	return !(info.FinalNumberOfChanges == 0 && info.InterferenceValue == 0)
}

// [startCommitNumber, endCommitNumber)
func (stat CommitChangesInterferenceStats) GetInterferenceWithFollowingCommits(startCommitNumber, endCommitNumber int) InterferenceInfo {
	initialNumberOfChanges := 0
	finalNumberOfChanges := 0

	for _, fileStat := range stat.FileStats {
		var maxInterference InterferenceInfo
		var interferenceFound = false

		for _, fileInterference := range fileStat.InterfereStats {
			for commitNumber := startCommitNumber; commitNumber < endCommitNumber; commitNumber++ {
				if fileInterference.CommitNumber == commitNumber {
					if fileInterference.InterferenceInfo.hasAnyInterference() {
						if !interferenceFound || maxInterference.InterferenceValue < fileInterference.InterferenceInfo.InterferenceValue {
							maxInterference = fileInterference.InterferenceInfo
							interferenceFound = true
						}
					}
				}
			}
		}

		if interferenceFound {
			initialNumberOfChanges += maxInterference.InitialNumberOfChanges
			finalNumberOfChanges += maxInterference.FinalNumberOfChanges
		} else {
			initialNumberOfChanges += fileStat.InterfereStats[0].InterferenceInfo.InitialNumberOfChanges
			finalNumberOfChanges += fileStat.InterfereStats[0].InterferenceInfo.InitialNumberOfChanges
		}
	}

	info := InterferenceInfo{
		InitialNumberOfChanges: initialNumberOfChanges,
		FinalNumberOfChanges:   finalNumberOfChanges,
	}

	//to keep the InterferenceInfo contact
	if initialNumberOfChanges != finalNumberOfChanges {
		info.calculateInterferenceValue()
	} else {
		info.FinalNumberOfChanges = 0
		info.InterferenceValue = 0
	}

	return info
}

func GetCommitFileChangesInterferenceStats(commitChanges []CommitChanges) []CommitChangesInterferenceStats {
	var result []CommitChangesInterferenceStats = nil

	for fromCommitNumber := 0; fromCommitNumber < len(commitChanges); fromCommitNumber++ {
		var commitFileChangesInterferenceStats []CommitFileChangesInterfereStats = nil
		fromFileOriginToCountMap := getFileOriginToChangesFromCommitMap(commitChanges[fromCommitNumber], fromCommitNumber)

		for fileOrigin := range fromFileOriginToCountMap {
			commitFileChangesInterferenceStats = append(commitFileChangesInterferenceStats, CommitFileChangesInterfereStats{
				Origin:         fileOrigin,
				InterfereStats: nil,
			})
		}

		if fromCommitNumber+1 < len(commitChanges) {
			for toCommitNumber := fromCommitNumber + 1; toCommitNumber < len(commitChanges); toCommitNumber++ {
				toFileOriginToCountMap := getFileOriginToChangesFromCommitMap(commitChanges[toCommitNumber], fromCommitNumber)

				for i := range commitFileChangesInterferenceStats {
					fromFileOrigin := commitFileChangesInterferenceStats[i].Origin
					fromCount := fromFileOriginToCountMap[fromFileOrigin]
					toCount, ok := toFileOriginToCountMap[fromFileOrigin]

					if !ok {
						toCount = 0
					}

					info := InterferenceInfo{
						InitialNumberOfChanges: fromCount,
						FinalNumberOfChanges:   toCount,
					}

					if ok {
						info.calculateInterferenceValue()
					}

					// todo log this
					//if (fromCount < toCount) {
					//	fromCount += 0
					//}

					commitFileChangesInterferenceStats[i].InterfereStats = append(commitFileChangesInterferenceStats[i].InterfereStats, FileInterfereStats{
						CommitNumber:     toCommitNumber,
						InterferenceInfo: info,
					})
				}
			}
		}

		result = append(result, CommitChangesInterferenceStats{commitFileChangesInterferenceStats})
	}

	return result
}

func getFileOriginToChangesFromCommitMap(commitFileChanges CommitChanges, commitNumber int) map[FileOrigin]int {
	result := map[FileOrigin]int{}

	for _, v := range commitFileChanges.FileChanges {
		count := 0
		for _, v2 := range v.Changes {
			if v2.CommitNumber == commitNumber {
				count += v2.AffectedLines
			}
		}
		if count > 0 {
			result[v.Origin] = count
		}
	}

	return result
}
