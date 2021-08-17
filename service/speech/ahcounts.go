package speech

import (
	"errors"

	"github.com/voice-of-colombo/service/database"
	"github.com/voice-of-colombo/service/graph/model"
	"gorm.io/gorm"
)

func GetAhCountsBySpeechId(id string) (*model.AhCounts, error) {
	var ahCounts *model.AhCounts
	result := database.Db.First(&ahCounts, "speech_id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return ahCounts, nil
}
