package collect

import (
	"bytes"
	"eebot/bot/service/analysis300/db"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var client = http.DefaultClient

var SearchMatchesUrl = "https://300report.jumpw.com/api/battle/searchMatchs?type=h5"

var SearchMatchInfoUrl = "https://300report.jumpw.com/api/battle/searchMatchinfo?type=h5"

var SearchRoleIDFromNameUrl = "https://300report.jumpw.com/api/battle/searchNormal?type=h5"

var contentType = "application/x-www-form-urlencoded; charset=UTF-8"

// SearchMatches 通过url获取战绩列表
func SearchMatches(RoleID uint64, MatchType int, searchIndex int) (matches []db.Match, err error) {
	values := url.Values{}
	values.Set("RoleID", strconv.FormatUint(RoleID, 10))
	values.Set("MatchType", fmt.Sprintf("%d", MatchType))
	values.Set("searchIndex", strconv.Itoa(searchIndex))
	retry := 0
A:
	resp, err := http.DefaultClient.Post(SearchMatchesUrl, contentType, bytes.NewBuffer([]byte(values.Encode())))
	if err != nil {
		return
	}

	body, _ := io.ReadAll(resp.Body)
	rs := map[string]interface{}{}
	err = json.Unmarshal(body, &rs)
	if err != nil {
		return
	}

	if success, ok := rs["success"].(bool); !ok || !success {
		if retry < 10 {
			retry++
			time.Sleep(5 * time.Second * time.Duration(retry))
			goto A
		} else {
			return nil, fmt.Errorf("error interface result, %+v", rs)
		}
	}

	b, err := json.Marshal(rs["data"].(map[string]interface{})["Matchs"].(map[string]interface{})["Matchs"])
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &matches)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 && searchIndex == 1 {
		return nil, errors.New("最近无战绩")
	}

	return
}

// SearchMatchInfo 通过url获取战绩详情
func SearchMatchInfo(matchID string) (match db.Match, err error) {
	values := url.Values{}
	values.Set("mtid", matchID)
	retry := 0
A:
	resp, err := http.DefaultClient.Post(SearchMatchInfoUrl, contentType, bytes.NewBuffer([]byte(values.Encode())))
	if err != nil {
		return
	}

	body, _ := io.ReadAll(resp.Body)
	rs := map[string]interface{}{}
	err = json.Unmarshal(body, &rs)
	if err != nil {
		return
	}

	if success, ok := rs["success"].(bool); !ok || !success {
		if retry < 10 {
			retry++
			time.Sleep(5 * time.Second * time.Duration(retry))
			goto A
		} else {
			return match, fmt.Errorf("error interface result, %+v", rs)
		}
	}

	b, _ := json.Marshal(rs["data"])
	err = json.Unmarshal(b, &match)
	return
}

// SearchRoleID 通过url获取id
func SearchRoleID(name string) (RoleID uint64, err error) {
	defer func() {
		if err != nil && err.Error() == "error interface result" {
			err = fmt.Errorf("不存在 %s 角色", name)
		}
	}()
	if strings.HasPrefix(name, MaskPrefix) {
		id, err := strconv.ParseInt(strings.TrimPrefix(name, MaskPrefix), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("无法将%s转换为int类型", strings.TrimPrefix(name, MaskPrefix))
		}
		_, err = SearchMatches(uint64(id), 1, 1)
		if err != nil {
			return 0, fmt.Errorf("查询id:%s对应的角色失败：%s", name, err.Error())
		}
		return uint64(id), nil
	}
	values := url.Values{}
	values.Set("AccountID", "0")
	values.Set("Guid", "0")
	values.Set("RoleName", name)
	v := values.Encode()

	resp, err := client.Post(SearchRoleIDFromNameUrl, contentType, bytes.NewBuffer([]byte(v)))
	if err != nil {
		return
	}

	tmp, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	rs := map[string]interface{}{}
	err = json.Unmarshal(tmp, &rs)
	if err != nil {
		return
	}

	if success, ok := rs["success"].(bool); !ok || !success {
		return 0, errors.New("error interface result")
	}

	if data, ok := rs["data"].(map[string]interface{}); ok {
		tmp := data["RoleID"].(float64)
		RoleID = uint64(tmp)
	} else {
		err = errors.New("data error")
	}
	return
}

var MaskPrefix = "id:"

// SearchName 通过url获取名称
func SearchName(RoleID uint64) (Name string) {
	defer func() {
		if Name == "*******" {
			Name = fmt.Sprintf("%s%d", MaskPrefix, RoleID)
		}
	}()
	name, err := db.RDB.Get(Ctx, fmt.Sprintf("%s_%d", PlayerIDToNameKey, RoleID)).Result()
	if err == nil {
		return name
	}

	matches, err := SearchMatches(RoleID, 1, 1)
	if err != nil {
		return
	}

	db.RDB.Set(Ctx, fmt.Sprintf("%s_%d", PlayerIDToNameKey, RoleID), matches[0].Players[0].Name, Expiration)
	return matches[0].Players[0].Name
}
