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
	start := time.Now()

	var companies []utils.CompanyId
	var companiesInfo []utils.Techstack
	var wg sync.WaitGroup

	pages := int(math.Round(crawling.GetCompanyListNum() / 12))

	wg.Add(pages)
	companyListChannel := make(chan []utils.CompanyId)

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

	//progressBar := utils.CreateProgressBar()

	companyInfoChannel := make(chan []utils.Techstack)
	wg.Add(len(companies))

	for _, company := range companies {
		go func(companyId int) {
			crawling.GetCompanyInfo(companyId, companyInfoChannel)
			wg.Done()
		}(company.ID)
	}

	go func() {
		wg.Wait()
		close(companyInfoChannel)
	}()

	for companyInfo := range companyInfoChannel {
		//progressBar.Add(len(companyInfo))
		companiesInfo = append(companiesInfo, companyInfo...)
	}

	fmt.Println("NUMBER_OF_DATA=", len(companiesInfo))
	fmt.Println("TOTAL_TIME=", time.Since(start))

	countedTechstacks := utils.CountTechstacks(companiesInfo)
	topTechstacks := utils.FindMaxCountPerCategory(countedTechstacks)

	graph.GenerateGraph(topTechstacks, countedTechstacks)
}
