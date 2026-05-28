package stat

import (
	"golang/packages/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{Db: db}
}

func (repo *StatRepository) AddClick(linkID uint) {
	var stat Stat
	currentDate := datatypes.Date(time.Now())
	repo.DB.Find(&stat, "link_id = ? and date = ?", linkID, currentDate)
	if stat.ID == 0 {
		repo.DB.Create(&Stat{
			LinkId: linkID,
			Clicks: 1,
			Date:   currentDate,
		})
	} else {
		stat.Clicks++
		repo.DB.Save(&stat)
	}
}

func (repo *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponce {
	var stats []GetStatResponce
	var selectQuery string
	switch by {
	case FilterByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case FilterByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, clicks"
	}
	repo.DB.Table("stats").Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Group("clicks").
		Scan(&stats)
	return stats
}
