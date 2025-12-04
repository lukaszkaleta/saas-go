package pguser

import (
	"testing"

	"github.com/lukaszkaleta/saas-go/database/pg"
	pgfilestore "github.com/lukaszkaleta/saas-go/filestore/pg"
	"github.com/lukaszkaleta/saas-go/universal"
)

func pgUserSetupTest(tb testing.TB) (func(tb testing.TB), *pg.PgDb) {
	db := pg.LocalPgWithName("saas-go", "user-test")
	pgfilestore.NewFilestoreSchema(db).Create()
	NewUserSchema(db).Create()

	return func(tb testing.TB) {
		NewUserSchema(db).Drop()
		pgfilestore.NewFilestoreSchema(db).Drop()
	}, db
}

func TestPgUser_UpdateSettings(t *testing.T) {
	teardownSuite, db := pgUserSetupTest(t)
	defer teardownSuite(t)

	users := PgUsers{Db: db}
	personModel := &universal.PersonModel{Phone: "01234"}
	user, err := users.Add(personModel)
	if err != nil {
		t.Fatal(err)
	}
	radarModel := &universal.RadarModel{
		Position: &universal.PositionModel{
			Lat: 49.7765893956647,
			Lon: 21.661360264888657,
		},
		Perimeter: 10,
	}
	user.Settings().Radar().Update(radarModel)
}
