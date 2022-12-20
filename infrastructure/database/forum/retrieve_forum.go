package forum

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *ForumRepository) RetrieveForum(forumId int) (forum *model.Forum, err error) {
	row := db.QueryRowx(`
		SELECT
			id,
			title,
			summary,
			guide,
			force_post_type
		FROM
			forums
		WHERE
			id = $1;
	`, forumId)

	var forumData model.Forum
	err = row.Scan(
		&forumData.Id,
		&forumData.Title,
		&forumData.Summary,
		&forumData.Guide,
		&forumData.ForcedPostType,
	)
	if err != nil {
		return nil, err
	}

	return &forumData, nil
}
