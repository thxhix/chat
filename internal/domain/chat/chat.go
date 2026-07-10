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

func GetChatTitle(c *Chat) string {
	if c == nil {
		return "Unknown"
	}

	switch c.Type {
	case 1: // Личный чат
		if c.Participants != nil {
			return c.Participants.UserLogin
		}
		return "Deleted User"
	case 2: // Группа
		if c.Title != nil && *c.Title != "" {
			return *c.Title
		}
		return "Group"
	case 3: // Избранное / Сохраненки
		return "Saved Messages"
	default:
		// Безопасная обработка на случай странных типов
		if c.Title != nil {
			return *c.Title
		}
		return "Chat"
	}
}
