package pgfilestore

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
	"github.com/lukaszkaleta/saas-go/filestore"
	"github.com/lukaszkaleta/saas-go/universal"
)

func setupTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("saas-go", "filestore_test")
	schema := NewFilestoreSchema(db)
	schema.Create()

	return func(tb testing.TB) {
		schema.Drop()
	}, db
}

func TestPgRecord_Update(t *testing.T) {
	teardownSuite, db := setupTest(t)
	defer teardownSuite(t)

	records := PgRecords{Db: db}
	record, err := records.Add(t.Context(), filestore.EmptyRecordModel())
	if err != nil {
		t.Fatal(err)
	}
	err = record.Update(t.Context(), &filestore.RecordModel{Name: &universal.NameModel{Value: "koza"}, Description: &universal.DescriptionModel{}})
	if err != nil {
		t.Fatal(err)
	}
	record, err = records.ById(t.Context(), record.Model(t.Context()).Id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPgRecord_Model(t *testing.T) {
	teardownSuite, db := setupTest(t)
	defer teardownSuite(t)

	records := PgRecords{Db: db}
	recordModel := &filestore.RecordModel{
		Name:        universal.SluggedName("file-name"),
		Description: &universal.DescriptionModel{Value: "file-description", ImageUrl: "file-image-url"},
	}
	record, err := records.Add(t.Context(), recordModel)
	if err != nil {
		t.Fatal(err)
	}
	record, err = records.ById(t.Context(), record.Model(t.Context()).Id)
	if err != nil {
		t.Fatal(err)
	}
	model := record.Model(t.Context())
	Equal[string](t, model.Name.Value, recordModel.Name.Value)
	Equal[string](t, model.Name.Slug, recordModel.Name.Slug)
	Equal[string](t, model.Description.Value, recordModel.Description.Value)
	Equal[string](t, model.Description.ImageUrl, recordModel.Description.ImageUrl)
}

func Equal[V comparable](t *testing.T, got, expected V) {
	t.Helper()

	if expected != got {
		t.Errorf(`assert.Equal(
t,
got:
%v
,
expected:
%v
)`, got, expected)
	}
}
