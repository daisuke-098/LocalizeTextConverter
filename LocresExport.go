package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const ExcelFile = "ue4res_text.xlsx"
const SheetName = "ue4res_text"

//エクセルのテキスト開始行
const N = 2

/*
* 翻訳する言語を指定したiniファイルを読み込む
 */
func LoadCultureIni() (string, error) {
	f, err := os.Open("LocresExport.ini")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var culture string
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, ":") {
			culture = line
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return culture, nil
}

/*
* InputとOutputDirを定義したiniファイルを読み込む
 */
func LoadDirIni() (string, string, error) {
	f, err := os.Open("Directory.ini")
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var InputDir, OutputDir string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "InputDir=") {
			InputDir = strings.TrimLeft(line, "InputDir=")
		} else if strings.HasPrefix(line, "OutputDir=") {
			OutputDir = strings.TrimLeft(line, "OutputDir=")
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", err
	}
	return InputDir, OutputDir, nil
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func main() {
	//翻訳対象の言語を読み込み
	exportCultures, err := LoadCultureIni()
	if err != nil {
		log.Fatal(err)
	}
	cultureCodes := strings.Split(exportCultures, ",")

	//Input,OutputDirを読み込み
	InputDir, OutputDir, err := LoadDirIni()
	if err != nil {
		log.Fatal(err)
	}

	cultureCol := map[string]int{}
	var keyCol int

	//エクセルを読み込む
	excel, err := excelize.OpenFile(InputDir + ExcelFile)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := excel.GetRows(SheetName)
	for col, cellValue := range rows[1] {
		if contains(cultureCodes, cellValue) {
			cultureCol[cellValue] = col
		} else if strings.Index(cellValue, "テキストID") != -1 {
			keyCol = col
		}
	}

	manifestOutput := `{
	"FormatVersion": 1,
	"Namespace": "",
	"Children": [
`
	archiveOutput := `{
	"FormatVersion": 1,
	"Namespace": "",
	"Children": [
`
	var i int
	//タグが入らないようにNから
	for _, cultureCode := range cultureCodes {
		for i = N; i < len(rows); i++ {
			keyCell, _ := excelize.CoordinatesToCellName(keyCol+1, i)
			keyValue, err := excel.GetCellValue(SheetName, keyCell)
			if err != nil {
				log.Fatal(err)
			}
			if keyValue != "" || strings.HasPrefix(keyValue, "_") {
				cultureCell, _ := excelize.CoordinatesToCellName(cultureCol[cultureCode]+1, i)
				cultureValue, err := excel.GetCellValue(SheetName, cultureCell)
				if err != nil {
					log.Fatal(err)
				}
				manifestOutput += `		{
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
			},
	`
				archiveOutput += `		{
				"Source":
				{
					"Text": "` + keyValue + `"
				},
				"Translation":
				{
					"Text": "` + cultureValue + `"
				}
			},
	`
			}

		}
		strings.TrimRight(manifestOutput, ",")
		manifestOutput += `	]
	}`
		strings.TrimRight(archiveOutput, ",")
		archiveOutput += `	]
	}`

		//言語ごとのローカライズデータを出力する
		manifestFile, err := os.Create(OutputDir + "Localization/ue4res_text.manifest")
		if err != nil {
			log.Fatal(err)
		}
		defer manifestFile.Close()
		archiveFile, err := os.Create(OutputDir + "Localization/" + cultureCode + "/ue4res_text.archive")
		if err != nil {
			log.Fatal(err)
		}
		defer archiveFile.Close()

		manifestFile.Write(([]byte)(manifestOutput))
		archiveFile.Write(([]byte)(archiveOutput))
	}
}
