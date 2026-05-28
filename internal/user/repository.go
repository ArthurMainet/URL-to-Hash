package user

import "golang/packages/db"

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (u *UserRepository) Create(user *User) (*User, error) {
	result := u.Database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (u *UserRepository) FindByEmail(emailID string) (*User, error) {
	var user User
	result := u.Database.DB.First(&user, "email = ?", emailID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
