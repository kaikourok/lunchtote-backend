package forum

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *ForumRepository) CreateForumTopic(characterId *int, identifier *string, forumId int, topic *model.ForumTopicCreateData) (topicId int, err error) {
	row := db.QueryRowx(`
		SELECT
			force_post_type
		FROM
			forums
		WHERE
			id = $1;
	`, forumId)

	var forcedPostType *string
	err = row.Scan(&forcedPostType)
	if err != nil {
		return 0, err
	}

	if topic.PostType != "ADMINISTRATOR" && forcedPostType != nil && *forcedPostType != topic.PostType {
		return 0, errors.New("投稿タイプがフォーラム設定に合致していません")
	}

	err = db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
			INSERT INTO forum_topics (
				character,
				identifier,
				forum,
				title,
				name,
				edit_password,
				post_type
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7
			)
			RETURNING
				id;
		`,
			characterId,
			identifier,
			forumId,
			topic.Title,
			topic.Name,
			topic.EditPassword,
			topic.PostType,
		)
		err := row.Scan(&topicId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO forum_topics_posts (
				character,
				identifier,
				topic,
				name,
				edit_password,
				content,
				post_type
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7
			);
		`,
			characterId,
			identifier,
			topicId,
			topic.Name,
			topic.EditPassword,
			topic.Content,
			topic.PostType,
		)
		return err
	})
	if err != nil {
		return 0, err
	}

	return topicId, nil
}
