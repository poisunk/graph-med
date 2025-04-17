package repository

import (
	"errors"
	"graph-med/internal/base/data"
	"graph-med/internal/model"
)

type UserRepository struct {
	data *data.Data
}

func NewUserRepository(data *data.Data) *UserRepository {
	return &UserRepository{
		data: data,
	}
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*model.ChatUser, error) {
	user := &model.ChatUser{}
	exist, err := r.data.DB.Where("username = ?", username).Get(user)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(email string) (*model.ChatUser, error) {
	user := &model.ChatUser{}
	exist, err := r.data.DB.Where("email = ?", email).Get(user)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// CreateUser 创建用户
func (r *UserRepository) CreateUser(user *model.ChatUser) error {
	_, err := r.data.DB.Insert(user)
	return err
}

// UpdateUser 更新用户
func (r *UserRepository) UpdateUser(user *model.ChatUser) error {
	_, err := r.data.DB.ID(user.ID).Update(user)
	return err
}

// DeleteUser 删除用户
func (r *UserRepository) DeleteUser(id int) error {
	_, err := r.data.DB.ID(id).Delete(&model.ChatUser{})
	return err
}

// FindUserByID 根据ID查找用户
func (r *UserRepository) FindUserByID(id int) (*model.ChatUser, error) {
	user := &model.ChatUser{}
	exists, err := r.data.DB.ID(id).Get(user)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
