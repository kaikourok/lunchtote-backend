package forum

func (db *ForumRepository) DeleteForumTopicPost(characterId *int, postId int) error {
	_, err := db.Exec(`
			UPDATE
				forum_topics_posts
			SET
				deleted_at = CURRENT_TIMESTAMP
			WHERE
				id = $1;
		`,
		postId,
	)

	return err
}
