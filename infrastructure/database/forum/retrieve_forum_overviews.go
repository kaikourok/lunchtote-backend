package forum

import (
	"encoding/json"
	"time"

	"github.com/kaikourok/lunchtote-backend/entity/model"
)

func (db *ForumRepository) RetrieveForumOverviews() (forumGroups *[]model.ForumGroup, err error) {
	rows, err := db.Queryx(`
		WITH forum_last_posts AS (
			SELECT
				last_posts.forum_id          AS forum_id,
				forum_topics.id              AS last_post_topic_id,
				forum_topics.title           AS last_post_topic_title,
				forum_topics_posts.posted_at AS last_post_posted_at,
				forum_topics_posts.post_type AS last_post_post_type,
				(
					CASE forum_topics_posts.post_type
						WHEN 'ADMINISTRATOR' THEN '管理者'
						WHEN 'SIGNED_IN'     THEN characters.username
						WHEN 'ANONYMOUS'     THEN COALESCE(forum_topics_posts.name, '')
					END
				) AS last_post_name,
				(
					CASE forum_topics_posts.post_type
						WHEN 'ADMINISTRATOR' THEN NULL
						WHEN 'SIGNED_IN'     THEN characters.id
						WHEN 'ANONYMOUS'     THEN NULL
					END
				) AS last_post_character
			FROM
				(
					SELECT
						forums.id                  AS forum_id,
						MAX(forum_topics_posts.id) AS last_post_id
					FROM
						forums
					JOIN
						forum_topics ON (forums.id = forum_topics.forum)
					JOIN
						forum_topics_posts ON (forum_topics.id = forum_topics_posts.topic AND forum_topics_posts.deleted_at IS NULL)
					GROUP BY
						forums.id
				) AS last_posts
			JOIN
				forum_topics_posts ON (last_posts.last_post_id = forum_topics_posts.id)
			JOIN
				forum_topics ON (forum_topics.id = forum_topics_posts.topic)
			LEFT JOIN
				characters ON (forum_topics_posts.character = characters.id)
		)

		SELECT
			t.forum_group_id,
			t.forum_group_title,
			JSON_AGG(JSON_BUILD_OBJECT(
				'forum_id',              t.forum_id,
				'forum_title',           t.forum_title,
				'forum_summary',         t.forum_summary,
				'last_post_topic_id',    t.last_post_topic_id,
				'last_post_topic_title', t.last_post_topic_title,
				'last_post_posted_at',   TO_CHAR(t.last_post_posted_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
				'last_post_post_type',   t.last_post_post_type,
				'last_post_name',        t.last_post_name,
				'last_post_character',   t.last_post_character
			))
		FROM
			(
				SELECT
					forum_groups.id    AS forum_group_id,
					forum_groups.title AS forum_group_title,
					forums.id          AS forum_id,
					forums.title       AS forum_title,
					forums.summary     AS forum_summary,
					forum_last_posts.last_post_topic_id,
					forum_last_posts.last_post_topic_title,
					forum_last_posts.last_post_posted_at,
					forum_last_posts.last_post_post_type,
					forum_last_posts.last_post_name,
					forum_last_posts.last_post_character
				FROM
					forum_groups
				JOIN
					forums ON (forums.forum_group = forum_groups.id)
				LEFT JOIN
					forum_last_posts ON (forum_last_posts.forum_id = forums.id)
				ORDER BY
					forum_groups.forum_groups_order,
					forums.forum_order
			) AS t
		GROUP BY
			t.forum_group_id,
			t.forum_group_title;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]model.ForumGroup, 0, 16)
	for rows.Next() {
		var forumGroup model.ForumGroup
		var forumsJson string
		rows.Scan(
			&forumGroup.Id,
			&forumGroup.Title,
			&forumsJson,
		)

		var forumsJsonReader []struct {
			ForumId            int        `json:"forum_id"`
			ForumTitle         string     `json:"forum_title"`
			ForumSummary       string     `json:"forum_summary"`
			LastPostTopicId    *int       `json:"last_post_topic_id"`
			LastPostTopicTitle *string    `json:"last_post_topic_title"`
			LastPostPostedAt   *time.Time `json:"last_post_posted_at"`
			LastPostPostType   *string    `json:"last_post_post_type"`
			LastPostName       *string    `json:"last_post_name"`
			LastPostCharacter  *int       `json:"last_post_character"`
		}

		err = json.Unmarshal([]byte(forumsJson), &forumsJsonReader)
		if err != nil {
			return nil, err
		}

		forums := make([]model.ForumOverview, len(forumsJsonReader))
		for i := range forums {
			forums[i].Id = forumsJsonReader[i].ForumId
			forums[i].Title = forumsJsonReader[i].ForumTitle
			forums[i].Summary = forumsJsonReader[i].ForumSummary

			if forumsJsonReader[i].LastPostTopicId != nil &&
				forumsJsonReader[i].LastPostTopicTitle != nil &&
				forumsJsonReader[i].LastPostPostedAt != nil &&
				forumsJsonReader[i].LastPostPostType != nil &&
				forumsJsonReader[i].LastPostName != nil {
				var lastPost model.ForumOverviewLastPost

				lastPost.Topic.Id = *forumsJsonReader[i].LastPostTopicId
				lastPost.Topic.Title = *forumsJsonReader[i].LastPostTopicTitle
				lastPost.Sender.PostType = *forumsJsonReader[i].LastPostPostType
				lastPost.Sender.Name = *forumsJsonReader[i].LastPostName
				lastPost.Sender.Character = forumsJsonReader[i].LastPostCharacter
				lastPost.PostedAt = *forumsJsonReader[i].LastPostPostedAt

				forums[i].LastPost = &lastPost
			}
		}

		forumGroup.Forums = forums
		groups = append(groups, forumGroup)
	}

	return &groups, nil
}
