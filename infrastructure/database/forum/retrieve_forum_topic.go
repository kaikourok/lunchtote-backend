package forum

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *ForumRepository) RetrieveForumTopic(topicId int) (topic *model.ForumTopic, err error) {
	row := db.QueryRowx(`
		SELECT
			forum_topics.id,
			forum_topics.forum,
			forum_topics.title,
			forum_topics.status,
			forums.force_post_type,
		  forum_topics.post_type,
			(
				CASE post_type
						WHEN 'ADMINISTRATOR' THEN '管理者'
						WHEN 'SIGNED_IN'     THEN (SELECT username FROM characters WHERE id = forum_topics.character)
						WHEN 'ANONYMOUS'     THEN COALESCE(forum_topics.name, '')
				END
			) AS name,
			forum_topics.character,
			forum_topics.identifier
		FROM
			forum_topics
		JOIN
			forums ON (forum_topics.forum = forums.id)
		WHERE
			forum_topics.id = $1;
	`, topicId)

	var topicData model.ForumTopic
	err = row.Scan(
		&topicData.Id,
		&topicData.Forum,
		&topicData.Title,
		&topicData.Status,
		&topicData.ForcePostType,
		&topicData.Sender.PostType,
		&topicData.Sender.Name,
		&topicData.Sender.Character,
		&topicData.Sender.Identifier,
	)
	if err != nil {
		return nil, err
	}

	return &topicData, nil
}
