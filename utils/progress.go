package utils

import "github.com/schollz/progressbar/v3"

func CreateProgressBar() *progressbar.ProgressBar {

	return progressbar.NewOptions(7963,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(70),
		progressbar.OptionSetDescription("[BigData] Crawling..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[blue]█[reset]",
			SaucerHead:    "[blue]█[reset]",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}))

}
