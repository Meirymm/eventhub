package auth

import (
    "database/sql"
    "errors"
    "time"
    "eventhub/internal/models"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    repo      *UserRepository
    jwtSecret string
}

func NewAuthService(repo *UserRepository, jwtSecret string) *AuthService {
    return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &models.User{
        Email:        req.Email,
        PasswordHash: string(hash),
        FirstName:    req.FirstName,
        LastName:     req.LastName,
        Role:         models.RoleUser,
    }

    if err := s.repo.CreateUser(user); err != nil {
        return nil, errors.New("email already exists")
    }

    token, err := s.generateToken(user)
    if err != nil {
        return nil, err
    }

    return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
    user, err := s.repo.GetUserByEmail(req.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("invalid email or password")
        }
        return nil, err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return nil, errors.New("invalid email or password")
    }

    token, err := s.generateToken(user)
    if err != nil {
        return nil, err
    }

    return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) GetMe(userID int) (*models.User, error) {
    return s.repo.GetUserByID(userID)
}

func (s *AuthService) UpdateMe(userID int, req *models.UpdateUserRequest) (*models.User, error) {
    err := s.repo.UpdateUser(userID, req)
    if err != nil {
        return nil, err
    }
    return s.repo.GetUserByID(userID)
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    string(user.Role),
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}
func (s *AuthService) GetAllUsers() ([]models.User, error) {
    return s.repo.GetAllUsers()
}

func (s *AuthService) UpdateUserRole(userID int, role string) error {
    return s.repo.UpdateUserRole(userID, role)
}