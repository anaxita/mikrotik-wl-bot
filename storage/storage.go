package storage

type User struct {
	ID            int64
	Username      string
	LastMessageID int
	Status        int
}

type Storage struct {
	Users map[int64]*User
}

func NewStorage() *Storage {
	return &Storage{
		Users: make(map[int64]*User),
	}
}
