package link

import (
	"golang/internal/stat"
	"math/rand"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	URL   string      `json:"url"`
	Hash  string      `json:"hash" gorm:"uniqueIndex"`
	Stats []stat.Stat `gorm:"constraits:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewLink(url string) *Link {
	link := &Link{
		URL: url,
	}
	link.GenerateHash(10)
	return link
}

func (l *Link) GenerateHash(n int) {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	l.Hash = string(b)
}

var letterRunes = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
