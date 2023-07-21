package repository

import (
	"errors"
	"time"

	"github.com/dimasyudhana/Qoin-Digital-Indonesia/app/middlewares"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/features/user"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/utils/identity"
	"github.com/dimasyudhana/Qoin-Digital-Indonesia/utils/password"
	"gorm.io/gorm"
)

var log = middlewares.Log()

type Query struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.Repository {
	return &Query{
		db: db,
	}
}

// Register implements user.UserData.
func (uq *Query) Register(request user.UserCore) (user.UserCore, error) {
	userId, err := identity.GenerateID()
	if err != nil {
		log.Error("error while creating id for user")
		return user.UserCore{}, errors.New("error while creating id for user")
	}

	hashed, err := password.HashPassword(request.Password)
	if err != nil {
		log.Error("error while hashing password")
		return user.UserCore{}, errors.New("error while hashing password")
	}

	request.UserID = userId
	request.Password = hashed
	req := userEntities(request)
	query := uq.db.Exec(`
		INSERT INTO users (user_id, username, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.UserID, req.Username, req.Email, req.Password, time.Now(), time.Now())

	if errors.Is(query.Error, gorm.ErrRegistered) {
		log.Error("error insert data, duplicate input")
		return user.UserCore{}, errors.New("error insert data, duplicate input")
	}

	rowAffect := query.RowsAffected
	if rowAffect == 0 {
		log.Warn("no row affected")
		return user.UserCore{}, errors.New("no row affected")
	}

	log.Sugar().Infof("new user has been created: %s", req.UserID)
	return userModels(req), nil
}

// Login implements user.UserData.
func (uq *Query) Login(request user.UserCore) (user.UserCore, string, error) {
	result := User{}
	query := uq.db.Raw(`
		SELECT user_id, username, password
		FROM users
		WHERE users.username = ?
		LIMIT 1
	`, request.Username).Scan(&result)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("user record not found")
		return user.UserCore{}, "", errors.New("invalid email and password")
	}

	rowAffect := query.RowsAffected
	if rowAffect == 0 {
		log.Warn("no row affected")
		return user.UserCore{}, "", errors.New("no row affected")
	}

	if !password.MatchPassword(request.Password, result.Password) {
		return user.UserCore{}, "", errors.New("password does not match")
	}

	token, err := middlewares.GenerateToken(result.UserID)
	if err != nil {
		log.Error("error while creating jwt token")
		return user.UserCore{}, "", errors.New("error while creating jwt token")
	}

	log.Sugar().Infof("user has been logged in: %s", result.UserID)
	return userModels(result), token, nil
}

// ProfileUser implements user.Repository.
func (uq *Query) Profile(userId string) (user.UserCore, error) {
	users := User{}
	query := uq.db.Raw("SELECT * FROM users WHERE user_id = ? AND is_deleted IS NULL", userId).Scan(&users)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("user profile record not found")
		return user.UserCore{}, errors.New("user profile record not found")
	}

	return userModels(users), nil
}
