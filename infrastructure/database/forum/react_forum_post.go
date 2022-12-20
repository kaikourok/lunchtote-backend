package forum

func (db *ForumRepository) ReactForumPost(characterId, postId int, emoji string) error {
	_, err := db.Exec(`
		INSERT INTO forum_topics_posts_reactions (
			character,
			post,
			emoji
		) VALUES (
			$1,
			$2,
			$3
		)
		ON CONFLICT DO NOTHING;
	`,
		characterId,
		postId,
		emoji,
	)

	return err
}
