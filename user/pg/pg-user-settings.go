package postgres

import (
	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PgUserSettings struct {
	Db *pg.PgDb
	Id int64
}

func NewPgUserSettings(db *pg.PgDb, id int64) user.UserSettings {
	return &PgUserSettings{db, id}
}

func (pgUserSetting *PgUserSettings) Model() *user.UserSettingsModel {
	return nil
}

func (pgUserSetting *PgUserSettings) Radar() universal.Radar {
	return nil
}
