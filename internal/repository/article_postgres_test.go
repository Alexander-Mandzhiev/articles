package repository

import (
	"articles/internal/models"
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestArticlePostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewArticlePostgres(db)

	type mockBehavior func(author *models.Article)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		input        *models.Article
		wantErr      bool
	}{
		// 1
		{
			name: "OK",
			input: &models.Article{
				ID:      "92b7f107-e3a6-4053-ae25-a374ed8cc1b7",
				Title:   "Author Ok",
				Text:    "test desc",
				Authors: "1",
			},
			mockBehavior: func(author *models.Article) {
				rows := []string{"id", "name", "text", "authors"}
				mock.ExpectQuery("INSERT INTO articles").
					WithArgs(author.ID, author.Title, author.Text, author.Authors).
					WillReturnRows(sqlmock.NewRows(rows).AddRow("92b7f107-e3a6-4053-ae25-a374ed8cc1b7", "Author Ok", "test desc", "1"))
			},
		},
		{
			name: "Empty Fields",
			input: &models.Article{
				ID:      "",
				Title:   "Author Ok",
				Text:    "test desc",
				Authors: "",
			},
			mockBehavior: func(author *models.Article) {
				rows := sqlmock.NewRows([]string{"id", "name", "text", "authors"}).
					AddRow("Author Ok", "test desc", "", "").
					RowError(1, errors.New("empti id field"))

				mock.ExpectQuery("INSERT INTO authors").
					WithArgs(author.Title, author.Text, author.Authors).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		//3
		{
			name: "Insert Error",
			input: &models.Article{
				ID:      "92b7f107-e3a6-4053-ae25-a374ed8cc1b7",
				Title:   "Author Ok",
				Text:    "test desc",
				Authors: "",
			},
			mockBehavior: func(author *models.Article) {
				mock.ExpectQuery("INSERT INTO authors").
					WithArgs(author.ID, author.Title, author.Text, author.Authors).
					WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)
			got, err := r.Create(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.input, got)
			}

		})
	}
}

func TestAuthorsPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewArticlePostgres(db)

	tests := []struct {
		name    string
		mock    func()
		want    []models.Article
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "text", "authors"}).
					AddRow("1", "title1", "description1", "1").
					AddRow("2", "title2", "description2", "2").
					AddRow("3", "title3", "description3", "3")

				mock.ExpectQuery("SELECT (.+) FROM articles").WillReturnRows(rows)
			},

			want: []models.Article{
				{"1", "title1", "description1", "1"},
				{"2", "title2", "description2", "2"},
				{"3", "title3", "description3", "3"},
			},
		},
		{
			name: "No Records",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "text", "authors"})

				mock.ExpectQuery("SELECT (.+) FROM articles").
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAll()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthorPostgres_GetOne(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewArticlePostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    models.Article
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "text", "authors"}).
					AddRow("1", "title1", "description1", "1")

				mock.ExpectQuery("SELECT (.+) FROM articles WHERE (.+)").
					WillReturnRows(rows)
			},
			input: "1",
			want:  models.Article{"1", "title1", "description1", "1"},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "text", "authors"}).
					AddRow("1", "title1", "description1", "1")

				mock.ExpectQuery("SELECT (.+) FROM articles WHERE (.+)").
					WillReturnRows(rows)
			},
			input:   "1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetOne(tt.input)
			if tt.wantErr {
				assert.Empty(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewArticlePostgres(db)

	tests := []struct {
		name    string
		mock    func()
		input   string
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("DELETE FROM articles WHERE (.+)").
					WithArgs("1").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: "1",
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec("DELETE FROM articles WHERE (.+)").
					WithArgs("1").WillReturnError(sql.ErrNoRows)
			},
			input:   "1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Delete(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewArticlePostgres(db)

	type args struct {
		id      string
		article models.Article
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK_AllFields",
			mock: func() {
				rows := []string{"id", "name", "text", "authors"}

				mock.ExpectQuery("UPDATE articles SET (.+) WHERE (.+)").
					WillReturnRows(sqlmock.NewRows(rows).
						AddRow("92b7f107-e3a6-4053-ae25-a374ed8cc1b7", "Author Ok", "test desc", "1"))
			},
			input: args{
				id: "92b7f107-e3a6-4053-ae25-a374ed8cc1b7",
				article: models.Article{
					Title:   "Author Ok",
					Text:    "test desc",
					Authors: "1",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			_, err := r.Update(tc.input.id, tc.input.article)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
