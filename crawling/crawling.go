package crawling

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokabook/big-data/utils"
	"strconv"
	"time"
)

type CompanyTotalNum struct {
	CompanyTotal int `json:"company_total"`
}
type Company struct {
	CompanyId int `json:"company_id"`
}

type CompanyInfo struct {
	TechstackList []Techstack `json:"techstack_list"`
}

type Techstack struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

var BaseUrl = "https://api.codenary.co.kr"

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	},
}

func GetCompanyListNum() float64 {
	var companyTotal CompanyTotalNum
	res, err := httpClient.Get(BaseUrl + "/company")
	utils.CheckErr(err)
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	utils.CheckErr(err)
	err = json.Unmarshal(bodyBytes, &companyTotal)
	fmt.Println("기업 수 : ", companyTotal.CompanyTotal)
	return float64(companyTotal.CompanyTotal)
}

func GetCompanyList(pageNum int, companyListChannel chan<- []Company) {
	var temp []Company
	res, err := httpClient.Get(BaseUrl + "/company/list?page=" + strconv.Itoa(pageNum))
	utils.CheckErr(err)
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	utils.CheckErr(err)

	err = json.Unmarshal(bodyBytes, &temp)
	utils.CheckErr(err)
	companyListChannel <- temp
}

func GetCompanyInfo(companyId int, companyInfoChannel chan<- []Techstack) {
	var temp CompanyInfo
	res, err := httpClient.Get(BaseUrl + "/company/detail/" + strconv.Itoa(companyId))
	utils.CheckErr(err)
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	utils.CheckErr(err)

	err = json.Unmarshal(bodyBytes, &temp)
	utils.CheckErr(err)
	companyInfoChannel <- temp.TechstackList
}
