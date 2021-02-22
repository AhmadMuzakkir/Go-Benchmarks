package obj_mapping_benchmarks

import (
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func getMock(t testing.TB) *sql.DB {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}

	createdAt := time.Now()
	columns := []string{"id", "name", "email", "created_at"}
	rows := sqlmock.NewRows(columns)

	for i := 1; i <= 10000; i++ {
		rows.AddRow(
			i,
			"Username "+strconv.Itoa(i),
			"Useremail"+strconv.Itoa(i)+"@gmail.com",
			createdAt,
		)
	}

	for i := 0; i < 100000; i++ {
		mock.ExpectQuery("SELECT * FROM users").WillReturnRows(rows)
	}

	return db
}

func TestManual(t *testing.T) {
	mockDb := getMock(t)
	defer mockDb.Close()

	var user []User

	gotRows, err := mockDb.Query("SELECT * FROM users")
	if err != nil {
		t.Fatal(err)
	}

	for gotRows.Next() {
		var u User
		err := gotRows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
		if err != nil {
			t.Fatal(err)
			return
		}

		user = append(user, u)
	}

	_ = user

	t.Log(len(user))
}

func BenchmarkManual(b *testing.B) {
	mockDb := getMock(b)
	defer mockDb.Close()

	var user []User

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		gotRows, err := mockDb.Query("SELECT * FROM users")
		if err != nil {
			b.Fatal(err)
		}

		user = make([]User, 0)

		for gotRows.Next() {
			var u User
			err := gotRows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
			if err != nil {
				b.Fatal(err)
				return
			}

			user = append(user, u)
		}
	}

	_ = user
}

func BenchmarkSqlx(b *testing.B) {
	mockDb := getMock(b)
	defer mockDb.Close()

	db := sqlx.NewDb(mockDb, "postgres")
	defer db.Close()

	var user []User

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		//for j := 0; j < 100; j++ {
		user = make([]User, 0)
		err := db.Select(&user, "SELECT * FROM users")
		if err != nil {
			b.Fatal(err)
			return
		}
		//}
	}

	_ = user
}

type User struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}
