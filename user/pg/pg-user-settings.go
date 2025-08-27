package postgres

import (
	"naborly/internal/api/common"
	"naborly/internal/api/user"
)

type PgUserSettings struct {
	Db *PgDb
	Id int64
}

func NewPgUserSettings(db *PgDb, id int64) user.UserSettings {
	return &PgUserSettings{db, id}
}

func (pgUserSetting *PgUserSettings) Model() *user.UserSettingsModel {
	return nil
}

func (pgUserSetting *PgUserSettings) Radar() common.Radar {
	return nil
}
