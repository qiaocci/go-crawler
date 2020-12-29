package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/url"
	"time"
)

func main() {
	douyinUrl := "https://v.douyin.com/JSdPMst/"
	userInfo := getUserInfo(douyinUrl)
	fmt.Printf("%v\n", userInfo)
	lastSignature := getLastSignature()
	fmt.Printf("老签名：%v, \n新签名：%v\n", lastSignature, userInfo.Signature)
	if lastSignature != userInfo.Signature {
		fmt.Println("签名改变啦")
		WxNotify("闪闪签名改变啦", "")
	} else {
		fmt.Println("签名没改")
	}

	Save(userInfo)
}

func WxNotify(text, desp string) {
	url := "https://sc.ftqq.com/SCU9800T1aa9ee59f94cfe6bcde0b23b4b91135d5959fab2590de.send"

	params := map[string]string{"text": text, "desp": desp}
	download(url, params)
}

func getLastSignature() string {
	db, err := sql.Open("mysql", "root:123@/douyin")
	if err != nil {
		fmt.Println("db connect error")
	}
	row := db.QueryRow("select id, signature from douyin order by id desc limit 1")
	user := new(ParsedUserInfo)
	if err := row.Scan(&user.ID, &user.Signature); err != nil {
		fmt.Println("scan error")
	}
	fmt.Printf("用户id:%v, 签名：%v", user.ID, user.Signature)
	return user.Signature
}

func download(url string, param map[string]string) *http.Response {
	client := &http.Client{CheckRedirect: CheckRedirect}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1")
	req.Header.Set("accept-encoding", "deflate")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("upgrade-insecure-requests", "1")

	q := req.URL.Query()
	if param != nil {
		for k, v := range param {
			q.Add(k, v)
		}
	}
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return nil
	}
	//defer resp.Body.Close() // 关闭连接

	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	//fmt.Println(resp.Location())
	//fmt.Println(resp.Header)

	return resp
}

func CheckRedirect(req *http.Request, via []*http.Request) error {
	fmt.Println(req.RequestURI)
	return http.ErrUseLastResponse
}

func test() {
	myUrl := "http://www.baidu.com/?name=tom"
	u, e := url.Parse(myUrl)
	if e != nil {
		fmt.Println(e)
	}
	q := u.Query()
	name := q["name"][0]

	fmt.Println(name)
}

type UserInfoResponse struct {
	UserInfo   UserInfo `json:"user_info"`
	StatusCode int      `json:"status_code"`
}

type UserInfo struct {
	FollowerCount   int    `json:"follower_count"`
	FavoritingCount int    `json:"favoriting_count"`
	TotalFavorited  string `json:"total_favorited"`
	UniqueID        string `json:"unique_id"`
	ShortID         string `json:"short_id"`
	AvatarLarge     AvatarLarge
	UID             string `json:"uid"`
	NickName        string `json:"nickname"`
	Signature       string `json:"signature"`
	FollowingCount  int    `json:"following_count"`
}

type AvatarLarge struct {
	URI        string   `json:"uri"`
	URLList    []string `json:"url_list"`
	AwemeCount int      `json:"aweme_count"`
}

func Save(userInfo ParsedUserInfo) {
	timezone := "'Asia/Shanghai'"

	db, err := sql.Open("mysql", "root:123@/douyin?charset=utf8mb4&parseTime=true&loc=Local&time_zone="+url.QueryEscape(timezone))
	if err != nil {
		fmt.Println("db connect error")
	}
	l, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(l)
	fmt.Println(now)
	result, err := db.Exec("insert into douyin (unique_id, nickname, signature, created) values (?,?,?,?)",
		userInfo.UniqueID, userInfo.Nickname, userInfo.Signature, now)
	if err != nil {
		fmt.Println("db execute error")
	}
	fmt.Println(result)

	//rows, err := db.Query("select * from douyin")
}

