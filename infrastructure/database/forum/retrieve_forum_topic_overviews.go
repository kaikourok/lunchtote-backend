package forum

import "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *ForumRepository) RetrieveForumTopicOverviews(forumId, skip, limit int) (topics *[]model.ForumTopicOverview, topicCounts int, err error) {
	rows, err := db.Queryx(`
		SELECT
			id,
			title,
			status,
			posts,
			created_at,
			last_posted_at,
			pinned_order IS NOT NULL,
		  post_type,
			(
				CASE post_type
						WHEN 'ADMINISTRATOR' THEN '管理者'
						WHEN 'SIGNED_IN'     THEN (SELECT username FROM characters WHERE id = forum_topics.character)
						WHEN 'ANONYMOUS'     THEN COALESCE(forum_topics.name, '')
				END
			) AS name,
			character,
			identifier
		FROM
			forum_topics
		WHERE
			forum = $1
		ORDER BY
			pinned_order DESC,
			last_posted_at
		OFFSET
			$2
		LIMIT
			$3;
	`,
		forumId,
		skip,
		limit,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	topicsSlice := make([]model.ForumTopicOverview, 0, limit)
	for rows.Next() {
		var topic model.ForumTopicOverview
		err = rows.Scan(
			&topic.Id,
			&topic.Title,
			&topic.Status,
			&topic.Posts,
			&topic.CreatedAt,
			&topic.LastPostedAt,
			&topic.IsPinned,
			&topic.Sender.PostType,
			&topic.Sender.Name,
			&topic.Sender.Character,
			&topic.Sender.Identifier,
		)
		if err != nil {
			return nil, 0, err
		}

		topicsSlice = append(topicsSlice, topic)
	}

	row := db.QueryRowx(`
		SELECT
			COUNT(*)
		FROM
			forum_topics
		WHERE
			forum = $1;
	`, forumId)

	err = row.Scan(&topicCounts)
	if err != nil {
		return nil, 0, err
	}

	return &topicsSlice, topicCounts, nil
}
