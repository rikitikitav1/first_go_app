package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	// "strconv"
)

func credentials_check(parameter string, value string) string {
	result := ""
	regexp_secret := regexp.MustCompile(`^\w{64}$`)
	regexp_uuid := regexp.MustCompile(`^\w{4}(\w{4}-{1}){4}\w{12}$`)

	if regexp_uuid.MatchString(value) {
		result = value
	}

	if regexp_secret.MatchString(value) {
		result = value
	}

	if result == "" {
		fmt.Println("Неверно задано значение", parameter)
		os.Exit(1)
	}

	return result
}

type ApiClient struct {
	uuid   string
	secret string
	client int64
}

func NewApiClient(uuid string, secret string) ApiClient {
	return ApiClient{
		uuid:   credentials_check("UUID", uuid),
		secret: credentials_check("SECRET", secret),
		client: get_client(uuid, secret),
	}
}

func get_client(uuid string, secret string) int64 {
	url := URL + "v1/objects/user"
	values := map[string]map[string]string{"filter": {"uuid": uuid}}
	jsonValue, _ := json.Marshal(values)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	req.Header.Add("X-WallarmAPI-UUID", uuid)
	req.Header.Add("X-WallarmAPI-Secret", secret)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	parse_string := "body.#[uuid==" + uuid + "]\".clientid"
	value := gjson.GetBytes(body, parse_string).Int()
	if value == 0 {
		fmt.Println("Неверно заданы UUID и SECRET")
		os.Exit(1)
	}
	return value

}

func (ac *ApiClient) get_users() {
	url := URL + "v1/objects/user"
	filter := map[string]map[string]int64{"filter": {"clientid": ac.client}}
	parse_string := "body.#[id>0]#.email"
	body := ac.req_builder("POST", url, filter, "json", "json")
	value := gjson.GetBytes(body, parse_string).Array()
	fmt.Println(value)
}

func (ac *ApiClient) req_builder(method string, url string, values map[string]map[string]int64, content string, accept string) []byte {
	jsonValue, _ := json.Marshal(values)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
	req.Header.Add("X-WallarmAPI-UUID", ac.uuid)
	req.Header.Add("X-WallarmAPI-Secret", ac.secret)
	if content == "json" {
		req.Header.Add("Content-Type", "application/json")
	}
	if accept == "json" {
		req.Header.Add("Accept", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

const URL = "https://api.wallarm.com/"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Недостаточно аргументов, должно быть два аргумента UUID и SECRET")
		os.Exit(1)
	}

	current := NewApiClient(os.Args[1], os.Args[2])
	current.get_users()

}
