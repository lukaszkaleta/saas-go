package user

import (
	"github.com/lukaszkaleta/saas-go/universal"
)

// API

type UserSettings interface {
	Model() *UserSettingsModel
	Radar() universal.Radar
	Avatar() universal.Description
}

// Model

type UserSettingsModel struct {
	Radar  *universal.RadarModel       `json:"radar"`
	Avatar *universal.DescriptionModel `json:"avatar"`
}

func NewUserSettingsModel() *UserSettingsModel {
	return &UserSettingsModel{
		Radar: universal.NewRadarModel(),
	}
}

// Solid

type SolidUserSettings struct {
	Id           int64
	model        *UserSettingsModel
	userSettings UserSettings
}

func NewSolidUserSettings(model *UserSettingsModel, userSettings UserSettings, id int64) UserSettings {
	return &SolidUserSettings{
		Id:           id,
		model:        model,
		userSettings: userSettings,
	}
}

func (u SolidUserSettings) Model() *UserSettingsModel {
	return u.model
}

func (u SolidUserSettings) Radar() universal.Radar {
	if u.userSettings != nil {
		return universal.NewSolidRadar(u.Model().Radar, u.userSettings.Radar())
	}
	return universal.NewSolidRadar(u.Model().Radar, nil)
}

func (u SolidUserSettings) Avatar() universal.Description {
	if u.userSettings != nil {
		return universal.NewSolidDescription(u.Model().Avatar, u.userSettings.Avatar())
	}
	return universal.NewSolidDescription(u.Model().Avatar, nil)
}
