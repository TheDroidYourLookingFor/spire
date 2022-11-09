package models

import (
	"github.com/volatiletech/null/v8"
)

type FactionBaseDatum struct {
	ClientFactionId   int16      `json:"client_faction_id" gorm:"Column:client_faction_id"`
	Min               int16 `json:"min" gorm:"Column:min"`
	Max               int16 `json:"max" gorm:"Column:max"`
	UnkHero1          int16 `json:"unk_hero_1" gorm:"Column:unk_hero1"`
	UnkHero2          int16 `json:"unk_hero_2" gorm:"Column:unk_hero2"`
	UnkHero3          int16 `json:"unk_hero_3" gorm:"Column:unk_hero3"`
}

func (FactionBaseDatum) TableName() string {
    return "faction_base_data"
}

func (FactionBaseDatum) Relationships() []string {
    return []string{}
}

func (FactionBaseDatum) Connection() string {
    return "eqemu_content"
}
