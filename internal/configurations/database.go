package configurations

import (
	. "github.com/BeamFoundry/rays-sample-todoers/pkg/domains"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	sloggorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"

	sqlitenocgo "github.com/glebarez/sqlite"
	"gorm.io/driver/sqlite"
)

// This serves as an example for Conditionals
//
// DbMigrator needs one gorm.DD pointer. It will be delivered by one of DbConfig function.
type DbMigrator struct {
	Service
	db *gorm.DB `@:"Inject"`
}

func (this *DbMigrator) PostConstruct() Error {
	return Try(func() {
		RaiseErr(this.db.AutoMigrate(&User{}))
		RaiseErr(this.db.AutoMigrate(&Todo{}))
		RaiseErr(this.db.AutoMigrate(&Item{}))
	})
}

//---

type DBOpener interface {
	Open(string) gorm.Dialector
}

type SQLiteDB struct{}

func (this *SQLiteDB) Open(dns string) gorm.Dialector {
	return sqlite.Open(dns)
}

type SQLiteDBNoCGO struct{}

func (this *SQLiteDBNoCGO) Open(dns string) gorm.Dialector {
	return sqlitenocgo.Open(dns)
}

//---

type DialectorConfig struct {
	Configuration
	logger Logging `@:"Inject"`
}

func (this *DialectorConfig) DBOpener(_ struct {
	ConditionalOnMissingBeam
},
) DBOpener {
	this.logger.Info("DBOpener with CGO")
	return &SQLiteDB{}
}

func (this *DialectorConfig) DBOpenerNoCGO(
	_ struct {
		ConditionalOnNoCGO
	},
) DBOpener {
	this.logger.Info("DBOpener without CGO")
	return &SQLiteDBNoCGO{}
}

// This serves as an example for Conditionals
type DbConfig struct {
	Configuration
	//
	logger Logging  `@:"Inject"`
	db     DBOpener `@:"Inject"`
}

// It should be a private function otherwise, with the "Configuration" annotation,
// Rays container will try to create it.
func (this *DbConfig) config() *gorm.Config {
	return &gorm.Config{
		Logger: sloggorm.New(
			sloggorm.WithHandler(this.logger.Handler()),
		),
	}
}

// This method will be called if no other *gorm.DB is defined (ConditionalOnMissingBeam).
func (this *DbConfig) DB(
	_ struct {
		ConditionalOnMissingBeam
	},
) (*gorm.DB, Error) {
	db, err := gorm.Open(this.db.Open("todolist.db"), this.config())
	if err != nil {
		return nil, NewErrorFrom(err)
	}

	return db, nil
}

// This method will be called if the "profile" named "test" is activated (ConditionalOnProfile), creating
// a Beam of type "*gorm/DB".
//
// See cmd/main.go for activating this profile.
//
// The Beam defined in DB() will fail its condition (ConditionalOnMissingBeam) as it's not the only one and as such won't be defined.
//
// As only one Beam will be defined of type "*gorm.DB", DbMigrator injection will not fail (only one Beam will satisfy the "db" field)
func (this *DbConfig) DBTesting(
	_ struct {
		ConditionalOnProfile `Name:"test"`
	},
) (*gorm.DB, Error) {
	db, err := gorm.Open(this.db.Open("file::memory:?cache=shared"), this.config())
	if err != nil {
		return nil, NewErrorFrom(err)
	}

	return db, nil
}
