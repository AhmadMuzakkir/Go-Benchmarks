package obj_mapping_benchmarks

import (
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func getMock(t testing.TB, nRows int) *sql.DB {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}

	createdAt := time.Now()
	columns := []string{"id", "name", "email", "created_at"}
	rows := sqlmock.NewRows(columns)

	for i := 1; i <= nRows; i++ {
		rows.AddRow(
			i,
			"Username "+strconv.Itoa(i),
			"Useremail"+strconv.Itoa(i)+"@gmail.com",
			createdAt,
		)
	}

	for i := 0; i < 3_000_000; i++ {
		mock.ExpectQuery("SELECT * FROM users").WillReturnRows(rows)
	}

	return db
}

func BenchmarkScan_Manual_10000(b *testing.B) {
	mockDb := getMock(b, 10_000)
	defer mockDb.Close()

	var users []User

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			gotRows, err := mockDb.Query("SELECT * FROM users")
			if err != nil {
				b.Fatal(err)
			}

			users = make([]User, 0)

			for gotRows.Next() {
				var u User
				err := gotRows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
				if err != nil {
					b.Fatal(err)
					return
				}

				users = append(users, u)
			}
		}
	}

	_ = users[0]
}

func BenchmarkScan_Manual_20000(b *testing.B) {
	mockDb := getMock(b, 20_000)
	defer mockDb.Close()

	var users []User

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			gotRows, err := mockDb.Query("SELECT * FROM users")
			if err != nil {
				b.Fatal(err)
			}

			users = make([]User, 0)

			for gotRows.Next() {
				var u User
				err := gotRows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
				if err != nil {
					b.Fatal(err)
					return
				}

				users = append(users, u)
			}
		}
	}

	_ = users[0]
}

func BenchmarkScan_Sqlx_10000(b *testing.B) {
	mockDb := getMock(b, 10_000)
	defer mockDb.Close()

	db := sqlx.NewDb(mockDb, "postgres")
	defer db.Close()

	var users []User

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			users = make([]User, 0)

			err := db.Select(&users, "SELECT * FROM users")
			if err != nil {
				b.Fatal(err)
				return
			}
		}
	}

	_ = users[0]
}

func BenchmarkScan_Sqlx_20000(b *testing.B) {
	mockDb := getMock(b, 20_000)
	defer mockDb.Close()

	db := sqlx.NewDb(mockDb, "postgres")
	defer db.Close()

	var users []User

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			users = make([]User, 0)

			err := db.Select(&users, "SELECT * FROM users")
			if err != nil {
				b.Fatal(err)
				return
			}
		}
	}

	_ = users
}

type User struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}
