package store

// Store -- app repositories
type Store interface {
	User() UserRepository
	Chat() ChatRepository
	LearningNode() LearningNodeRepository
	Document() DocumentRepository
}
