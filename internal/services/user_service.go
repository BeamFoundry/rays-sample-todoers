package services

import (
	"context"
	"errors"
	"fmt"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	"gorm.io/gorm"

	. "github.com/BeamFoundry/rays-orm/pkg/core"
	. "github.com/BeamFoundry/rays-web-security/pkg/core"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/core"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
)

func paginate(page, pageSize int) OrmQueryFunc {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type UserServiceImpl struct {
	Service
	transactional  GormTransactional `@:"Inject"`
	userRepository UserRepository    `@:"Inject"`
	encoder        PasswordEncoder   `@:"Inject"`
	logger         Logging           `@:"Inject"`
}

var _ UserService = &UserServiceImpl{}

func (this *UserServiceImpl) GetFullUser(id uint) Optional[*User] {
	user := User{}
	res := this.userRepository.Query(func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id).Preload("Todos").Preload("Todos.Items")
	}).First(&user)
	this.logger.Debug("GetFullUser", "user", user)
	if res.RowsAffected == 0 {
		return OptionalEmpty[*User]()
	}
	return OptionalOf(&user)
}

func (this *UserServiceImpl) LoadUserByUsername(username string) UserDetails {
	if user, ok := this.userRepository.FindByUsername(username).Get(); ok {
		return user
	}
	return nil
}

func (this *UserServiceImpl) AddUser(name, username, password string) (*User, error) {
	pass, err := this.encoder.Encode([]byte(password))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	tx := this.transactional.Start(context.Background())
	defer this.transactional.End()
	if _, ok := this.userRepository.FindByUsername(username).Get(); ok {
		tx.Rollback()
		return nil, errors.New("entity 'User' already exists")
	}
	user := &User{
		Name:     name,
		Username: username,
		Password: string(pass),
	}
	if this.CountAll() == 0 {
		user.Admin = true
	}
	u, err := this.userRepository.Save(user)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%w", err)
	}
	tx.Commit()
	return u, nil
}

func (this *UserServiceImpl) UpdateName(name string, user *User) (*User, error) {
	tx := this.transactional.Start(context.Background())
	if err := tx.Model(user).Update("name", name).Error; err != nil {
		tx.Rollback()
		this.transactional.End()
		return user, err
	}
	tx.Commit()
	this.transactional.End()
	user.Name = name
	return user, nil
}

func (this *UserServiceImpl) UpdatePassword(password string, user *User) (*User, error) {
	tx := this.transactional.Start(context.Background())
	defer this.transactional.End()

	pass, err := this.encoder.Encode([]byte(password))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	user.Password = string(pass)
	u, err := this.userRepository.Save(user)
	if err != nil {
		tx.Rollback()
		return user, fmt.Errorf("%w", err)
	}
	tx.Commit()
	return u, nil
}

func (this *UserServiceImpl) FindAllUsersByUsername(query string, page, pageSize int) []User {
	return this.userRepository.FindAllUsersByUsername("%"+query+"%", paginate(page, pageSize))
}

func (this *UserServiceImpl) CountAll() int64 {
	var count int64
	this.userRepository.Query().Count(&count)
	return count
}
