package errors

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"io/fs"
	"path/filepath"
	"regexp"
)

var (
	languages map[language.Tag]*i18n.Localizer
	i10nPath  = "./assets/locales/errors"
)

func init() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	if err := filepath.Walk(i10nPath, func(path string, info fs.FileInfo, err error) error {
		r, _ := regexp.Compile("[a-z]{2,3}.toml")
		if r.MatchString(info.Name()) {
			if _, e := language.Parse(info.Name()); e != nil {
				return e
			}
			_, err = bundle.LoadMessageFile(path)
			if err != nil {
				log.WithError(err).Errorf("failed to add file path: %v", info.Name())
			}
		}
		return nil
	}); err != nil {
		log.WithError(err).Error("invalid language tag")
	}

	languages := make(map[language.Tag]*i18n.Localizer)
	for _, tag := range bundle.LanguageTags() {
		languages[tag] = i18n.NewLocalizer(bundle, tag.String())
	}
}

func translate(messageId string) map[language.Tag]string {
	translated := make(map[language.Tag]string)
	for tag, l := range languages {
		str, err := l.Localize(&i18n.LocalizeConfig{
			MessageID: messageId,
		})
		if err != nil {
			translated[tag] = messageId
		} else {
			translated[tag] = str
		}
	}
	return translated
}
