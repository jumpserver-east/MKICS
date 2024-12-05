package logic

import (
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MaxkbLogic struct {
	conf   model.MaxkbConf
	mutex  sync.Mutex
	loaded bool
}

type IMaxkbLogic interface {
	GetConfig(c *gin.Context) (dto.MaxkbConf, error)
	UpdateConfig(c *gin.Context, info dto.MaxkbConf) error
	ChatMessage(chatID string, message string) (string, error)
	GetChatID() (string, error)
	getChatIdByKH(khid string) (string, error)
}

func NewIMaxkbLogic() IMaxkbLogic {
	return &MaxkbLogic{}
}

func (u *MaxkbLogic) loadConfig() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	if u.loaded {
		return nil
	}
	conf, err := maxkbRepo.GetConfig()
	if err != nil {
		global.ZAPLOG.Error("failed to get Maxkb config", zap.Error(err))
		return err
	}
	u.conf = conf
	u.loaded = true
	return nil
}

func (u *MaxkbLogic) UpdateConfig(c *gin.Context, info dto.MaxkbConf) error {
	if err := u.loadConfig(); err != nil {
		return err
	}
	if info.BaseURL != "" {
		u.conf.BaseURL = info.BaseURL
	}
	if info.ApiKey != "" {
		u.conf.ApiKey = info.ApiKey
	}
	if err := maxkbRepo.UpdateConfig(u.conf); err != nil {
		global.ZAPLOG.Error("failed to update Maxkb config", zap.Error(err))
		return err
	}
	return nil
}

func (u *MaxkbLogic) GetConfig(c *gin.Context) (dto.MaxkbConf, error) {
	if err := u.loadConfig(); err != nil {
		return dto.MaxkbConf{}, err
	}
	return dto.MaxkbConf{
		BaseURL: u.conf.BaseURL,
		ApiKey:  u.conf.ApiKey,
	}, nil
}

func (u *MaxkbLogic) ChatMessage(chatID string, message string) (string, error) {
	if err := u.loadConfig(); err != nil {
		return "", err
	}
	headers := map[string]string{
		"AUTHORIZATION": u.conf.ApiKey,
		"accept":        "application/json",
		"Content-Type":  "application/json",
	}
	url := fmt.Sprintf("%s/application/chat_message/%s", u.conf.BaseURL, chatID)
	data := map[string]interface{}{
		"message": message,
		"re_chat": false,
		"stream":  true,
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	resp, err := makeRequest("POST", url, headers, bytes.NewBuffer(dataBytes))
	if err != nil {
		return "", err
	}
	return parseResponse(resp)
}

func (u *MaxkbLogic) GetAppID() (string, error) {
	if err := u.loadConfig(); err != nil {
		return "", err
	}
	headers := map[string]string{
		"AUTHORIZATION": u.conf.ApiKey,
		"accept":        "application/json",
		"Content-Type":  "application/json",
	}
	url := u.conf.BaseURL + "/application/profile"
	resp, err := makeRequest("GET", url, headers, nil)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}
	return result["data"].(map[string]interface{})["id"].(string), nil
}

func (u *MaxkbLogic) GetChatID() (string, error) {
	if err := u.loadConfig(); err != nil {
		return "", err
	}
	headers := map[string]string{
		"AUTHORIZATION": u.conf.ApiKey,
		"accept":        "application/json",
		"Content-Type":  "application/json",
	}
	appID, err := u.GetAppID()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/application/%s/chat/open", u.conf.BaseURL, appID)
	resp, err := makeRequest("GET", url, headers, nil)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", err
	}
	return result["data"].(string), nil
}

func (m *MaxkbLogic) getChatIdByKH(khid string) (string, error) {
	khinfo, err := kHRepo.Get(kHRepo.WithKHID(khid))
	if err != nil {
		global.ZAPLOG.Error("数据库没有找到该客户信息", zap.Error(err))
		if err := kHRepo.Create(model.KH{
			KHID: khid,
		}); err != nil {
			return "", err
		}
		global.ZAPLOG.Info("录入客户信息:", zap.String("KHID", khid))
	}
	if khinfo.ChatID == "" {
		global.ZAPLOG.Info("该客户没有和机器人的聊天ID", zap.String("KHID", khid))
		newchatid, err := m.GetChatID()
		if err != nil {
			return "", err
		}
		kh := model.KH{
			KHID:   khid,
			ChatID: newchatid,
		}
		if err := kHRepo.UpdatebyKHID(kh); err != nil {
			return "", err
		}
		global.ZAPLOG.Info("生成聊天id更新客户信息:", zap.String("newchatid", newchatid), zap.String("khid", khid))
		return newchatid, nil
	}
	return khinfo.ChatID, nil
}

func parseResponse(resp []byte) (string, error) {
	contentAll := ""
	jsonStrList := strings.Split(string(resp), "\n")
	for _, jsonStr := range jsonStrList {
		if strings.TrimSpace(jsonStr) != "" {
			splitData := strings.SplitN(jsonStr, ": ", 2)
			if len(splitData) > 1 {
				var jsonData map[string]interface{}
				if err := json.Unmarshal([]byte(splitData[1]), &jsonData); err != nil {
					return "", err
				}
				content, ok := jsonData["content"].(string)
				if !ok {
					return "", fmt.Errorf("invalid content format")
				}
				if content != "" {
					contentAll += content
				}
			}
		}
	}
	contentAll = cleanContent(contentAll)
	return contentAll, nil
}

func cleanContent(content string) string {
	re := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	content = re.ReplaceAllString(content, "$2")

	content = regexp.MustCompile(`\#\s*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\*\*|\_\_`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\*|\_`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`\~\~`).ReplaceAllString(content, "")
	content = regexp.MustCompile("`{1,3}").ReplaceAllString(content, "")

	content = strings.ReplaceAll(content, "\n\n", "\n")
	reNewline := regexp.MustCompile(`\n+`)
	content = reNewline.ReplaceAllString(content, "\n")

	return content
}

func makeRequest(method, url string, headers map[string]string, body *bytes.Buffer) ([]byte, error) {
	client := &http.Client{}
	var req *http.Request
	var err error
	if body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, body)
	}
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}
