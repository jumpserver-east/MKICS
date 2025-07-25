package validator

import (
	"MKICS/backend/constant"
	"MKICS/backend/global"
	"reflect"
	"regexp"
	"strconv"
	"sync"
	"time"

	val "github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Once     sync.Once
	Validate *val.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{}
}

func (v *CustomValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.LazyInit()
		if err := v.Validate.Struct(obj); err != nil {
			global.ZAPLOG.Error(err.Error())
			return err
		}
	}
	return nil
}

func (v *CustomValidator) Engine() interface{} {
	v.LazyInit()
	return v.Validate
}

func (v *CustomValidator) LazyInit() {
	v.Once.Do(func() {
		v.Validate = val.New()
		v.Validate.SetTagName("validate")

		v.Validate.RegisterValidation("int64", ValidateInt64Filed)
		v.Validate.RegisterValidation("phone", ValidatePhoneFiled)
		v.Validate.RegisterValidation("corp_id", ValidateCorpIDFiled)
		v.Validate.RegisterValidation("word", ValidateWordFiled)
		v.Validate.RegisterValidation("ext_id", ValidateExtIDFiled)
		v.Validate.RegisterValidation("time", ValidateTimeFiled)
		v.Validate.RegisterValidation("date", ValidateDateFiled)
		v.Validate.RegisterValidation("boolean", ValidateBooleanFiled)
		v.Validate.RegisterValidation("weekday", ValidateWeekdayFiled)
		v.Validate.RegisterValidation("weekFormat", ValidateWeekFormat)
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func ValidateInt64Filed(fl val.FieldLevel) bool {
	_, err := strconv.ParseInt(fl.Field().String(), 10, 64)
	return err == nil
}

func ValidatePhoneFiled(fl val.FieldLevel) bool {
	ok, _ := regexp.MatchString(`^1[3-9][0-9]{9}$`, fl.Field().String())
	return ok
}

func ValidateWordFiled(fl val.FieldLevel) bool {
	ok, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, fl.Field().String())
	return ok
}

const (
	ExtIDLen  = 32
	CorpIDLen = 18
)

func ValidateExtIDFiled(fl val.FieldLevel) bool {
	matchString, _ := regexp.MatchString(`^[a-zA-Z0-9\-_]+$`, fl.Field().String())
	return matchString && (len(fl.Field().String()) == ExtIDLen)
}

func ValidateCorpIDFiled(fl val.FieldLevel) bool {
	matchString, _ := regexp.MatchString(`^[a-zA-Z0-9\-_]+$`, fl.Field().String())
	return matchString && (len(fl.Field().String()) == CorpIDLen)
}

func ValidateTimeFiled(fl val.FieldLevel) bool {
	_, err := time.Parse(constant.TimeLayout, fl.Field().String())
	return err == nil
}

func ValidateDateFiled(fl val.FieldLevel) bool {
	_, err := time.Parse(constant.DateLayout, fl.Field().String())
	return err == nil
}

func ValidateBooleanFiled(fl val.FieldLevel) bool {
	val := fl.Field().Int()
	return val == 1 || val == 2
}

func ValidateWeekdayFiled(fl val.FieldLevel) bool {
	day := fl.Field().String()
	switch day {
	case "周一", "周二", "周三", "周四", "周五", "周六", "周日":
		return true
	default:
		return false
	}
}

func ValidateWeekFormat(fl val.FieldLevel) bool {
	week := fl.Field().String()
	if len(week) != 7 {
		return false
	}
	for _, char := range week {
		if char != '0' && char != '1' {
			return false
		}
	}
	return true
}
