package repository

type Repository interface {
	characterRepository
	diaryRepository
	generalRepository
	mailRepository
	roomRepository
	forumRepository
}
