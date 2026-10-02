package config

import (
	. "github.com/BeamFoundry/rays-web-security/pkg/config"
	. "github.com/BeamFoundry/rays-web-security/pkg/core"
	. "github.com/BeamFoundry/rays-web-security/pkg/rememberme"
	. "github.com/BeamFoundry/rays/pkg/core"
)

type SecurityConfig struct {
	Configuration
}

func (this *SecurityConfig) RememberMe(
	that struct {
		userDetailsServicer UserDetailsService `@:"Inject"`
	},
) RememberMeService {
	return TokenBasedRememberMeService{}.New("rememberme", that.userDetailsServicer)
}

func (this *SecurityConfig) HttpSecurity(
	that struct {
		rememberme RememberMeService `@:"Inject"`
	},
) *HttpSecurity {
	return (&HttpSecurityBuilder{}).
		AuthorizedHTTPRequests().
		Paths("/profile").Authenticated().
		Paths("/todolists**").HasAuthority("ROLE_USER").
		Paths("/admin/**").HasAuthority("ROLE_ADMIN").
		Paths("/static/**").PermitAll().
		Paths("/.well-known/**").PermitAll().
		Paths("/*").PermitAll().
		And().
		AnyRequest().DenyAll().
		FormLogin().
		ViewName("views/login").
		And().
		RememberMe(that.rememberme).
		Build()
}
