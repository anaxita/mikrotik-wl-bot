package storage

const (
	anaxitaUsername  = "anaxita"
	mishaglUsername  = "Mishagl"
	kmsControlChatID = -1001700493413
)

type User struct {
	ID            int64
	Username      string
	LastMessageID int
	Status        int
	Role          int
}

type Storage struct {
	Users        map[int64]*User
	Admins       map[string]struct{}
	AllowChatIDs []int64
}

func NewStorage() *Storage {
	return &Storage{
		Users: make(map[int64]*User),
		Admins: map[string]struct{}{
			anaxitaUsername: {},
			mishaglUsername: {},
		},
		AllowChatIDs: []int64{kmsControlChatID},
	}
}

func (s *Storage) AddAdmin(username string) {
	s.Admins[username] = struct{}{}
}

func (s *Storage) RemoveAdmin(username string) {
	delete(s.Admins, username)
}
