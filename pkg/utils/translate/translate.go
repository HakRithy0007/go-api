package utils

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Translate(MessageID string, param *string, c *fiber.Ctx) string {
	var translate string
	if param != nil {
		translate = fiberi18n.MustLocalize(c, &i18n.LocalizeConfig{
			MessageID: MessageID,
			TemplateData: map[string]interface{}{
				"name": param,
			},
		})
	} else {
		translate = fiberi18n.MustLocalize(c, &i18n.LocalizeConfig{
			MessageID: MessageID,
		})
	}
	return translate

}
