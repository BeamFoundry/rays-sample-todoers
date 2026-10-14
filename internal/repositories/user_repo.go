package repositories

import (
	. "github.com/BeamFoundry/rays-orm/pkg/core"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/core"
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
	. "github.com/BeamFoundry/rays/pkg/lang"
)

type UserRepositoryImpl struct {
	OrmRepository[User, uint]
	findByUsername                      func(string) Optional[*User]
	findAllUsersByUsernameLikeOrderByID func(string, ...OrmQueryFunc) []User
}

func (this *UserRepositoryImpl) FindByUsername(username string) Optional[*User] {
	return this.findByUsername(username)
}

// FindAllUsersByUsernameLike implements [UserRepository].
func (this *UserRepositoryImpl) FindAllUsersByUsername(query string, queries ...OrmQueryFunc) []User {
	return this.findAllUsersByUsernameLikeOrderByID(query, queries...)
}

var _ UserRepository = &UserRepositoryImpl{}
