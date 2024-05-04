package memdb

type Store struct {
	db []string
}

// Создание нового хранилища
func New() (*Store, error) {
	return &Store{db: []string{
		"qwerty",
		"йцукен",
		"zxcvbn",
		"asdfgh",
	}}, nil
}

func (s *Store) Words() []string {
	return s.db
}
