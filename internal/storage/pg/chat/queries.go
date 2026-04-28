package chat

const queryGetByUUID = `SELECT id, idempotency_key, type, title, created_at FROM chats WHERE idempotency_key = $1`
const queryGetIDByUUID = `SELECT id FROM chats WHERE idempotency_key = $1`

const queryCreateChat = `INSERT INTO chats (idempotency_key, type) 
	VALUES ($1, $2) 
	ON CONFLICT (idempotency_key) 
	DO UPDATE SET 
		idempotency_key = EXCLUDED.idempotency_key
	RETURNING id, idempotency_key, (XMAX = 0) AS is_created;`

const queryGetUserChats = `
	SELECT c.*
	FROM chats c
	JOIN users_to_chats utc ON c.id = utc.chat_id
	WHERE utc.user_id = $1
	ORDER BY c.created_at DESC
	LIMIT $2
`
