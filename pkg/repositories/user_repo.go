package repositories

import (
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
	. "github.com/BeamFoundry/rays-orm"
	. "github.com/BeamFoundry/rays/pkg/lang"
)

type UserRepository interface {
	OrmRepositoryInterface[User, uint]
	FindByUsername(string) Optional[*User]
}

type UserRepositoryImpl struct {
	OrmRepositoryComposite[User, uint]
	findByUsername func(string) Optional[*User]
}

func (this *UserRepositoryImpl) FindByUsername(username string) Optional[*User] {
	return this.findByUsername(username)
}

var _ UserRepository = &UserRepositoryImpl{}
