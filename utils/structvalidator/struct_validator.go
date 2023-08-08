package structvalidator

import (
	"reflect"
	"strings"

	enLocales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

// custom validation tags (prefix with "_")
const (
	PasswordTag = "_password"
)

// StructValidator wraps https://github.com/go-playground/validator utility to create human-readable validation messages.
type StructValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

// New creates new struct validator. Golang's field names will be used for validation messages.
func New() *StructValidator {
	enTranslator := enLocales.New()
	validate := validator.New()

	universalTranslator := ut.New(enTranslator)
	translator, _ := universalTranslator.GetTranslator(enTranslator.Locale())
	_ = enTranslations.RegisterDefaultTranslations(validate, translator)

	return &StructValidator{
		validator:  validate,
		translator: translator,
	}
}

// NewWithOptions creates new struct validator with options
func NewWithOptions(options ...func(*StructValidator)) *StructValidator {
	structValidator := New()
	for _, o := range options {
		if o != nil {
			o(structValidator)
		}
	}

	return structValidator
}

// FieldTag will be used for validation messages field name.
func WithFieldTag(tag string) func(*StructValidator) {
	if tag == "" {
		return nil
	}

	return func(sv *StructValidator) {
		sv.validator.RegisterTagNameFunc(tagRegister{Tag: tag}.registerTag)
	}
}

// Will add custom translation to the specified tag
// example: tag: "startswith" -> translation: "{0} should start with {1}"
func WithCustomTranslation(tag, translation string) func(*StructValidator) {
	if tag == "" || translation == "" {
		return nil
	}

	return func(sv *StructValidator) {
		ct := customTranslation{
			Tag:         tag,
			Translation: translation,
		}
		ct.RegisterCustomTranslation(sv.validator, sv.translator)
	}
}

// ValidateStruct validates struct, returning translated human-readable error messages. It returns error if input is not struct.
func (j *StructValidator) ValidateStruct(s interface{}) ([]string, error) {
	if err := j.validator.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return j.translateValidationErrors(validationErrors), nil
		}
		return nil, err
	}

	return nil, nil
}

func (j *StructValidator) translateValidationErrors(validationErrors validator.ValidationErrors) []string {
	var messages []string
	for _, validationError := range validationErrors {
		messages = append(messages, validationError.Translate(j.translator))
	}
	return messages
}

type tagRegister struct {
	Tag string
}

func (t tagRegister) registerTag(field reflect.StructField) string {
	name := strings.SplitN(field.Tag.Get(t.Tag), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

type customTranslation struct {
	Tag         string
	Translation string
}

func (ct customTranslation) RegisterCustomTranslation(validate *validator.Validate, translator ut.Translator) {
	validate.RegisterTranslation(
		ct.Tag,
		translator,
		func(ut ut.Translator) error {
			return ut.Add(ct.Tag, ct.Translation, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(ct.Tag, fe.Field(), fe.Param())
			return t
		},
	)
}
