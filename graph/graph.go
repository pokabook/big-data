package graph

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
	"pokabook/big-data/crawling"
)

type TechstackCount struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Count    int    `json:"count"`
}

func CountTechstacks(techstacks []crawling.Techstack) (result []TechstackCount) {
	counts := make(map[string]map[string]int)

	for _, tech := range techstacks {
		if _, ok := counts[tech.Name]; !ok {
			counts[tech.Name] = make(map[string]int)
		}
		counts[tech.Name][tech.Category]++
	}

	for name, categories := range counts {
		for category, count := range categories {
			result = append(result, TechstackCount{name, category, count})
		}
	}

	return
}

func FindMaxCountPerCategory(techstacks []TechstackCount) (result []TechstackCount) {
	maxCounts := make(map[string]TechstackCount)

	for _, tech := range techstacks {
		if maxTech, ok := maxCounts[tech.Category]; !ok || tech.Count > maxTech.Count {
			maxCounts[tech.Category] = tech
		}
	}

	for _, tech := range maxCounts {
		result = append(result, tech)
	}

	return result
}

func GenerateBar(topTechstacks []TechstackCount) {

	bar := charts.NewBar()

	bar.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{Title: "기술 스택 사용량"},
		),
		charts.WithInitializationOpts(
			opts.Initialization{
				Width:           "1600px",
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
			},
		),
	)

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

	f, _ := os.Create("bar.html")
	bar.Render(f)
}

func GeneratePie(techstacks []TechstackCount) {

	categoryData := make(map[string][]TechstackCount)
	for _, techStack := range techstacks {
		categoryData[techStack.Category] = append(categoryData[techStack.Category], techStack)
	}

	page := components.NewPage()

	for category, techStacks := range categoryData {
		pie := charts.NewPie()

		total := 0
		for _, techStack := range techStacks {
			total += techStack.Count
		}

		items := make([]opts.PieData, 0, len(techStacks))
		others := 0.0
		for _, techStack := range techStacks {
			percentage := float64(techStack.Count) / float64(total) * 100
			if percentage <= 2.5 {
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
			}),
			charts.WithInitializationOpts(
				opts.Initialization{
					BackgroundColor: "#000000",
					Theme:           "dark",
				}),
			charts.WithLegendOpts(opts.Legend{
				Show:   true,
				Orient: "vertical",
				Right:  "right",
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
				},
			),
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

		page.AddCharts(pie)
	}

	f, _ := os.Create("pie.html")
	page.Render(f)
}
