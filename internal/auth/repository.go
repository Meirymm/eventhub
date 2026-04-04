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

func (r *UserRepository) UpdateUser(id int, req *models.UpdateUserRequest) error {
    query := `UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4`
    _, err := r.db.Exec(query, req.FirstName, req.LastName, req.Email, id)
    return err
}
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
    rows, err := r.db.Query(
        `SELECT id, email, first_name, last_name, role, created_at 
         FROM users ORDER BY id`, )
     if err != nil {
        return nil, err}
     defer rows.Close()
     var users []models.User
      for rows.Next() {
        var u models.User
        rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Role, &u.CreatedAt)
        users = append(users, u)
    }
    return users, nil
}
func (r *UserRepository) UpdateUserRole(userID int, role string) error {
    _, err := r.db.Exec(
        `UPDATE users SET role=$1 WHERE id=$2`,
        role, userID,
    )
    return err
}