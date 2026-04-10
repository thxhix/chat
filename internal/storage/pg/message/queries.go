package message

const queryGetByChatMessagesFirstPage = `
SELECT *
FROM messages
WHERE chat_id = $1
  AND deleted_at IS NULL
ORDER BY created_at DESC, id DESC
LIMIT $2 + 1;
`
const queryGetByChatMessagesWithCursor = `
SELECT *
FROM messages
WHERE chat_id = $1
  AND deleted_at IS NULL
  AND (created_at, id) < ($2, $3)
ORDER BY created_at DESC, id DESC
LIMIT $4 + 1;
`

const addMessageQuery = `INSERT INTO messages (message_id, chat_id, user_id, text)
	SELECT $1, $2, $3, $4
	WHERE EXISTS (
		SELECT 1 FROM users_to_chats
		WHERE chat_id = $2 AND user_id = $3
	)`
