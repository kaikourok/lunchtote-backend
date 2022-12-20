package forum

import (
	"encoding/json"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *ForumRepository) RetrieveForumTopicPosts(topicId int, characterId *int) (posts *[]model.ForumTopicPost, err error) {
	rows, err := db.Queryx(`
		WITH
			post_revisions AS (
				SELECT
					forum_topics_posts_revisions.post      AS post_id,
					forum_topics_posts_revisions.content   AS content,
					forum_topics_posts_revisions.posted_at AS posted_at
				FROM
					forum_topics_posts
				JOIN
					forum_topics_posts_revisions ON (forum_topics_posts.id = forum_topics_posts_revisions.post)
				WHERE
					forum_topics_posts.topic = $1
			),
			post_reactions AS (
				SELECT
					forum_topics_posts_reactions.post  AS post_id,
					forum_topics_posts_reactions.emoji AS emoji,
					COUNT(forum_topics_posts_reactions.character) AS reacted_counts,
					MIN(forum_topics_posts_reactions.reacted_at)  AS first_reacted_at,
					(
						CASE
							WHEN $2::INT IS     NULL THEN false
							WHEN $2::INT IS NOT NULL THEN BOOL_OR(forum_topics_posts_reactions.character = $2)
						END
					) AS is_reacted
				FROM
					forum_topics_posts
				JOIN
					forum_topics_posts_reactions ON (forum_topics_posts.id = forum_topics_posts_reactions.post)
				WHERE
					forum_topics_posts.topic = $1
				GROUP BY
					forum_topics_posts_reactions.post,
					forum_topics_posts_reactions.emoji
			)

		SELECT
			forum_topics_posts.id,
			forum_topics_posts.content,
			forum_topics_posts.posted_at,
			forum_topics_posts.updated_at,
			COALESCE(
				(
					SELECT
						JSON_AGG(rev)
					FROM
						(
							SELECT
								post_revisions.content   AS content,
								post_revisions.posted_at AS postedAt
							FROM
								post_revisions
							WHERE
								post_revisions.post_id = forum_topics_posts.id
							ORDER BY
								post_revisions.posted_at
						) rev
				),
				'[]'
			),
			COALESCE(
				(
					SELECT
						JSON_AGG(rea)
					FROM
						(
							SELECT
								post_reactions.emoji          AS emoji,
								post_reactions.reacted_counts AS reactedCounts,
								post_reactions.is_reacted     AS isReacted
							FROM
								post_reactions
							WHERE
								post_reactions.post_id = forum_topics_posts.id
							ORDER BY
								post_reactions.first_reacted_at
						) rea
				),
				'[]'
			),
			forum_topics_posts.post_type,
			(
				CASE post_type
						WHEN 'ADMINISTRATOR' THEN '管理者'
						WHEN 'SIGNED_IN'     THEN (SELECT username FROM characters WHERE id = forum_topics_posts.character)
						WHEN 'ANONYMOUS'     THEN COALESCE(forum_topics_posts.name, '')
				END
			) AS name,
			forum_topics_posts.character,
			forum_topics_posts.identifier
		FROM
			forum_topics_posts
		WHERE
			forum_topics_posts.topic = $1
		ORDER BY
			forum_topics_posts.posted_at;
	`, topicId, characterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postsSlice := make([]model.ForumTopicPost, 0, 64)
	for rows.Next() {
		var post model.ForumTopicPost
		var revisionsJson, reactionsJson string
		err = rows.Scan(
			&post.Id,
			&post.Content,
			&post.PostedAt,
			&post.UpdatedAt,
			&revisionsJson,
			&reactionsJson,
			&post.Sender.PostType,
			&post.Sender.Name,
			&post.Sender.Character,
			&post.Sender.Identifier,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(revisionsJson), &post.Revisions)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(reactionsJson), &post.Reactions)
		if err != nil {
			return nil, err
		}

		postsSlice = append(postsSlice, post)
	}

	return &postsSlice, nil
}
