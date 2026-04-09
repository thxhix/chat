package message

const addMessageQuery = `INSERT INTO messages (message_id, chat_id, user_id, text)
	SELECT $1, $2, $3, $4
	WHERE EXISTS (
		SELECT 1 FROM users_to_chats
		WHERE chat_id = $2 AND user_id = $3
	)`
