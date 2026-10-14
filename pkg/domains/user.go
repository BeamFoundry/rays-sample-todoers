package domains

import (
	"encoding/gob"
	"time"

	// "gorm.io/gorm"
	// "gorm.io/gorm/clause"

	. "github.com/BeamFoundry/rays-web-security/pkg/core"
)

//nolint:govet // ignoring
type User struct {
	ID        uint
	Name      string
	Username  string
	Password  string
	Admin     bool
	Todos     []Todo
	CreatedAt time.Time
	UpdatedAt time.Time
}

/*
func (this *User) AfterDelete(tx *gorm.DB) (err error) {
	tx.Clauses(clause.Returning{}).Where("user_id = ?", this.ID).Delete(&Todo{})
	return
}
*/

// GetAuthorities implements core.UserDetails.
func (this *User) GetAuthorities() []GrantedAuthority {
	if this.Admin {
		return []GrantedAuthority{Authority("ROLE_ADMIN"), Authority("ROLE_USER")}
	}
	return []GrantedAuthority{Authority("ROLE_USER")}
}

func (this *User) EraseCredentials() {
	this.Password = ""
}

// GetPassword implements core.UserDetails.
func (this *User) GetPassword() string {
	return this.Password
}

// GetUsername implements core.UserDetails.
func (this *User) GetUsername() string {
	return this.Username
}

// IsCredentialExpired implements core.UserDetails.
func (this *User) IsCredentialExpired() bool {
	return false
}

// IsDisabled implements core.UserDetails.
func (this *User) IsDisabled() bool {
	return false
}

// IsExpired implements core.UserDetails.
func (this *User) IsExpired() bool {
	return false
}

// IsLocked implements core.UserDetails.
func (this *User) IsLocked() bool {
	return false
}

var _ UserDetails = &User{}

func init() {
	gob.Register(&User{})
}
