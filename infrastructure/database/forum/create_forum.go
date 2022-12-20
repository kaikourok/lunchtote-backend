package forum

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *ForumRepository) CreateForum(forumGroupId int, forum *model.ForumCreateData) (id int, err error) {
	row := db.QueryRowx(`
		INSERT INTO forums (
			forum_group,
			title,
			summary,
			guide,
			force_post_type,
			forum_order
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			(COALESCE((SELECT forum_order FROM forums ORDER BY forum_order DESC LIMIT 1), 0) + 1)
		)
		RETURNING
			id;
	`,
		forumGroupId,
		forum.Title,
		forum.Summary,
		forum.Guide,
		forum.ForcePostType,
	)

	err = row.Scan(&id)
	return
}
