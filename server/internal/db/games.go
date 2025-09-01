package db

func (s *Store) CreateGame(game *Game) error {
	return s.db.Create(game).Error
}

func (s *Store) GetGameByID(id string) (*Game, error) {
	var game Game
	if err := s.db.First(&game, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &game, nil
}

func (s *Store) GetActiveGames() ([]Game, error) {
	var games []Game
	if err := s.db.Find(&games).Where("status = ?", "active").Error; err != nil {
		return nil, err
	}
	return games, nil
}

func (s *Store) UpdateGame(game *Game) error {
	return s.db.Save(game).Error
}

func (s *Store) DeleteGameByID(id string) error {
	return s.db.Delete(&Game{}, "id = ?", id).Error
}
