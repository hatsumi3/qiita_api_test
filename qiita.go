// 参考URL
// https://qiita.com/jpshadowapps/items/463b2623209479adcd88
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"gopkg.in/ini.v1"
)

// QiitaConfig is saved token information
type QiitaConfig struct {
	QiitaToken string
}

// QiitaItem is the response data for the qiita api
type QiitaItem struct {
	Created_at string `json:"created_at,omitempty"`
	Title      string `json:"title,omitempty"`
}

var (
	BaseUrl       string = "https://qiita.com/api/v2/items"
	UrlParameters url.Values
	Conf          QiitaConfig
)

func makeUrl() string {
	UrlParameters = url.Values{}
	UrlParameters.Add("query", "tag:Web")
	UrlParameters.Add("per_page", "5")
	return BaseUrl

	// http.Client.GetのときはURL Parameterをこのタイミングでつくる。
	// return BaseUrl + UrlParameters.Encode()

}

func (item *QiitaItem) values() string {
	return item.Created_at + " " + item.Title
}

// RequestForQiita is function
func RequestForQiita() {

	Url := makeUrl()

	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		log.Panicln(err)
	}
	req.URL.RawQuery = UrlParameters.Encode()
	req.Header.Set("Content-Type", "Application/json")
	req.Header.Set("Authorization", "Bearer "+Conf.QiitaToken)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	var qiitaItems []QiitaItem

	// json.NewDecoder(response.Body).Decodeはstreamのときに使うのだ。へけ。
	if err := json.NewDecoder(resp.Body).Decode(&qiitaItems); err != nil {
		log.Fatalln(err)
	}

	// json.Unmarshalはstream以外のときに使うのね、ハム太郎？
	// bodyArray, err := ioutil.ReadAll(resp.Body)
	// if err := json.Unmarshal(bodyArray, &qiitaItems); err != nil {
	// 	log.Fatalln(err)
	// }

	for i, item := range qiitaItems {
		log.Println(i, item.values())
	}

}

// init() にてiniファイルからtokenを読み込む。
func init() {
	conn, err := ini.Load("./qiita.ini")
	if err != nil {
		log.Panic(err)
	}
	Conf = QiitaConfig{
		QiitaToken: conn.Section("web").Key("QiitaToken").String(),
	}
}

func main() {
	log.Println("start")
	RequestForQiita()
	log.Println("end")
}
