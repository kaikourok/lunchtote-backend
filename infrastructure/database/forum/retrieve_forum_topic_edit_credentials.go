package forum

func (db *ForumRepository) RetrieveForumTopicEditCredentials(topicId int) (masterCharacter *int, editPassword *string, postType string, err error) {
	row := db.QueryRowx(`
		SELECT
			post_type,
			character,
			edit_password
		FROM
			forum_topics
		WHERE
			id = $1;
	`, topicId)

	err = row.Scan(
		&postType,
		&masterCharacter,
		&editPassword,
	)

	return
}
