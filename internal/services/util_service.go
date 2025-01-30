package services

import (
	"time"
)

type UtilService struct {
}

func NewUtilService() *UtilService {
	return &UtilService{}
}

func (utilService *UtilService) GetNumberEmoji(number int) string {
	emojis := map[int]string{
		1:  "1️⃣",
		2:  "2️⃣",
		3:  "3️⃣",
		4:  "4️⃣",
		5:  "5️⃣",
		6:  "6️⃣",
		7:  "7️⃣",
		8:  "8️⃣",
		9:  "9️⃣",
		10: "🔟",
		0:  "0️⃣",
	}
	return emojis[number]
}

func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02T15:04:05Z", date)
	return err == nil
}
