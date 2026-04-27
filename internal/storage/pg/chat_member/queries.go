package chat_member

const queryIsMember = `
        SELECT EXISTS (
            SELECT 1
            FROM users_to_chats
            WHERE chat_id = $1 AND user_id = $2
        )
`
const queryLink = `INSERT INTO users_to_chats (chat_id, user_id) VALUES ($1, $2)`

const queryLinkBatch = `
	INSERT INTO users_to_chats (chat_id, user_id)
	SELECT $1, unnest($2::bigint[])
	ON CONFLICT DO NOTHING
`
