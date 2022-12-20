package forum

import (
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *ForumRepository) UpdateForumTopic(characterId *int, topicId int, topic *model.ForumTopicEditData) error {
	_, err := db.Exec(`
		UPDATE
			forum_topics
		SET
			title,
			status
		WHERE
			id = $1;
	`,
		topicId,
		topic.Title,
		topic.Status,
	)

	return err
}
