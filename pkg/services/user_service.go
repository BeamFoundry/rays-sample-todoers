package services

import (
	"context"
	"errors"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	"gorm.io/gorm"

	. "github.com/BeamFoundry/rays-orm"
	. "github.com/BeamFoundry/rays-web-security/pkg/core"

	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/repositories"
)

type UserService interface {
	UserDetailsService
	AddUser(name, username, password string) (*User, error)
	GetFullUser(id uint) Optional[*User]
	UpdateName(name string, user *User) (*User, error)
	UpdatePassword(password string, user *User) (*User, error)
}

type UserServiceImpl struct {
	Service
	transactionManager *GormTxManager  `@:"Inject"`
	userRepository     UserRepository  `@:"Inject"`
	encoder            PasswordEncoder `@:"Inject"`
	logger             Logging         `@:"Inject"`
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
	} else {
		return OptionalOf(&user)
	}
}

func (this *UserServiceImpl) LoadUserByUsername(username string) UserDetails {
	if user, ok := this.userRepository.FindByUsername(username).Get(); ok {
		return user
	}
	return nil
}

func (this *UserServiceImpl) AddUser(name, username, password string) (*User, error) {
	if pass, err := this.encoder.Encode([]byte(password)); err != nil {
		return nil, err
	} else {
		tx := this.transactionManager.Start(context.Background())
		defer this.transactionManager.End()
		if _, ok := this.userRepository.FindByUsername(username).Get(); ok {
			tx.Rollback()
			return nil, errors.New("User already exists")
		}
		user := &User{
			Name:     name,
			Username: username,
			Password: string(pass),
		}
		if u, err := this.userRepository.Save(user); err != nil {
			tx.Rollback()
			return nil, err
		} else {
			tx.Commit()
			return u, err
		}
	}
}

func (this *UserServiceImpl) UpdateName(name string, user *User) (*User, error) {
	tx := this.transactionManager.Start(context.Background())
	if err := tx.Model(user).Update("name", name).Error; err != nil {
		tx.Rollback()
		this.transactionManager.End()
		return user, err
	} else {
		tx.Commit()
		this.transactionManager.End()
		user.Name = name
		return user, nil
	}
}

func (this *UserServiceImpl) UpdatePassword(password string, user *User) (*User, error) {
	tx := this.transactionManager.Start(context.Background())
	defer this.transactionManager.End()

	if pass, err := this.encoder.Encode([]byte(password)); err != nil {
		return nil, err
	} else {
		user.Password = string(pass)
		if u, err := this.userRepository.Save(user); err != nil {
			tx.Rollback()
			return user, err
		} else {
			tx.Commit()
			return u, nil
		}
	}
}
