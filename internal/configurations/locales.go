package configurations

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/i18n"
	. "github.com/BeamFoundry/rays/pkg/i18n/locales"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/web"
)

var defaultLocale = Locales.EN

type LocalizationProperties struct {
	ConfigurationProperties `Name:"app"`
	messages                map[string]LocalizedMessages
}

type LocalizationConfig struct {
	Configuration
	localizations *LocalizationProperties `@:"Inject"`
}

func (this *LocalizationConfig) MessageSource() (MessageSource, Error) {
	messageSource := DefaultMessageSource{}.New(defaultLocale)

	bundle := messageSource.DefaultMessageBundle()

	if err := bundle.LoadValidationTranslation(ValidationTranslations.EN); err != nil {
		return nil, NewErrorFrom(err)
	}
	bundle.LoadLocalizedMessages(this.localizations.messages["en"])

	bundle = messageSource.NewMessageBundle(Locales.FR)

	if err := bundle.LoadValidationTranslation(ValidationTranslations.FR); err != nil {
		return nil, NewErrorFrom(err)
	}
	bundle.LoadLocalizedMessages(this.localizations.messages["fr"])
	return messageSource, nil
}

func (this *LocalizationConfig) LocaleContextHolder(
	_ struct {
		ConditionalOnMissingBeam
	},
) LocaleContextHolder {
	return ThreadLocaleContextHolder{}.New(defaultLocale)
}

func (this *LocalizationConfig) MessageSourceValidator(
	that struct {
		localeContextHolder LocaleContextHolder `@:"Inject"`
		messageSource       MessageSource       `@:"Inject"`
	},
) MessageSourceValidator {
	return DefaultMessageSourceValidator{}.New(that.localeContextHolder, that.messageSource)
}

func (this *LocalizationConfig) LocaleFuncMap(
	that struct {
		localeContextHolder LocaleContextHolder `@:"Inject"`
		messageSource       MessageSource       `@:"Inject"`
	},
) *LocaleFuncMap {
	return LocaleFuncMap{}.New(that.localeContextHolder, that.messageSource)
}
