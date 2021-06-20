package user

import (
	"database/sql"

	sqlutil "github.com/alexyslozada/gosqlutils"
	"github.com/jpastorm/dialogflowbot/infraestructure/postgres"
	"github.com/jpastorm/dialogflowbot/model"
)

const table = "users"

var fields = []string{
	"name",
	"telegram_id",
}

var (
	psqlInsert = postgres.BuildSQLInsert(table, fields)
	psqlUpdate = postgres.BuildSQLUpdateByID(table, fields)
	psqlDelete = "DELETE FROM " + table + " WHERE id = $1"
	psqlGetAll = postgres.BuildSQLSelect(table, fields)
)

// CronLog struct that implement the interface domain.cronlog.Storage
type User struct {
	db *sql.DB
}

// New return a new CronLog
func New(db *sql.DB) User {
	return User{db}
}

// Create this method creates a model.User in postgres db
func (c User) Create(m *model.User) error {
	stmt, err := c.db.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		m.Name,
		m.TelegramID,
	).Scan(&m.ID, &m.CreatedAt)
	if err != nil {
		if errPsql := postgres.CheckError(err); errPsql != nil {
			return errPsql
		}

		return err
	}

	return nil
}

// Update this method updates a model.User in postgres db
func (c User) Update(m *model.User) error {
	stmt, err := c.db.Prepare(psqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = sqlutil.ExecAffectingOneRow(
		stmt,
		m.Name,
		m.TelegramID,
		m.ID,
	)
	if err != nil {
		if errPsql := postgres.CheckError(err); errPsql != nil {
			return errPsql
		}
		return err
	}

	return nil
}

// Delete this method deletes a model.User by ID in postgres db
func (c User) Delete(ID uint) error {
	stmt, err := c.db.Prepare(psqlDelete)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return sqlutil.ExecAffectingOneRow(stmt, ID)
}

// GetWhere this method get a model.User with filters in postgres  db
func (c User) GetWhere(filter model.Fields, sort model.SortFields) (model.User, error) {
	condition, args := postgres.BuildSQLWhere(filter)
	query := psqlGetAll + " " + condition

	query += " " + postgres.BuildSQLOrderBy(sort)

	stmt, err := c.db.Prepare(query)
	if err != nil {
		return model.User{}, err
	}
	defer stmt.Close()

	return c.scanRow(stmt.QueryRow(args...))
}

// GetAllWhere this method get all ordered model.CronLogs with filters in postgres db
func (c User) GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Users, error) {
	conditions, args := postgres.BuildSQLWhere(filter)
	query := psqlGetAll + " " + conditions

	query += " " + postgres.BuildSQLOrderBy(sort)
	query += " " + postgres.BuildSQLPagination(pag)

	stmt, err := c.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := make(model.Users, 0)
	for rows.Next() {
		m, err := c.scanRow(rows)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (c User) scanRow(s sqlutil.RowScanner) (model.User, error) {
	user := model.User{}
	updatedAtNull := sql.NullTime{}

	err := s.Scan(
		&user.ID,
		&user.Name,
		&user.TelegramID,
		&user.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return user, err
	}

	user.UpdatedAt = updatedAtNull.Time

	return user, nil
}
