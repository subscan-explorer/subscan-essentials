package validator

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/itering/subscan/util/address"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var validate = validator.New()

func init() {
	_ = registerCustomValidator(validate)
}

func Validate(data interface{}, model interface{}) (err error) {
	var b []byte
	switch v := data.(type) {
	case []byte:
		b = v
	case io.ReadCloser:
		b, _ = io.ReadAll(v)
	default:
		b, _ = json.Marshal(data)
	}
	if err = json.Unmarshal(b, model); err != nil {
		return err
	}
	return validate.Struct(model)
}

type translator struct {
	tag           string
	registerFn    validator.RegisterTranslationsFunc
	translationFn validator.TranslationFunc
}

var translatorList = []translator{
	{
		tag: "addr",
		registerFn: func(ut ut.Translator) error {
			return ut.Add("addr", "The {0} an invalid address", true)
		},
		translationFn: func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("addr", fe.Field())
			return t
		},
	},
	{
		tag: "time_range",
		registerFn: func(ut ut.Translator) error {
			return ut.Add("time_range", "The {0} is invalid time range", true)
		},
		translationFn: func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("time_range", fe.Field())
			return t
		},
	},
	{
		tag: "date_range",
		registerFn: func(ut ut.Translator) error {
			return ut.Add("date_range", "The {0} is invalid data range", true)
		},
		translationFn: func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("date_range", fe.Field())
			return t
		},
	},
	{
		tag: "num_range",
		registerFn: func(u ut.Translator) error {
			return u.Add("num_range", "The {0} is invalid number range", true)
		},
		translationFn: func(u ut.Translator, fe validator.FieldError) string {
			t, _ := u.T("num_range", fe.Field())
			return t
		},
	},
}

var trans ut.Translator

func registerCustomValidator(v *validator.Validate) error {
	_ = v.RegisterValidation("addr", addr)
	_ = v.RegisterValidation("time_range", timeRange)
	_ = v.RegisterValidation("date_range", dateRange)
	_ = v.RegisterValidation("num_range", numberRange)
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)
		if len(name) > 0 && name[0] != "-" {
			return name[0]
		}
		return ""
	})
	eng := english.New()
	uni := ut.New(eng, eng)
	trans, _ = uni.GetTranslator("en")
	if err := en.RegisterDefaultTranslations(v, trans); err != nil {
		return err
	}
	for _, ts := range translatorList {
		_ = v.RegisterTranslation(ts.tag, trans, ts.registerFn, ts.translationFn)
	}
	return nil
}

func TranslationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			return e.Translate(trans) // return first error
		}
	}
	return err.Error()
}

func RegisterCustomValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := registerCustomValidator(v); err != nil {
			panic(err)
		}
	}
}

var addr validator.Func = func(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		panic("unsupported struct type")
	}
	str := fl.Field().String()
	if len(str) == 0 {
		return false
	}
	return len(address.Decode(str)) != 0
}

var numberRange validator.Func = func(fl validator.FieldLevel) bool {
	var val int64
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = fl.Field().Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = int64(fl.Field().Uint())
	case reflect.String:
		var err error
		if val, err = strconv.ParseInt(fl.Field().String(), 10, 64); err != nil {
			return false
		}
	default:
		panic("unsupported struct type")
	}
	numberRange := strings.Split(fl.Param(), "-")
	if len(numberRange) != 2 {
		panic(fmt.Sprintf("Bad number_range value %s", numberRange))
	}
	var (
		minNumber, maxNumber int64
		err                  error
	)
	if minNumber, err = strconv.ParseInt(numberRange[0], 10, 64); err != nil {
		panic(fmt.Sprintf("Bad number_range Type. err: %s", err.Error()))
	}
	if maxNumber, err = strconv.ParseInt(numberRange[1], 10, 64); err != nil {
		panic(fmt.Sprintf("Bad number_range Type. err: %s", err.Error()))
	}
	if val < minNumber || val > maxNumber {
		return false
	}
	return true
}

