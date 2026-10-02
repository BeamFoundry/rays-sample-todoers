# Init

```go
type Stereotype interface {
	_Stereotype()
}

type Component interface {
	Stereotype
	_Component()
}

type Service interface {
	Stereotype
	_Service()
}

type Configuration interface {
	Stereotype
	_Configuration()
}

type ConfigurationProperties interface {
	Stereotype
	_Property() string
}

type Controller interface {
	Stereotype
	_Controller()
}

type Repository interface {
	Stereotype
	_Repository()
}
```
