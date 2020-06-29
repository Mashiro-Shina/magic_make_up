package repositories

import "magicMakeup/entities"

func GetAllStars() ([]*entities.Star, error) {
	var stars []*entities.Star
	res := mysqlConn.Find(&stars)
	if res.Error != nil {
		return nil, res.Error
	}

	return stars, nil
}
