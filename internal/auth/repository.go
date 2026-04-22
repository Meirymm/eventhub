package auth

import (
    "database/sql"
    "eventhub/internal/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
    query := `
        INSERT INTO users (email, password_hash, first_name, last_name, role)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at`
    return r.db.QueryRow(query,
        user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Role,
    ).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, email, password_hash, first_name, last_name, role, created_at
              FROM users WHERE email = $1`
    err := r.db.QueryRow(query, email).Scan(
        &user.ID, &user.Email, &user.PasswordHash,
        &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, email, first_name, last_name, role, created_at
              FROM users WHERE id = $1`
    err := r.db.QueryRow(query, id).Scan(
        &user.ID, &user.Email,
        &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return user, nil
}
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, email, first_name, last_name, role, created_at FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Role, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
