package postgres

import (
	"testing"
	"time"

	"github.com/jpastorm/dialogflowbot/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildSQLInsert(t *testing.T) {
	tableTest := []struct {
		table  string
		fields []string
		want   string
	}{
		{
			table:  "cashboxes",
			fields: []string{"responsable", "country", "user_id", "account"},
			want:   "INSERT INTO cashboxes (responsable,country,user_id,account) VALUES ($1,$2,$3,$4) RETURNING id, created_at",
		},
		{
			table:  "nothing",
			fields: []string{},
			want:   "INSERT INTO nothing () VALUES () RETURNING id, created_at",
		},
		{
			table:  "one",
			fields: []string{"one_field"},
			want:   "INSERT INTO one (one_field) VALUES ($1) RETURNING id, created_at",
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, BuildSQLInsert(tt.table, tt.fields))
	}
}

func TestBuildSQLUpdateByID(t *testing.T) {
	tableTest := []struct {
		table  string
		fields []string
		want   string
	}{
		{
			table:  "cashboxes",
			fields: []string{"responsable", "country", "user_id", "account"},
			want:   "UPDATE cashboxes SET responsable = $1, country = $2, user_id = $3, account = $4, updated_at = now() WHERE id = $5",
		},
		{
			table:  "nothing",
			fields: []string{},
			want:   "",
		},
		{
			table:  "one",
			fields: []string{"one_field"},
			want:   "UPDATE one SET one_field = $1, updated_at = now() WHERE id = $2",
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, BuildSQLUpdateByID(tt.table, tt.fields))
	}
}

func TestBuildSQLSelect(t *testing.T) {
	tableTest := []struct {
		table  string
		fields []string
		want   string
	}{
		{
			table:  "cashboxes",
			fields: []string{"responsable", "country", "user_id", "account"},
			want:   "SELECT id, responsable, country, user_id, account, created_at, updated_at FROM cashboxes",
		},
		{
			table:  "nothing",
			fields: []string{},
			want:   "",
		},
		{
			table:  "one",
			fields: []string{"one_field"},
			want:   "SELECT id, one_field, created_at, updated_at FROM one",
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, BuildSQLSelect(tt.table, tt.fields))
	}
}

func TestBuildSQLSelectFields(t *testing.T) {
	tableTest := []struct {
		table  string
		fields []string
		want   string
	}{
		{
			table:  "cashboxes",
			fields: []string{"responsable", "country", "user_id", "account"},
			want:   "SELECT responsable, country, user_id, account FROM cashboxes",
		},
		{
			table:  "nothing",
			fields: []string{},
			want:   "",
		},
		{
			table:  "one",
			fields: []string{"one_field"},
			want:   "SELECT one_field FROM one",
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, BuildSQLSelectFields(tt.table, tt.fields))
	}
}

