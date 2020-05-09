package analysis

const minInterferenceToAssumeHigh = 0.4

type CommitReviewInterferenceInfo struct {
	FileInfos []ReviewFileInterferenceInfo
}

type ReviewFileInterferenceInfo struct {
	File             string
	Origin           FileOrigin
	InterferenceInfo InterferenceInfo
	CommitNumber     int
}

func GetOptimalReviewFileInterferenceInfo(combinedCommitStats []CombinedCommitsInterferenceStats, commitChanges []CommitChanges) []CommitReviewInterferenceInfo {
	var result []CommitReviewInterferenceInfo

	for _, combinedCommitStat := range combinedCommitStats {
		var filesInfos []ReviewFileInterferenceInfo

		for _, fileStat := range combinedCommitStat.InterferenceStats.FileStats {
			var maxFileInterfereStats *FileInterfereStats

			for _, interferenceStat := range fileStat.InterfereStats {
				if maxFileInterfereStats == nil || maxFileInterfereStats.InterferenceInfo.InterferenceValue < interferenceStat.InterferenceInfo.InterferenceValue {
					tmp := interferenceStat
					maxFileInterfereStats = &tmp
				}
			}

			if maxFileInterfereStats != nil && maxFileInterfereStats.InterferenceInfo.InterferenceValue >= minInterferenceToAssumeHigh {

				var file = ""
				for i := combinedCommitStat.EndCommitNumber - 1; i >= combinedCommitStat.StartCommitNumber; i-- {
					for _, v := range commitChanges[i].FileChanges {
						if v.Origin == fileStat.Origin {
							file = v.File
						}
					}
				}

				if len(file) == 0 {
					panic("GetOptimalReviewFileInterferenceInfo - empty file name")
				}

				fileInterferenceInfo := ReviewFileInterferenceInfo{
					File:             file,
					Origin:           fileStat.Origin,
					InterferenceInfo: maxFileInterfereStats.InterferenceInfo,
					CommitNumber:     maxFileInterfereStats.CommitNumber,
				}
				filesInfos = append(filesInfos, fileInterferenceInfo)
			}
		}
		result = append(result, CommitReviewInterferenceInfo{FileInfos: filesInfos})
	}

	return result
}
