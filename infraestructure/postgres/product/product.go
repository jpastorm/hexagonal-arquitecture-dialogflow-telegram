package product

import (
	"database/sql"

	sqlutil "github.com/alexyslozada/gosqlutils"
	"github.com/jpastorm/dialogflowbot/infraestructure/postgres"
	"github.com/jpastorm/dialogflowbot/model"
)

const table = "products"

var fields = []string{
	"name",
	"price",
}

var (
	psqlInsert = postgres.BuildSQLInsert(table, fields)
	psqlUpdate = postgres.BuildSQLUpdateByID(table, fields)
	psqlDelete = "DELETE FROM " + table + " WHERE id = $1"
	psqlGetAll = postgres.BuildSQLSelect(table, fields)
)

// Products struct that implement the interface domain.products.Storage
type Product struct {
	db *sql.DB
}

// New return a new Product
func New(db *sql.DB) Product {
	return Product{db}
}

// Create this method creates a model.Product in postgres db
func (p Product) Create(m *model.Product) error {
	stmt, err := p.db.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		m.Name,
		m.Price,
	).Scan(&m.ID, &m.CreatedAt)
	if err != nil {
		if errPsql := postgres.CheckError(err); errPsql != nil {
			return errPsql
		}

		return err
	}

	return nil
}

// Update this method updates a model.Product in postgres db
func (p Product) Update(m *model.Product) error {
	stmt, err := p.db.Prepare(psqlUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = sqlutil.ExecAffectingOneRow(
		stmt,
		m.Name,
		m.Price,
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

// Delete this method deletes a model.Product by ID in postgres db
func (p Product) Delete(ID uint) error {
	stmt, err := p.db.Prepare(psqlDelete)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return sqlutil.ExecAffectingOneRow(stmt, ID)
}

// GetWhere this method get a model.Product with filters in postgres  db
func (p Product) GetWhere(filter model.Fields, sort model.SortFields) (model.Product, error) {
	condition, args := postgres.BuildSQLWhere(filter)
	query := psqlGetAll + " " + condition

	query += " " + postgres.BuildSQLOrderBy(sort)

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return model.Product{}, err
	}
	defer stmt.Close()

	return p.scanRow(stmt.QueryRow(args...))
}

// GetAllWhere this method get all ordered model.Products with filters in postgres db
func (p Product) GetAllWhere(filter model.Fields, sort model.SortFields, pag model.Pagination) (model.Products, error) {
	conditions, args := postgres.BuildSQLWhere(filter)
	query := psqlGetAll + " " + conditions

	query += " " + postgres.BuildSQLOrderBy(sort)
	query += " " + postgres.BuildSQLPagination(pag)

	stmt, err := p.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := make(model.Products, 0)
	for rows.Next() {
		m, err := p.scanRow(rows)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func (p Product) scanRow(s sqlutil.RowScanner) (model.Product, error) {
	product := model.Product{}
	updatedAtNull := sql.NullTime{}

	err := s.Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return product, err
	}

	product.UpdatedAt = updatedAtNull.Time

	return product, nil
}
