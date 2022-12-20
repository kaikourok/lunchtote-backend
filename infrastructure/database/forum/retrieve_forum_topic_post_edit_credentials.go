package forum

func (db *ForumRepository) RetrieveForumTopicPostEditCredentials(postId int) (masterCharacter *int, editPassword *string, postType string, err error) {
	row := db.QueryRowx(`
		SELECT
			post_type,
			character,
			edit_password
		FROM
			forum_topics_posts
		WHERE
			id = $1;
	`, postId)

	err = row.Scan(
		&postType,
		&masterCharacter,
		&editPassword,
	)

	return
}
