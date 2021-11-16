package storage

const (
	anaxitaChatID = 404907387
	mishaglChatID = 2102715403
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
	Admins       []string
	AllowChatIDs []int64
}

func NewStorage() *Storage {
	return &Storage{
		Users:        make(map[int64]*User),
		Admins:       []string{"anaxita", "Mishagl"},
		AllowChatIDs: []int64{anaxitaChatID, mishaglChatID, -1001287143568},
	}
}
