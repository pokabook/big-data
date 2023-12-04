package utils

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

type TechstackCount struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Count    int    `json:"count"`
}
