package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

//TagCol エクセルのtag記述行
const TagCol = 2

//Key key列を示すタグ文字
const Key = "key"

// contains スライスに文字列が存在するか判定
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// directoryExists 指定したディレクトリが無ければ作成する
func directoryExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
	}
}

func main() {
	//実行オプション定義
	var (
		outputDir      = flag.String("o", ".", "Output Directory")
		targetCultures = flag.String("l", "", "Loclize target Cultures")
	)
	flag.Parse()

	//翻訳対象の言語を設定
	cultureCodes := strings.Split(*targetCultures, ",")
	cultureCol := map[string]int{}

	var keyCol int
	//エクセルを読み込む
	for _, Arg := range flag.Args() {
		//Argに実行オプションも含まれるので、ハイフンから始まるものは無視する
		if !strings.HasPrefix(Arg, "-") {
			excel, err := excelize.OpenFile(Arg)
			if err != nil {
				log.Fatal(err)
			}

			var manifestOutput string
			for _, cultureCode := range cultureCodes {
				manifestOutput = `{
	"FormatVersion": 1,
	"Namespace": "",
	"Children": [`
				archiveOutput := `{
	"FormatVersion": 1,
	"Namespace": "",
	"Children": [`

				for _, sheet := range excel.GetSheetList() {
					rows, _ := excel.GetRows(sheet)
					for col, cellValue := range rows[TagCol-1] {
						//key列と各言語の列番号を取得
						if contains(cultureCodes, cellValue) {
							cultureCol[cellValue] = col + 1
						} else if strings.Contains(cellValue, Key) {
							keyCol = col + 1
						}
					}

					var i int
					//タグが入らないようにTagCol+1から
					for i = TagCol + 1; i < len(rows); i++ {
						keyCell, _ := excelize.CoordinatesToCellName(keyCol, i)
						keyValue, err := excel.GetCellValue(sheet, keyCell)
						if err != nil {
							log.Fatal(err)
						}
						//keyの先頭がアンダースコアの場合は無視する
						if keyValue != "" && !strings.HasPrefix(keyValue, "_") {
							cultureCell, _ := excelize.CoordinatesToCellName(cultureCol[cultureCode], i)
							cultureValue, err := excel.GetCellValue(sheet, cultureCell)
							if err != nil {
								log.Fatal(err)
							}
							manifestOutput += `
		{
			"Source":
			{
				"Text": "` + keyValue + `"
			},
			"Keys": [
				{
					"Key": "` + keyValue + `",
					"Path": ""
				}
			]
		},`
							archiveOutput += `
		{
			"Source":
			{
				"Text": "` + keyValue + `"
			},
			"Translation":
			{
				"Text": "` + cultureValue + `"
			}
		},`
						}
					}
				}
				manifestOutput = strings.TrimRight(manifestOutput, ",")
				manifestOutput += `
	]
}`
				archiveOutput = strings.TrimRight(archiveOutput, ",")
				archiveOutput += `
	]
}`

				//エクセルのファイル名取得
				//windows以外で実行したらバグる(多分)
				if strings.Contains(Arg, "\\") {
					Arg = Arg[strings.LastIndex(Arg, "\\"):]
				}
				directoryExists(*outputDir + "/" + cultureCode)
				archiveFile, err := os.Create(*outputDir + "/" + cultureCode + "/" + strings.TrimRight(Arg, ".xlsx") + ".archive")
				if err != nil {
					log.Fatal(err)
				}
				defer archiveFile.Close()

				archiveFile.Write(([]byte)(archiveOutput))
			}
			manifestFile, err := os.Create(*outputDir + "/" + strings.TrimRight(Arg, ".xlsx") + ".manifest")
			if err != nil {
				log.Fatal(err)
			}
			defer manifestFile.Close()

			manifestFile.Write(([]byte)(manifestOutput))
		}
	}
}
