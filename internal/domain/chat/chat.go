package chat

import "fmt"

func GenerateIdempotencyStr(cType int8, userIDs []int64) string {
	id1, id2 := getSortedPair(userIDs[0], userIDs[1])
	return getKeyStr(cType, id1, id2)
}

func getSortedPair(a, b int64) (int64, int64) {
	if a < b {
		return a, b
	}
	return b, a
}

func getKeyStr(cType int8, id1, id2 int64) string {
	return fmt.Sprintf("%d:%d:%d", cType, id1, id2)
}