func Parse() {
	body := `{
    "status_code": 1,
    "user_info": {
        "follower_count": 53000,
        "favoriting_count": 3668,
        "total_favorited": "1697000",
        "custom_verify": "",
        "unique_id": "Call03",
        "short_id": "143112824",
        "avatar_larger": {
            "uri": "319d40009dc5cf0ec02a4",
            "url_list": [
                "https://p3-dy-ipv6.byteimg.com/aweme/1080x1080/319d40009dc5cf0ec02a4.jpeg?from=4010531038",
                "https://p9-dy.byteimg.com/aweme/1080x1080/319d40009dc5cf0ec02a4.jpeg?from=4010531038",
                "https://p1-dy-ipv6.byteimg.com/aweme/1080x1080/319d40009dc5cf0ec02a4.jpeg?from=4010531038"
            ]
        },
        "aweme_count": 95,
        "original_musician": {
            "music_count": 0,
            "music_used_count": 0
        },
        "is_gov_media_vip": false,
        "avatar_thumb": {
            "uri": "319d40009dc5cf0ec02a4",
            "url_list": [
                "https://p26-dy.byteimg.com/aweme/100x100/319d40009dc5cf0ec02a4.jpeg?from=4010531038",
                "https://p9-dy.byteimg.com/aweme/100x100/319d40009dc5cf0ec02a4.jpeg?from=4010531038",
                "https://p3-dy-ipv6.byteimg.com/aweme/100x100/319d40009dc5cf0ec02a4.jpeg?from=4010531038"
            ]
        },
        "region": "CN",
        "secret": 0,
        "uid": "61490865008",
        "nickname": "奶茶闪闪✨吃什么",
        "type_label": null,
        "verification_type": 1,
        "followers_detail": null,
        "platform_sync_info": null,
        "geofencing": null,
        "policy_version": null,
        "signature": "✨听起来很好吃\n🌈在厦门上班的治愈系美食博主\n💜吃喝日常，蟹蟹你的关注",
        "avatar_medium": {
            "uri": "319d40009dc5cf0ec02a4",
            "url_list": [
                "https://p3-dy-ipv6.byteimg.com/aweme/720x720/319d40009dc5cf0ec02a4.jpeg?from=4010531038",
                "https://p9-dy.byteimg.com/aweme/720x720/319d40009dc5cf0ec02a4.jpeg?from=4010531038",
                "https://p26-dy.byteimg.com/aweme/720x720/319d40009dc5cf0ec02a4.jpeg?from=4010531038"
            ]
        },
        "following_count": 166
    },
    "extra": {
        "now": 1609168074000,
        "logid": "20201228230754010198062098361DD3FC"
    }
}`

	var r UserInfoResponse
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Printf("err was %v", err)
	}
	fmt.Printf("状态码=%v, 用户信息：%v, 粉丝数：%v", r.StatusCode, r.UserInfo, r.UserInfo.FollowerCount)
}

type ParsedUserInfo struct {
	ID        string `json:"id"`
	UniqueID  string `json:"unique_id"`
	Nickname  string `json:"nickname"`
	Signature string `json:"signature"`
}

func getUserInfo(theUrl string) ParsedUserInfo {
	realUrl := getRealAddress(theUrl)
	// [https://www.iesdouyin.com/share/user/61490865008?sec_uid=MS4wLjABAAAAI7jMi3c1toaNF49sQcWWt7fKtZAQ6eR-mp7I0PcQXXg&u_code=198jab20e&timestamp=1599819622&utm_source=copy&utm_campaign=client_share&utm_medium=android&share_app_name=douyin]
	u, err := url.Parse(realUrl)
	if err != nil {
		fmt.Println("parse error")
	}
	hostName := u.Host
	fmt.Printf("hostName=%v\n", hostName)
	q := u.Query()
	secUid := q["sec_uid"][0]
	fmt.Printf("secUid=%v\n", secUid)
	userInfoUrl := fmt.Sprintf("https://%v/web/api/v2/user/info/", hostName)
	userInfoParams := make(map[string]string)
	userInfoParams["sec_uid"] = secUid
	fmt.Printf("userInfoUrl=%v\n", userInfoUrl)
	resp := download(userInfoUrl, userInfoParams)

	var r UserInfoResponse
	json.NewDecoder(resp.Body).Decode(&r)
	fmt.Printf("%v, %v, %v\n", r, r.StatusCode, r.UserInfo)

	return ParsedUserInfo{
		UniqueID:  r.UserInfo.UniqueID,
		Nickname:  r.UserInfo.NickName,
		Signature: r.UserInfo.Signature,
	}
}

func getRealAddress(url string) string {
	resp := download(url, nil)
	fmt.Println(resp.Header)
	if resp.StatusCode == 302 {
		return resp.Header["Location"][0]
	}

	return ""
}
