package utils

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// InitTrans 初始化翻译器
func InitTrans(locale string) (Trans ut.Translator, err error) {
	//修改gin框架中的validator属性，实现自定制  .(*validator.Validate)类型断言，
	//类型断言用于检查一个接口值是否包含特定类型的值
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器

		uni := ut.New(enT, zhT, enT)
		var ok bool
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return nil, fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, Trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		}
		return Trans, nil
	}
	return Trans, nil
}
