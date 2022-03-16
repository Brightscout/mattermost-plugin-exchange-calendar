package utils

import "encoding/base64"

func EncodeString(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func DecodeString(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetTotalNumberOfBatches(total int, maxNumRequestsPerBatch int) int {
	numOfBatches := total / maxNumRequestsPerBatch
	if total % maxNumRequestsPerBatch != 0 {
		numOfBatches++
	}

	return numOfBatches
}
