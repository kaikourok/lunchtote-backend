package forum

func (db *ForumRepository) CreateForumGroup(title string) (id int, err error) {
	row := db.QueryRowx(`
		INSERT INTO forum_groups (
			title,
			forum_groups_order
		) VALUES (
			$1,
			(COALESCE((SELECT forum_groups_order FROM forum_groups ORDER BY forum_groups_order DESC LIMIT 1), 0) + 1)
		)
		RETURNING
			id;
	`, title)

	err = row.Scan(&id)
	return
}
