package auth

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/log"
	"io"
	"net/http"
	"void-studio.net/fiesta/fapi/config"
)

//type UserRoles struct {
//	Roles []string `json:"roles"`
//}

//type UserToken struct {
//	AccessToken *string `json:"access_token"`
//}
//
//type TokenForm struct {
//	ClientID     string `json:"client_id"`
//	ClientSecret string `json:"client_secret"`
//	GrantType    string `json:"grant_type"`
//	Code         string `json:"code"`
//	RedirectURI  string `json:"redirect_uri"`
//}

func GetAccessRole(accessToken string) []string {
	req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/v10/users/@me/guilds/206367403733549056/member", nil)
	if err != nil {
		log.Errorf("failed to create request: %v", err)
		return nil
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("failed to get response: %v", err)
		return nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read response: %v", err)
		return nil
	}

	var user map[string]interface{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Errorf("failed to unmarshal response: %v", err)
		return nil
	}

	//for srv, accessRole := range config.Config.Discord.AccessRoles {
	//	if srv == server {
	//		return slices.Contains(roles, accessRole)
	//	}
	//}

	return checkAccessRoles(user["roles"].([]string))
}

func GetUserToken(code string) string {
	//form := TokenForm{
	//	ClientID:     config.Config.Discord.ID,
	//	ClientSecret: config.Config.Discord.Secret,
	//	GrantType:    "authorization_code",
	//	Code:         code,
	//	RedirectURI:  config.Config.Discord.AuthURI,
	//}

	form := map[string]string{
		"code":          code,
		"client_id":     config.Config.Discord.ID,
		"client_secret": config.Config.Discord.Secret,
		"grant_type":    "authorization_code",
		"redirect_uri":  config.Config.Discord.AuthURI,
	}

	jsonForm, err := json.Marshal(form)
	if err != nil {
		log.Errorf("failed to marshal form: %v", err)
		return ""
	}

	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/oauth2/token", bytes.NewBuffer(jsonForm))
	if err != nil {
		log.Errorf("failed to create request: %v", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("failed to get response: %v", err)
		return ""
	}

	defer resp.Body.Close()

	//var userToken UserToken
	var userToken map[string]string
	if err = json.NewDecoder(resp.Body).Decode(&userToken); err != nil {
		log.Errorf("failed to decode response: %v", err)
		return ""
	}

	return userToken["access_token"]
}

func checkAccessRoles(roles []string) []string {
	var comparedRoles []string

	for _, accessRole := range config.Config.Discord.AccessRoles {
		for _, userRole := range roles {
			if accessRole == userRole {
				comparedRoles = append(comparedRoles, userRole)
			}
		}
	}

	//for srv, accessRole := range config.Config.Discord.AccessRoles {
	//	if srv == server {
	//		return slices.Contains(roles, accessRole)
	//	}
	//}

	return comparedRoles
}
