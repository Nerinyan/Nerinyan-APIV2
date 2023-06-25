package config

import (
	"github.com/goccy/go-json"
	"github.com/pterm/pterm"
	"io/ioutil"
	"os"
)

type config struct {
	Port         string   `json:"port"`
	TargetDir    string   `json:"targetDir"`
	SlaveServers []string `json:"slave"`
	Debug        bool     `json:"debug"`
	Discord      struct {
		Webhook struct {
			Startup string `json:"startup"`
			Info    string `json:"info"`
			Error   string `json:"error"`
		} `json:"webhook"`
	} `json:"discord"`
	Sql struct {
		Url string `json:"url"`
	} `json:"sql"`
	Osu struct {
		Username string `json:"username"`
		Passwd   string `json:"passwd"`
		Token    struct {
			TokenType    string `json:"token_type"`
			ExpiresIn    int64  `json:"expires_in"`
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"token"`
		BeatmapUpdate struct {
			UpdatedAsc struct {
				CursorString string `json:"cursor_string"`
			} `json:"updated_asc"`
			UpdatedDesc struct {
				CursorString string `json:"cursor_string"`
			} `json:"updated_desc"`
			GraveyardAsc struct {
				CursorString string `json:"cursor_string"`
			} `json:"graveyard_asc"`
		} `json:"beatmapUpdate"`
	} `json:"osu"`
	Log struct {
		Crawler bool `json:"crawler"`
	} `json:"log"`
}

var Config config

func LoadConfig() {
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		out, err := os.Create("./config.json")
		if err != nil {
			pterm.Error.Println("Can't create ./config.json")
			panic(err)
		}
		defer out.Close()
		body, err := json.MarshalIndent(Config, "", "    ")
		if err != nil {
			pterm.Error.Println("Error. Marshal json")
			panic(err)
		}
		// Write the body to file
		if _, err = out.Write(body); err != nil {
			pterm.Error.Println("Can't Write ./config.json")
			panic(err)
		}
	}

	err = json.Unmarshal(b, &Config)
	if err != nil {
		pterm.Error.Println("fail to parse config.json")
		return
	}
	pterm.Success.Println("Success load config.json")

}
func (v *config) Save() {
	file, _ := json.MarshalIndent(v, "", "  ")
	_ = ioutil.WriteFile("config.json", file, 0755)
}
