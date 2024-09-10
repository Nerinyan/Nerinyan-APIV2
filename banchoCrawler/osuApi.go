package banchoCrawler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type banchoJWT struct {
	/*
		iss: 토큰 발급자 (issuer)
		sub: 토큰 제목 (subject)
		aud: 토큰 대상자 (audience)
		exp: 토큰의 만료시간 (expiraton), 시간은 NumericDate 형식으로 되어있어야 하며 (예: 1480849147370) 언제나 현재 시간보다 이후로 설정되어있어야합니다.
		nbf: Not Before 를 의미하며, 토큰의 활성 날짜와 비슷한 개념입니다. 여기에도 NumericDate 형식으로 날짜를 지정하며, 이 날짜가 지나기 전까지는 토큰이 처리되지 않습니다.
		iat: 토큰이 발급된 시간 (issued at), 이 값을 사용하여 토큰의 age 가 얼마나 되었는지 판단 할 수 있습니다.
		jti: JWT의 고유 식별자로서, 주로 중복적인 처리를 방지하기 위하여 사용됩니다. 일회용 토큰에 사용하면 유용합니다.
	*/
	Aud    string   `json:"aud"`
	Jti    string   `json:"jti"`
	Iat    float64  `json:"iat"`
	Nbf    float64  `json:"nbf"`
	Exp    float64  `json:"exp"`
	Sub    string   `json:"sub"`
	Scopes []string `json:"scopes"`
}

func parseTokenExpiraton() (t int64) {

	s := strings.Split(config.Config.Osu.Token.AccessToken, ".")
	if len(s) != 3 {
		return
	}

	decodeString, err := base64.RawStdEncoding.DecodeString(s[1])
	if err != nil {
		decodeString, err = base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			decodeString, err = base64.RawURLEncoding.DecodeString(s[1])
			if err != nil {
				decodeString, err = base64.URLEncoding.DecodeString(s[1])
				if err != nil {
					return
				}
			}
		}
	}
	var j banchoJWT
	if err = json.Unmarshal(decodeString, &j); err != nil {
		return
	}
	t = int64(j.Exp) - time.Now().Unix()
	return

}

func tryLogin() (err error) {
	spinner, _ := pterm.DefaultSpinner.Start("Trying Login Bancho...")
	if err = login(true); err != nil {
		spinner.Fail("fail refresh Bancho Token: ", err)
		if err = login(false); err != nil {
			return
		}
		spinner.Success("successful New Login Bancho")
	} else {
		spinner.Success("successful refresh Bancho Token")
	}
	config.Config.Save()
	return
}

func LoadBancho(ch chan struct{}) {
	b := false
	if parseTokenExpiraton() > 120 {
		pterm.Info.Println("Bancho Token Alive.")
		b = true
		ch <- struct{}{}
	}
	for {
		if parseTokenExpiraton() < 3600 || pause { //1시간
			err := tryLogin()
			if err != nil {
				pterm.Error.Println("LoadBancho tryLogin(). :", err)
			} else {
				pause = false
				if !b {
					b = true
					ch <- struct{}{}
				}
			}
		}
		time.Sleep(time.Second * 10)
	}
}

func login(refresh bool) (err error) {

	url := "https://osu.ppy.sh/oauth/token"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("client_id", "5")
	_ = writer.WriteField("client_secret", "FGc9GAtyHzeQDshWP5Ah7dega8hJACAJpQtw6OXk")
	_ = writer.WriteField("scope", "*")

	if refresh {
		_ = writer.WriteField("grant_type", "refresh_token")
		_ = writer.WriteField("refresh_token", config.Config.Osu.Token.RefreshToken)
	} else {
		_ = writer.WriteField("username", config.Config.Osu.Username)
		_ = writer.WriteField("password", config.Config.Osu.Passwd)
		_ = writer.WriteField("grant_type", "password")
	}

	err = writer.Close()
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if refresh {
		req.Header.Set("Authorization", config.Config.Osu.Token.TokenType+" "+config.Config.Osu.Token.AccessToken)
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprint(res.StatusCode, string(body)))
	}
	return json.Unmarshal(body, &config.Config.Osu.Token)
}
