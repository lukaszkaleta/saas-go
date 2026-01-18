package pgcategory

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
)

func pgCategorySetupTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("saas-go", "category_test")
	schema := NewCategorySchema(db)
	err := schema.CreateTest()
	if err != nil {
		tb.Fatal(err)
	}
	return func(tb testing.TB) {
		err := schema.DropTest()
		if err != nil {
			tb.Fatal(err)
		}
	}, db
}

func TestPgCategory_Parent(t *testing.T) {
	teardownSuite, db := pgCategorySetupTest(t)
	defer teardownSuite(t)

	categories := PgCategories{Db: db}
	parentName := "Parent"
	categoryParent, err := categories.AddWithName(t.Context(), parentName)
	if err != nil {
		t.Fatal(err)
	}
	categoryChild, err := categories.AddWithParent(t.Context(), categoryParent, "Child")
	if err != nil {
		t.Fatal(err)
	}
	parent, err := categoryChild.Parent(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if parent.Model(t.Context()).Name.Value != parentName {
		t.Fatal("parent name does not match")
	}
}