var dateRange validator.Func = func(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		panic("unsupported struct type")
	}
	timeLimit := strings.Split(fl.Param(), " ")
	if len(timeLimit) != 2 {
		panic(fmt.Sprintf("Bad date_range value %s", timeLimit))
	}
	var (
		startTime, endTime, valTime time.Time
		timeFormat                  = "2006-01-02"
		err                         error
	)
	if startTime, err = time.Parse(timeFormat, timeLimit[0]); err != nil {
		panic(fmt.Sprintf("Bad time_range Type. err: %s", err.Error()))
	}
	if endTime, err = time.Parse(timeFormat, timeLimit[1]); err != nil {
		panic(fmt.Sprintf("Bad time_range Type. err: %s", err.Error()))
	}
	if valTime, err = time.Parse(timeFormat, fl.Field().String()); err != nil {
		return false
	}
	if endTime.Before(startTime) {
		panic("Bad time_range value. start time should be before end time")
	}
	if valTime.Before(startTime) || valTime.After(endTime) {
		return false
	}
	return true
}

var timeRange validator.Func = func(fl validator.FieldLevel) bool {
	if tm, ok := fl.Field().Interface().(string); ok {
		timeLimit := strings.Split(fl.Param(), "+")
		timeFormat, startTimeStr, EndTimeStr := "2006-01-02", "", ""
		switch len(timeLimit) {
		case 2:
			startTimeStr, EndTimeStr = timeLimit[0], timeLimit[1]
		case 3:
			timeFormat, startTimeStr, EndTimeStr = timeLimit[0], timeLimit[1], timeLimit[2]
		default:
			panic(fmt.Sprintf("Bad time_range value %s", timeLimit))
		}
		var startTime, endTime, valTime time.Time
		switch timeFormat {
		case "unix":
			if startUnix, err := strconv.ParseInt(startTimeStr, 10, 64); err != nil {
				panic(fmt.Sprintf("Bad time_range Type. err: %s", err.Error()))
			} else {
				startTime = time.Unix(startUnix, 0)
			}
			if endUnix, err := strconv.ParseInt(EndTimeStr, 10, 64); err != nil {
				panic(fmt.Sprintf("Bad time_range Type. err: %s", err.Error()))
			} else {
				endTime = time.Unix(endUnix, 0)
			}

			if valUnix, err := strconv.ParseInt(tm, 10, 64); err != nil {
				return false
			} else {
				valTime = time.Unix(valUnix, 0)
			}

		default:
			var err error
			if startTime, err = time.Parse(timeFormat, startTimeStr); err != nil {
				panic(fmt.Sprintf("Bad time_range Type. err: %s", err.Error()))
			}
			if endTime, err = time.Parse(timeFormat, EndTimeStr); err != nil {
				panic(fmt.Sprintf("Bad time_range Type. err: %s", err.Error()))
			}
			if valTime, err = time.Parse(timeFormat, tm); err != nil {
				return false
			}
		}
		if endTime.Before(startTime) {
			panic("Bad time_range value. start time should be before end time")
		}
		if valTime.Before(startTime) || valTime.After(endTime) {
			return false
		}
		return true
	}
	var val int64
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64:
		val = fl.Field().Int()
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		val = int64(fl.Field().Uint())
	default:
		panic("unsupported struct type")
	}
	timeLimit := strings.Split(fl.Param(), "+")
	var startTime, endTime, valTime time.Time
	if len(timeLimit) != 2 {
		panic(fmt.Sprintf("Bad time_range value %s", timeLimit))
	}
	if startUnix, err := strconv.ParseInt(timeLimit[0], 10, 64); err != nil {
		panic(fmt.Sprintf("Bad time_range Type. err: %s", err.Error()))
	} else {
		startTime = time.Unix(startUnix, 0)
	}
	if endUnix, err := strconv.ParseInt(timeLimit[1], 10, 64); err != nil {
		panic(fmt.Sprintf("Bad time_range Type. err: %s", err.Error()))
	} else {
		endTime = time.Unix(endUnix, 0)
	}
	valTime = time.Unix(val, 0)
	if endTime.Before(startTime) {
		panic("Bad time_range value. start time should be before end time")
	}
	if valTime.Before(startTime) || valTime.After(endTime) {
		return false
	}
	return true
}
