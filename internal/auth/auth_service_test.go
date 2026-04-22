package auth

import (
	"fmt"
	"testing"

	"eventhub/internal/models"
)

// ══════════════════════════════════════════════════
//  MOCK — фейковая база данных для тестов
//  Не нужен реальный PostgreSQL!
// ══════════════════════════════════════════════════

type mockUserRepo struct {
	users  map[string]*models.User
	nextID int
}

func newMockRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*models.User)}
}

func (m *mockUserRepo) CreateUser(u *models.User) error {
	if _, exists := m.users[u.Email]; exists {
		return fmt.Errorf("email already exists")
	}
	m.nextID++
	u.ID = m.nextID
	copied := *u
	m.users[u.Email] = &copied
	return nil
}

func (m *mockUserRepo) GetUserByEmail(email string) (*models.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (m *mockUserRepo) GetUserByID(id int) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserRepo) GetAllUsers() ([]models.User, error) {
	var list []models.User
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, nil
}

func (m *mockUserRepo) UpdateUser(id int, req *models.UpdateUserRequest) error { return nil }
func (m *mockUserRepo) UpdateUserRole(id int, role string) error               { return nil }

// ══════════════════════════════════════════════════
//  TABLE-DRIVEN ТЕСТ: Register
//  Запуск: go test ./internal/auth/... -v -run TestRegister
// ══════════════════════════════════════════════════

func TestRegister(t *testing.T) {
	tests := []struct {
		name    string                 // название теста
		req     models.RegisterRequest // входные данные
		wantErr bool                   // ожидаем ошибку?
	}{
		{
			name:    "успешная регистрация",
			req:     models.RegisterRequest{Email: "alice@test.com", Password: "123456", FirstName: "Alice", LastName: "Smith"},
			wantErr: false,
		},
		{
			name:    "дублирующий email",
			req:     models.RegisterRequest{Email: "alice@test.com", Password: "654321", FirstName: "Bob", LastName: "Jones"},
			wantErr: true, // alice уже зарегистрирована выше
		},
		{
			name:    "новый пользователь",
			req:     models.RegisterRequest{Email: "bob@test.com", Password: "password1", FirstName: "Bob", LastName: "Jones"},
			wantErr: false,
		},
	}

	repo := newMockRepo()
	svc  := NewAuthService(repo, "test-secret")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.Register(&tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("FAIL: ожидали ошибку, но её нет")
				}
				return
			}
			if err != nil {
				t.Errorf("FAIL: не ожидали ошибку: %v", err)
				return
			}
			if resp.Token == "" {
				t.Error("FAIL: токен пустой")
			}
			if resp.User.Role != models.RoleUser {
				t.Errorf("FAIL: ожидали роль 'user', получили '%s'", resp.User.Role)
			}
		})
	}
}

// ══════════════════════════════════════════════════
//  TABLE-DRIVEN ТЕСТ: Login
//  Запуск: go test ./internal/auth/... -v -run TestLogin
// ══════════════════════════════════════════════════

func TestLogin(t *testing.T) {
	repo := newMockRepo()
	svc  := NewAuthService(repo, "test-secret")

	// Регистрируем пользователя перед тестами
	svc.Register(&models.RegisterRequest{
		Email: "user@test.com", Password: "correct123",
		FirstName: "Test", LastName: "User",
	})

	tests := []struct {
		name    string
		email   string
		pass    string
		wantErr bool
	}{
		{
			name:    "правильный логин",
			email:   "user@test.com",
			pass:    "correct123",
			wantErr: false,
		},
		{
			name:    "неверный пароль",
			email:   "user@test.com",
			pass:    "wrongpass",
			wantErr: true,
		},
		{
			name:    "несуществующий email",
			email:   "nobody@test.com",
			pass:    "123456",
			wantErr: true,
		},
		{
			name:    "пустой пароль",
			email:   "user@test.com",
			pass:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.Login(&models.LoginRequest{
				Email:    tt.email,
				Password: tt.pass,
			})

			if tt.wantErr {
				if err == nil {
					t.Errorf("FAIL: ожидали ошибку — её нет")
				}
				return
			}
			if err != nil {
				t.Errorf("FAIL: не ожидали ошибку: %v", err)
				return
			}
			if resp.Token == "" {
				t.Error("FAIL: токен пустой после логина")
			}
		})
	}
}

// ══════════════════════════════════════════════════
//  TABLE-DRIVEN ТЕСТ: GetMe
//  Запуск: go test ./internal/auth/... -v -run TestGetMe
// ══════════════════════════════════════════════════

func TestGetMe(t *testing.T) {
	repo := newMockRepo()
	svc  := NewAuthService(repo, "test-secret")

	resp, _ := svc.Register(&models.RegisterRequest{
		Email: "me@test.com", Password: "123456",
		FirstName: "Мейрім", LastName: "Султан",
	})

	tests := []struct {
		name    string
		userID  int
		wantErr bool
	}{
		{
			name:    "существующий пользователь",
			userID:  resp.User.ID,
			wantErr: false,
		},
		{
			name:    "несуществующий ID",
			userID:  999,
			wantErr: true,
		},
		{
			name:    "нулевой ID",
			userID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := svc.GetMe(tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Error("FAIL: ожидали ошибку — её нет")
				}
				return
			}
			if err != nil {
				t.Errorf("FAIL: не ожидали ошибку: %v", err)
				return
			}
			if user.Email != "me@test.com" {
				t.Errorf("FAIL: неверный email: %s", user.Email)
			}
		})
	}
}