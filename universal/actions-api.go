package universal

type Actions interface {
	List() map[string]*Action
	WithName(name string) Action
	Model() *ActionsModel
}

type ActionsModel struct {
	List map[string]*ActionModel
}

func EmptyActionsModel() *ActionsModel {
	return &ActionsModel{List: make(map[string]*ActionModel)}
}

type SolidActions struct {
	actions Actions
	model   *ActionsModel
}

func (s *SolidActions) List() map[string]*Action {
	return s.actions.List()
}

func (s *SolidActions) Model() *ActionsModel {
	return s.model
}

func (s *SolidActions) WithName(name string) Action {
	return s.actions.WithName(name)
}
