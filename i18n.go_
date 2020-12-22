package fibers;

import (
    "github.com/gofiber/fiber/v2"
    i18ns "github.com/nicksnyder/go-i18n/v2/i18n"
    "golang.org/x/text/language"
    "github.com/pelletier/go-toml"
)


func I18nMiddleware() func(c *fiber.Ctx) error {
    bundle := i18ns.NewBundle(language.English)
    bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
    bundle.LoadMessageFile("locales/active.en.toml")
    bundle.LoadMessageFile("locales/active.de.toml")

    return func(c *fiber.Ctx) error {
        accept          := c.Get("Accept-Language")
        langcookie      := c.Cookies("lang")
        localizer := i18ns.NewLocalizer(bundle, langcookie, "en", accept)
        c.Locals("i18n", localizer)
        return c.Next();
    }
}

func i18n(c *fiber.Ctx) func(id string, args...map[string]interface{}) string {
    localizer := c.Locals("i18n");

    return func(id string, args...map[string]interface{}) string {

        tplargs := map[string]interface{}{};
        if len(args) > 0 {
            tplargs = args[0];
        }

        s, _ := localizer.(*i18ns.Localizer).Localize(&i18ns.LocalizeConfig{
            DefaultMessage: &i18ns.Message{
                ID:     id,
                Other:  id,
            },
            TemplateData: tplargs,
        })
        return s;
    }
}
