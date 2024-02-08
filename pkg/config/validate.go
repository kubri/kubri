package config

import (
	"io/fs"
	"reflect"
	"strings"
	"unicode"

	locales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
	"golang.org/x/mod/semver"

	"github.com/kubri/kubri/pkg/version"
)

type Error struct{ Errors []string }

func (e *Error) Error() string {
	return "invalid config:\n  " + strings.Join(e.Errors, "\n  ")
}

func Validate(c *config) error {
	v := validator.New(validator.WithRequiredStructEnabled(), validator.WithOmitAnonymousName())
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name, _, _ := strings.Cut(fld.Tag.Get("yaml"), ",")
		if name == "-" {
			return ""
		}
		return name
	})
	_ = v.RegisterValidation("dirname", isDirname)
	_ = v.RegisterValidation("version", isVersion)
	_ = v.RegisterValidation("version_constraint", isConstraint)

	uni := ut.New(locales.New())
	trans, _ := uni.GetTranslator("en")
	_ = registerTranslations(v, trans)

	err := v.Struct(c)
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	errs := &Error{Errors: make([]string, len(validationErrors))}
	for i, ve := range validationErrors {
		_, ns, _ := strings.Cut(ve.Namespace(), ".")
		errs.Errors[i] = ns + strings.TrimPrefix(ve.Translate(trans), ve.Field())
	}

	return errs
}

func registerTranslations(v *validator.Validate, trans ut.Translator) error {
	err := translations.RegisterDefaultTranslations(v, trans)
	if err != nil {
		return err
	}

	translations := []struct {
		name    string
		message string
	}{
		{
			name:    "dir",
			message: "{0} must be a valid path to a directory",
		},
		{
			name:    "dirname",
			message: "{0} must be a valid folder name",
		},
		{
			name:    "version_constraint",
			message: "{0} must be a valid version constraint",
		},
		{
			name:    "version",
			message: "{0} must be a valid semver version",
		},
		{
			name:    "http_url",
			message: "{0} must be a valid URL",
		},
		{
			name:    "fqdn|http_url",
			message: "{0} must be a valid URL or FQDN",
		},
	}

	for _, t := range translations {
		err = v.RegisterTranslation(t.name, trans, func(ut ut.Translator) error {
			return ut.Add(t.name, t.message, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(fe.Tag(), fe.Field())
			return t
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// isDirname validates if a dirname is valid in an URL, as well as common targets.
//
//   - RFC3986 allows the following in a URL path (i.e. pchar): A-Z a-z 0-9 - . _ ~ ! $ & ' ( ) * + , ; = : @
//   - GCS forbids # [ ] * ? : " < > | (See https://cloud.google.com/storage/docs/objects#naming)
//   - S3 forbids \ { ^ } % ` ] " > [ ~ < # | (See https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-keys.html#object-key-guidelines-avoid-characters)
//   - This leaves us with letters, numbers and - . _ ! $ ' ( ) + , ; = @
func isDirname(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return fs.ValidPath(s) && !strings.ContainsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) && !strings.ContainsRune("/-._!$&'()+,;=@ ", r)
	})
}

func isVersion(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	return semver.IsValid(v)
}

func isConstraint(fl validator.FieldLevel) bool {
	_, err := version.NewConstraint(fl.Field().String())
	return err == nil
}
