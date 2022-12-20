package forum

func (db *ForumRepository) RetrieveForumForcedPostType(forumId int) (forcedPostType *string, err error) {
	row := db.QueryRowx(`
		SELECT
			force_post_type
		FROM
			forums
		WHERE
			id = $1;
	`, forumId)

	err = row.Scan(&forcedPostType)
	return
}
