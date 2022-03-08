package utils

import "encoding/base64"

func EncodeString(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func GetTotalNumberOfBatches(total int, maxNumRequestsPerBatch int) int {
	numOfBatches := total / maxNumRequestsPerBatch
	if total % maxNumRequestsPerBatch != 0 {
		numOfBatches++
	}

	return numOfBatches
}
