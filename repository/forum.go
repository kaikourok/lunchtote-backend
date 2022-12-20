package repository

import "github.com/kaikourok/lunchtote-backend/entity/model"

type forumRepository interface {
	CreateForumGroup(title string) (id int, err error)
	CreateForum(forumGroupId int, forum *model.ForumCreateData) (id int, err error)
	RetrieveForumOverviews() (forumGroups *[]model.ForumGroup, err error)
	RetrieveForum(forumId int) (forum *model.Forum, err error)
	RetrieveForumForcedPostType(forumId int) (forcedPostType *string, err error)
	RetrieveForumTopicOverviews(forumId, skip, limit int) (topics *[]model.ForumTopicOverview, topicCounts int, err error)
	RetrieveForumTopic(topicId int) (topic *model.ForumTopic, err error)
	RetrieveForumTopicPosts(topicId int, characterId *int) (posts *[]model.ForumTopicPost, err error)
	CreateForumTopic(characterId *int, identifier *string, forumId int, topic *model.ForumTopicCreateData) (topicId int, err error)
	RetrieveForumTopicEditCredentials(topicId int) (masterCharacter *int, editPassword *string, postType string, err error)
	UpdateForumTopic(characterId *int, topicId int, topic *model.ForumTopicEditData) error
	PostForumTopicPost(characterId *int, identifier *string, topicId int, post *model.ForumTopicPostSendData) (postId int, err error)
	RetrieveForumTopicPostEditCredentials(postId int) (masterCharacter *int, editPassword *string, postType string, err error)
	ReviseForumTopicPost(characterId *int, postId int, post *model.ForumTopicPostReviseData) error
	DeleteForumTopicPost(characterId *int, postId int) error
	ReactForumPost(characterId, postId int, emoji string) error
	CancelReactForumPost(characterId, postId int, emoji string) error
}
