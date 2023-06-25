package download

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/db"
	"github.com/labstack/echo/v4"
	"github.com/pterm/pterm"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
)

//go:embed * embed.html
var beatmapSetEmbed embed.FS

type TemplateRenderer struct {
	Templates *template.Template
}

var Renderer = TemplateRenderer{
	Templates: template.Must(template.ParseFS(beatmapSetEmbed, "*.html")),
}

var statusWithIcon = map[int]string{
	4:  "💟 Loved",
	3:  "✅ Qualified",
	2:  "🔥 Approved",
	1:  "⏫ Ranked",
	0:  "❔ Pending",
	-1: "🛠️ WIP",
	-2: "⚰️ Graveyard",
}
var status = map[int]string{
	4:  "Loved",
	3:  "Qualified",
	2:  "Approved",
	1:  "Ranked",
	0:  "Pending",
	-1: "WIP",
	-2: "Graveyard",
}
var modeString = map[int]string{
	3: "mania",
	2: "catch",
	1: "taiko",
	0: "osu!",
}

type setEmbed struct {
	TITLE           string
	CREATOR         string
	RANKED          int
	BPM             float64
	FAVOURITE_COUNT string
}
type mapEmbed struct {
	MODE_INT          int
	VERSION           string
	DIFFICULTY_RATING float64
	TOTAL_LENGTH      int64
	CS                float64
	AR                float64
}

//

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.Templates.ExecuteTemplate(w, name, data)
}

func Embed(next echo.HandlerFunc) echo.HandlerFunc {

	//return
	return func(c echo.Context) error {
		if !strings.Contains(strings.ToLower(c.Request().Header.Get("User-Agent")), "discord") {
			return next(c)
		}
		setId := c.Param("setId")
		var set setEmbed
		var Map []mapEmbed
		db.Gorm.Raw("SELECT TITLE, CREATOR, RANKED, BPM, FAVOURITE_COUNT FROM BEATMAPSET WHERE BEATMAPSET_ID = ?", setId).Scan(&set)
		db.Gorm.Raw("SELECT MODE_INT, VERSION, DIFFICULTY_RATING, TOTAL_LENGTH, CS, AR FROM BEATMAP WHERE DELETED_AT IS NULL AND BEATMAPSET_ID = ? ORDER BY DIFFICULTY_RATING", setId).Scan(&Map)
		pterm.Info.Println(set)
		pterm.Info.Println(Map)
		pterm.Info.Println("================")
		return c.Render(
			http.StatusOK, "embed.html", map[string]interface{}{
				"setId": setId,
				"title": set.TITLE,
				"content": func() (content string) {
					var res string
					//Ranked osu! beatmap by Kyuukai.
					res += status[set.RANKED] + " osu! beatmap by " + set.CREATOR + "\n"
					//               ⏫ Ranked             · 📚               9                Difficulties · 🎵                       190                  · ❤️ 601
					res += statusWithIcon[set.RANKED] + " · 📚 " + strconv.Itoa(len(Map)) + " Difficulties · 🎵 " + fmt.Sprintf("%.0f", set.BPM) + " · ❤ " + set.FAVOURITE_COUNT + "\n"
					res += "\n"
					for _, m := range Map {
						//                       (osu!) Normal - ⭐ 2.07 · ⏳ 2:17 | CS 3.2 · AR 4.5
						res += fmt.Sprintf("(%s) %s - ⭐ %.2f · ⏳ %s | CS %.1f · AR  %.1f \n", modeString[m.MODE_INT], m.VERSION, m.DIFFICULTY_RATING, parseTime(m.TOTAL_LENGTH), m.CS, m.AR)
					}

					return res
				}(),
			},
		)
	}
}

func parseTime(t int64) (ms string) {
	return fmt.Sprintf("%d:%d", t/60, t%60)
}
