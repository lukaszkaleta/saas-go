package universal

import (
	"regexp"
	"strings"

	"github.com/lukaszkaleta/saas-go/universal"
)

// API

type Name interface {
	Model() *NameModel
	Update(newModel *NameModel) error
}

// Builder

func NameFromModel(model *NameModel) Name {
	return SolidName{
		model: model,
		Name:  nil,
	}
}

// Model

type NameModel struct {
	Value string `json:"value"`
	Slug  string `json:"slug"`
}

func EmptyNameModel() *NameModel {
	return &NameModel{
		Value: "",
		Slug:  "",
	}
}

func SluggedName(name string) *NameModel {
	return &NameModel{
		Value: name,
		Slug:  CreateSlug(name),
	}
}

func (model *NameModel) Change(name string) {
	model.Value = name
	model.Slug = CreateSlug(name)
}

// Solid

type SolidName struct {
	model *NameModel
	Name  Name
}

func NewSolidName(model *NameModel, Name Name) Name {
	return &SolidName{model, Name}
}

func (addr SolidName) Update(newModel *NameModel) error {
	addr.model.Change(newModel.Value)
	if addr.Name == nil {
		return nil
	}
	return addr.Name.Update(addr.model)
}

func (addr SolidName) Model() *NameModel {
	return addr.model
}

// createSlug converts a name into a URL-friendly slug
func CreateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove all non-alphanumeric and non-hyphen characters
	re := regexp.MustCompile("[^a-z0-9-]+")
	slug = re.ReplaceAllString(slug, "")

	// Replace multiple hyphens with a single one
	reHyphen := regexp.MustCompile("-+")
	slug = reHyphen.ReplaceAllString(slug, "-")

	// Trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}
