package forum

import (
	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *ForumRepository) ReviseForumTopicPost(characterId *int, postId int, post *model.ForumTopicPostReviseData) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		_, err := db.Exec(`
			WITH old_post AS (
				SELECT
					content,
					(CASE WHEN updated_at IS NULL THEN posted_at ELSE updated_at END) AS posted_at
				FROM
					forum_topics_posts
				WHERE
					id = $1
			)

			INSERT INTO forum_topics_posts_revisions (
				post,
				content,
				updated_at
			) VALUES (
				$1,
				(SELECT content   FROM old_post),
				(SELECT posted_at FROM posted_at)
			);
		`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`
			UPDATE
				forum_topics_posts
			SET
				content    = $2,
				updated_at = CURRENT_TIMESTAMP
			WHERE
				id = $1;
		`,
			postId,
			post.Content,
		)

		return err
	})
}
