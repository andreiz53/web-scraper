package db

import "github.com/andreiz53/web-scraper/types"

func (s *SqliteStorage) CreateWebsite(w *types.Website) error {
	result := s.Db.Create(w)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SqliteStorage) GetWebsitesByUserID(userID uint) ([]types.Website, error) {
	var websites []types.Website
	result := s.Db.Where("user_id = ?", userID).Find(&websites)
	if result.Error != nil {
		return nil, result.Error
	}
	return websites, nil
}
