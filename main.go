package main

import (
	"fmt"
	"math"
	"pokabook/big-data/crawling"
	"pokabook/big-data/graph"
	"pokabook/big-data/utils"
	"sync"
	"time"
)

func main() {
	fmt.Println("크롤링 시작")
	start := time.Now()

	var companies []crawling.Company
	var companiesInfo []crawling.Techstack
	var wg sync.WaitGroup

	pages := int(math.Round(crawling.GetCompanyListNum() / 12))

	wg.Add(pages)
	companyListChannel := make(chan []crawling.Company)

	for i := 1; i <= pages; i++ {
		go func(i int) {
			crawling.GetCompanyList(i, companyListChannel)
			wg.Done()
		}(i)
	}
	go func() {
		wg.Wait()
		close(companyListChannel)
	}()

	for companyList := range companyListChannel {
		companies = append(companies, companyList...)
	}

	progressBar := utils.CreateProgressBar()

	companyInfoChannel := make(chan []crawling.Techstack)
	wg.Add(len(companies))

	for _, company := range companies {
		go func(companyId int) {
			crawling.GetCompanyInfo(companyId, companyInfoChannel)
			wg.Done()
		}(company.CompanyId)
	}

	go func() {
		wg.Wait()
		close(companyInfoChannel)
	}()

	for companyInfo := range companyInfoChannel {
		progressBar.Add(len(companyInfo))
		companiesInfo = append(companiesInfo, companyInfo...)
	}

	fmt.Println()
	fmt.Println("데이터 개수 : ", len(companiesInfo))
	fmt.Println("크롤링 완료 시간 : ", time.Since(start))

	countedTechstacks := graph.CountTechstacks(companiesInfo)
	topTechstacks := graph.FindMaxCountPerCategory(countedTechstacks)

	graph.GenerateBar(topTechstacks)
	graph.GeneratePie(countedTechstacks)
}