func TestBuildSQLWhere(t *testing.T) {
	fakeDate := time.Date(2021, 4, 28, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	tableTest := []struct {
		name      string
		fields    model.Fields
		wantQuery string
		wantArgs  []interface{}
	}{
		{
			name: "where with ILIKE",
			fields: model.Fields{
				{Name: "id", Value: []uint{1, 2, 3}, Operator: model.In},
			},
			wantQuery: "WHERE id IN (1,2,3)",
			wantArgs:  []interface{}{},
		},
		{
			name: "where with all operators",
			fields: model.Fields{
				{Name: "name", Value: "Alejandro"},
				{Name: "age", Value: 30, ChainingKey: model.Or},
				{Name: "course", Value: "Go"},
				{Name: "id", Value: []uint{1, 4, 9}, Operator: model.In},
				{Name: "DESCRIPTION", Value: "%golang%", Operator: model.Ilike},
				{Name: "certificates", Value: 3, Operator: model.GreaterThan},
				{Name: "is_active", Value: true},
			},
			wantQuery: "WHERE name = $1 AND age = $2 OR course = $3 AND id IN (1,4,9) AND description ILIKE $4 AND certificates > $5 AND is_active = $6",
			wantArgs:  []interface{}{"Alejandro", 30, "Go", "%golang%", 3, true},
		},
		{
			name: "where with operators and string ILIKE",
			fields: model.Fields{
				{Name: "country", Value: "COLOMBIA"},
				{Name: "currency_id", Value: 3, ChainingKey: model.Or},
				{Name: "enable", Value: true},
				{Name: "code", Value: []string{"COL", "COP"}, Operator: model.In},
			},
			wantQuery: "WHERE country = $1 AND currency_id = $2 OR enable = $3 AND code IN ('COL','COP')",
			wantArgs:  []interface{}{"COLOMBIA", 3, true},
		},
		{
			name: "where with operators and NOT NULL",
			fields: model.Fields{
				{Name: "country", Value: "COLOMBIA"},
				{Name: "currency_id", Value: 3, ChainingKey: model.Or},
				{Name: "begins_at", Value: "fake", Operator: model.IsNull},
				{Name: "enable", Value: true},
				{Name: "code", Value: []string{"COL", "COP"}, Operator: model.In},
			},
			wantQuery: "WHERE country = $1 AND currency_id = $2 OR begins_at IS NULL AND enable = $3 AND code IN ('COL','COP')",
			wantArgs:  []interface{}{"COLOMBIA", 3, true},
		},
		{
			name: "where with aliased",
			fields: model.Fields{
				{Source: "contracts", Name: "employer_id", Value: 777},
				{Source: "contracts", Name: "pay_frequency_id", Value: 2, ChainingKey: model.Or},
				{Source: "contracts", Name: "is_active", Value: true},
				{Source: "contract_statuses", Name: "description", Value: "ACTIVE", Operator: model.Ilike},
			},
			wantQuery: "WHERE contracts.employer_id = $1 AND contracts.pay_frequency_id = $2 OR contracts.is_active = $3 AND contract_statuses.description ILIKE $4",
			wantArgs:  []interface{}{777, 2, true, "ACTIVE"},
		},
		{
			name: "where with aliased where some fields have missing source",
			fields: model.Fields{
				{Name: "employer_id", Value: 19},
				{Name: "pay_frequency_id", Value: 1, ChainingKey: model.Or},
				{Name: "is_active", Value: false},
				{Source: "contract_statuses", Name: "description", Value: "CREATED", Operator: model.Ilike},
			},
			wantQuery: "WHERE employer_id = $1 AND pay_frequency_id = $2 OR is_active = $3 AND contract_statuses.description ILIKE $4",
			wantArgs:  []interface{}{19, 1, false, "CREATED"},
		},
		{
			name: "where with group conditions",
			fields: model.Fields{
				{Name: "employer_id", Value: 1},
				{Name: "pay_frequency_id", Value: 2},
				{GroupOpen: true, Name: "is_active", Value: true, ChainingKey: model.Or},
				{GroupClose: true, Name: "is_staff", Value: false},
				{Source: "contract_statuses", Name: "description", Value: "ACTIVE", Operator: model.Ilike},
			},
			wantQuery: "WHERE employer_id = $1 AND pay_frequency_id = $2 AND (is_active = $3 OR is_staff = $4) AND contract_statuses.description ILIKE $5",
			wantArgs:  []interface{}{1, 2, true, false, "ACTIVE"},
		},
		{
			name: "where with group conditions and with missing GroupClose key",
			fields: model.Fields{
				{Name: "employer_id", Value: 1},
				{Name: "pay_frequency_id", Value: 2},
				{GroupOpen: true, Name: "is_active", Value: true, ChainingKey: model.Or},
				{Source: "contract_statuses", Name: "description", Value: "ACTIVE", Operator: model.Ilike},
			},
			wantQuery: "WHERE employer_id = $1 AND pay_frequency_id = $2 AND (is_active = $3 OR contract_statuses.description ILIKE $4)",
			wantArgs:  []interface{}{1, 2, true, "ACTIVE"},
		},
		{
			name: "where with group conditions and aliases - complex",
			fields: model.Fields{
				{Source: "c", Name: "employer_id", Value: 1},
				{Source: "c", Name: "termination_date", Operator: model.IsNotNull},
				{Source: "c", Name: "pay_frequency_id", Value: 2},
				{GroupOpen: true, Source: "cs", Name: "description", Operator: model.Ilike, Value: "ACTIVE", ChainingKey: model.Or},
				{GroupOpen: true, Source: "cs", Name: "description", Operator: model.Ilike, Value: "2021-04-28"},
				{GroupClose: true, Source: "c", Name: "hire_date", Operator: model.LessThanOrEqualTo, Value: fakeDate},
			},
			wantQuery: "WHERE c.employer_id = $1 AND c.termination_date IS NOT NULL AND c.pay_frequency_id = $2 AND (cs.description ILIKE $3 OR (cs.description ILIKE $4 AND c.hire_date <= $5))",
			wantArgs:  []interface{}{1, 2, "ACTIVE", "2021-04-28", "2021-04-28"},
		},
	}

	for _, tt := range tableTest {
		gotQuery, gotArgs := BuildSQLWhere(tt.fields)
		assert.Equal(t, tt.wantQuery, gotQuery)
		assert.Equal(t, tt.wantArgs, gotArgs)
	}
}

func TestColumnsAliased(t *testing.T) {
	tableTest := []struct {
		aliased string
		fields  []string
		want    string
	}{
		{
			aliased: "b",
			fields:  []string{"title", "slug", "content", "poster"},
			want:    "b.id, b.title, b.slug, b.content, b.poster, b.created_at, b.updated_at",
		},
		{
			aliased: "nothing",
			fields:  []string{},
			want:    "",
		},
		{
			aliased: "one",
			fields:  []string{"one_field"},
			want:    "one.id, one.one_field, one.created_at, one.updated_at",
		},
	}

	for _, tt := range tableTest {
		assert.Equal(t, tt.want, ColumnsAliased(tt.fields, tt.aliased))
	}
}

func TestBuildSQLOrderBy(t *testing.T) {
	tests := []struct {
		name  string
		sorts model.SortFields
		want  string
	}{
		{
			name:  "Without sort order specification",
			sorts: model.SortFields{{Name: "id"}, {Name: "begins_at"}},
			want:  "ORDER BY id ASC, begins_at ASC",
		},
		{
			name:  "With sort order specification",
			sorts: model.SortFields{{Name: "id", Order: model.Desc}, {Name: "begins_at", Order: model.Asc}},
			want:  "ORDER BY id DESC, begins_at ASC",
		},
		{
			name:  "One field sort",
			sorts: model.SortFields{{Name: "id"}},
			want:  "ORDER BY id ASC",
		},
		{
			name:  "Without field sorts",
			sorts: model.SortFields{},
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildSQLOrderBy(tt.sorts)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_buildIN(t *testing.T) {
	tableTest := []struct {
		field     model.Field
		wantQuery string
	}{
		{
			field: model.Field{
				Name: "id", Value: []uint{1, 2, 3}, Operator: model.In,
			},
			wantQuery: "id IN (1,2,3)",
		},
		{
			field: model.Field{
				Name: "employee_id", Value: []int{5, 6, 7}, Operator: model.In,
			},
			wantQuery: "employee_id IN (5,6,7)",
		},
		{
			field: model.Field{
				Name: "marital_status", Value: []string{"SINGLE"}, Operator: model.In,
			},
			wantQuery: "marital_status IN ('SINGLE')",
		},
		{
			field: model.Field{
				Name: "employee_id", Value: "fake", Operator: model.In,
			},
			wantQuery: "employee_id = 0",
		},
		{
			field: model.Field{
				Name: "contract_id", Value: []uint{}, Operator: model.In,
			},
			wantQuery: "contract_id = 0",
		},
	}

	for _, tt := range tableTest {
		gotQuery := BuildIN(tt.field)
		assert.Equal(t, tt.wantQuery, gotQuery)
	}
}
