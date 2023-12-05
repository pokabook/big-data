package graph

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
	"pokabook/big-data/utils"
	"sort"
)

func generateBar(topTechstacks []utils.TechstackCount) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: "기술 스택 사용량",
				Left:  "60px",
			},
		),
		charts.WithInitializationOpts(
			opts.Initialization{
				Width:           "80%",
				PageTitle:       "기술 스택 사용량",
				BackgroundColor: "#000000",
				Theme:           "dark",
			},
		),
		charts.WithTooltipOpts(
			opts.Tooltip{
				Show:    true,
				Trigger: "axis",
			},
		),
		charts.WithToolboxOpts(
			opts.Toolbox{
				Show: true,
				Feature: &opts.ToolBoxFeature{
					SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
						Show:  true,
						Title: "Save as Png",
					},
				},
			}),
	)

	sort.Slice(topTechstacks, func(i, j int) bool {
		return topTechstacks[i].Count > topTechstacks[j].Count
	})

	names := make([]string, 0)
	values := make([]opts.BarData, 0)
	for _, tech := range topTechstacks {
		names = append(names, tech.Name+"\n("+tech.Category+")")
		values = append(values, opts.BarData{Value: tech.Count})
	}

	bar.SetXAxis(names).AddSeries("사용하는 기업 수", values).SetSeriesOptions(
		charts.WithLineChartOpts(
			opts.LineChart{
				Smooth: true,
			}),
		charts.WithLabelOpts(
			opts.Label{
				Show:       true,
				Position:   "top",
				FontWeight: "normal",
			},
		),
		charts.WithSeriesAnimation(
			true,
		),
	)

	return bar
}

func generatePie(category string, techStacks []utils.TechstackCount) *charts.Pie {
	pie := charts.NewPie()

	total := 0
	others := 0.0

	sort.Slice(techStacks, func(i, j int) bool {
		return techStacks[i].Count > techStacks[j].Count
	})

	items := make([]opts.PieData, 0, len(techStacks))
	for _, techStack := range techStacks {
		total += techStack.Count
	}

	for _, techStack := range techStacks {
		percentage := float64(techStack.Count) / float64(total) * 100
		if percentage <= 2 {
			others += percentage
		} else {
			items = append(items, opts.PieData{
				Value: percentage,
				Name:  techStack.Name,
			})
		}
	}

	if others > 0 {
		items = append(items, opts.PieData{
			Value: others,
			Name:  "기타",
		})
	}

	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("%s 기술 스택 사용량", category),
			Left:  "60px",
		}),
		charts.WithInitializationOpts(
			opts.Initialization{
				PageTitle:       "기술 스택 사용량",
				Width:           "80%",
				BackgroundColor: "#000000",
				Theme:           "dark",
			}),
		charts.WithLegendOpts(opts.Legend{
			Show:   true,
			Orient: "vertical",
			Left:   "70%",
			Top:    "middle",
		}),
		charts.WithToolboxOpts(
			opts.Toolbox{
				Show: true,
				Feature: &opts.ToolBoxFeature{
					SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
						Show:  true,
						Title: "Save as Png",
					},
				},
			}),
		charts.WithTooltipOpts(
			opts.Tooltip{
				Show:      true,
				Trigger:   "item",
				Formatter: "{b} : {d}%",
			}),
	)

	pie.AddSeries(category, items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      true,
				Formatter: "{b}: {d}%",
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius: []string{"30%", "75%"},
				Center: []string{"35%", "50%"},
			}),
			charts.WithSeriesAnimation(
				true,
			),
		)
	return pie
}

func generateWordCloud(techstacks []utils.TechstackCount) *charts.WordCloud {

	wordCloud := charts.NewWordCloud()
	wordCloud.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "기술 스택",
		}),
	)

	data := make([]opts.WordCloudData, len(techstacks))
	for i, techStack := range techstacks {
		data[i] = opts.WordCloudData{
			Name:  techStack.Name,
			Value: techStack.Count,
		}
	}

	wordCloud.AddSeries("기술 스택", data)

	return wordCloud
}

func GenerateGraph(topTechstacks []utils.TechstackCount, techstacks []utils.TechstackCount) {
	categoryData := make(map[string][]utils.TechstackCount)
	for _, techStack := range techstacks {
		categoryData[techStack.Category] = append(categoryData[techStack.Category], techStack)
	}

	page := components.NewPage()
	page.PageTitle = "기업별 사용 기술 스택 분석"
	page.AddCharts(generateWordCloud(topTechstacks))
	page.AddCharts(generateBar(topTechstacks))

	for category, techStacks := range categoryData {
		page.AddCharts(generatePie(category, techStacks))
	}

	page.AddCustomizedCSSAssets("./css/graph.css")
	f, _ := os.Create("index.html")
	page.Render(f)
}
