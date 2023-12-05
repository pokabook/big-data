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

	var companies []utils.Company
	var companiesInfo []utils.Techstack
	var wg sync.WaitGroup

	pages := int(math.Round(crawling.GetCompanyListNum() / 12))

	wg.Add(pages)
	companyListChannel := make(chan []utils.Company)

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
		}(company.CompanyId)
	}

	go func() {
		wg.Wait()
		close(companyInfoChannel)
	}()

	for companyInfo := range companyInfoChannel {
		//progressBar.Add(len(companyInfo))
		companiesInfo = append(companiesInfo, companyInfo...)
	}

	fmt.Println(`echo "NUMBER_OF_DATA=`, "데이터 수: ", len(companiesInfo), `" >> $GITHUB_OUTPUT`)
	fmt.Println(`echo "TOTAL_TIME=`, "걸린 시간: ", time.Since(start), `" >> $GITHUB_OUTPUT`)

	countedTechstacks := utils.CountTechstacks(companiesInfo)
	topTechstacks := utils.FindMaxCountPerCategory(countedTechstacks)

	graph.GenerateGraph(topTechstacks, countedTechstacks)
}
