package forum

func (db *ForumRepository) CancelReactForumPost(characterId, postId int, emoji string) error {
	_, err := db.Exec(`
		DELETE FROM
			forum_topics_posts_reactions
		WHERE
			character = $1 AND
			post      = $2 AND
			emoji     = $3;
	`,
		characterId,
		postId,
		emoji,
	)

	return err
}
