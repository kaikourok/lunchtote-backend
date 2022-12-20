package forum

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *ForumRepository) PostForumTopicPost(characterId *int, identifier *string, topicId int, post *model.ForumTopicPostSendData) (postId int, err error) {
	row := db.QueryRowx(`
		SELECT
			forums.force_post_type
		FROM
			forum_topics
		JOIN
			forums ON (forum_topics.id = $1 AND forum_topics.forum = forums.id);
	`, topicId)

	var forcedPostType *string
	err = row.Scan(&forcedPostType)
	if err != nil {
		return 0, err
	}

	if post.PostType != "ADMINISTRATOR" && forcedPostType != nil && *forcedPostType != post.PostType {
		return 0, errors.New("投稿タイプがフォーラム設定に合致していません")
	}

	err = db.ExecTx(func(tx *sqlx.Tx) error {
		row := tx.QueryRowx(`
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
			)
			RETURNING
				id;
		`,
			characterId,
			identifier,
			topicId,
			post.Name,
			post.EditPassword,
			post.Content,
			post.PostType,
		)
		err := row.Scan(&postId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE
				forum_topics
			SET
				posts          = posts + 1,
				last_posted_at = CURRENT_TIMESTAMP
			WHERE
				id = $1;
		`, topicId)

		return err
	})

	return
}
