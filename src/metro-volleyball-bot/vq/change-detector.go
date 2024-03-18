package vq

import "fmt"

func DetectLadderChanges(old GetLadderResponseBody, new GetLadderResponseBody) bool {
	// check the lengths are the same...
	if len(old.Records) != len(new.Records) {
		return true
	}

	// slowish algo. If records are in order then the running time should be (On)
	// worst case (O)n2
	// check to see if the positions have changed.
	for i := 0; i < len(old.Records); i++ {
		recordFound := false
		oldRecord := old.Records[i]

		for j := 0; j < len(new.Records); j++ {
			newRecord := new.Records[j]

			if oldRecord.Fields.TeamNameLookup != newRecord.Fields.TeamNameLookup {
				continue
			}

			recordFound = true

			if oldRecord.Fields.Rank != newRecord.Fields.Rank {
				return true
			}

			if oldRecord.Fields.CompetitionPoints != newRecord.Fields.CompetitionPoints {
				return true
			}

			break
		}

		if !recordFound {
			fmt.Println("Record Not Found", oldRecord)
			return true
		}
	}

	// check for points changing
	return false
}
