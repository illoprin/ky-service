package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"ky-id-backend/src/storage"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type UserModel struct {
	Id       int64
	RoleId   int64
	Login    string
	Email    string
	Password string
}

// Returns interface to interact with sqlite database
func New(storage_path string) (*Storage, error) {

	// Open connection
	db, err := sql.Open("sqlite3", storage_path)

	if err != nil {
		return nil, fmt.Errorf("could not open storage")
	}

	// Create tables for users and roles
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS user (
			id INTEGER PRIMARY KEY,
			role_id INTEGER NOT NULL,
			login TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL,
			password TEXT NOT NULL -- bcrypt hash
		);
		CREATE INDEX IF NOT EXISTS idx_login ON user(login);
		CREATE INDEX IF NOT EXISTS idx_email ON user(email);
	`)

	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec()
	_ = res

	if err != nil {
		return nil, fmt.Errorf("could not create table 'user'")
	}

	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS role(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL UNIQUE
		);
	`)

	if err != nil {
		return nil, err
	}

	res, err = stmt.Exec()
	_ = res

	if err != nil {
		return nil, fmt.Errorf("could not create table 'role'")
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddUser(login, email, password string) (int64, error) {
	res, err := s.db.Exec(`
		INSERT INTO user(role_id, login, email, password) VALUES (1, ?, ?, ?);
	`, login, email, password)

	if err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			return 0, storage.ErrSameLoginExist
		}
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Pls, validate 'login' field, to prevent sql-injection
func (s *Storage) GetUserByLogin(login string) (*UserModel, error) {
	var user UserModel

	row := s.db.QueryRow(`SELECT id, role_id, login, email, password FROM user WHERE login = ?`, login)

	err := row.Scan(&user.Id, &user.RoleId, &user.Login, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, storage.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}

	return &user, nil
}

func (s *Storage) GetUserById(id int64) (*UserModel, error) {
	var user UserModel

	row := s.db.QueryRow(`SELECT id, role_id, login, email, password FROM user WHERE id = ?`, id)

	err := row.Scan(&user.Id, &user.RoleId, &user.Login, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, storage.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}

	return &user, nil
}

func (s *Storage) GetUsers() ([]UserModel, error) {
	rows, err := s.db.Query("SELECT id, role_id, login, email, password FROM user")
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	users := make([]UserModel, 0)
	for rows.Next() {
		var user UserModel
		if err := rows.Scan(&user.Id, &user.RoleId, &user.Login, &user.Email, &user.Password); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return users, nil
}

func (s *Storage) DeleteUser(id int64) (int64, error) {
	res, err := s.db.Exec("DELETE FROM user WHERE id = ?", id)

	if err == sqlite3.ErrNotFound {
		return 0, storage.ErrUserNotFound
	} else if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rows, nil
}
