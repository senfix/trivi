package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/senfix/trivi-upload/trivi"

	"github.com/senfix/abo"
)

const (
	BANK_ACCOUNT_NO   = "1427237044/3030"
	BANC_ACCOUNT_NAME = "Živnostenský účet"
)

func main() {
	balance := 50073.58
	balance = exportMonth("vypis_10.csv", "1910_vypis.gpc", balance, "2019-10-01", "2019-12-31")
	balance = exportMonth("vypis_11.csv", "1911_vypis.gpc", balance, "2019-11-01", "2019-12-30")
	balance = exportMonth("vypis_12.csv", "1912_vypis.gpc", balance, "2019-12-01", "2019-12-31")

	uploadStatementToTrivi("1910_vypis.gpc")
	uploadStatementToTrivi("1911_vypis.gpc")
	//uploadStatementToTrivi("1912_vypis.gpc")
}

func uploadStatementToTrivi(filename string) {
	tc := trivi.NewTriviApi()
	err, _ := tc.Authenticate(trivi.AuthenticationRequest{
		AppID:     TRIVI_APP_ID,
		AppSecret: TRIVI_APP_SECRET,
	})
	if err != nil {
		panic(err)
	}

	err, resp := tc.UploadBankStatement(fmt.Sprintf("data/%v", filename))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", resp)
}

func exportMonth(sourceName, destinationName string, pocatecniZustatek float64, from, to string) (endBalance float64) {

	csvFile, _ := os.Open(fmt.Sprintf("data/%v", sourceName))
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ';'

	exportFrom, _ := time.Parse("2006-01-02", from)
	exportTo, _ := time.Parse("2006-01-02", to)

	export := abo.New(
		BANK_ACCOUNT_NO,
		BANC_ACCOUNT_NAME,
		pocatecniZustatek,
		exportFrom,
		exportTo,
	)

	first := true
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if first {
			first = false
			continue
		}
		castka, _ := strconv.ParseFloat(strings.Replace(line[5], ",", ".", -1), 64)
		valuta, _ := time.Parse("02/01/2006", line[31])
		datumZauctovani, _ := time.Parse("02/01/2006", line[31])

		export.AddTransaction(abo.Transaction{
			CisloProtiuctu:   line[10],
			CisloDokladu:     line[32],
			Castka:           castka,
			VariabilniSymbol: line[12],
			KonstantniSymbol: line[13],
			SpecifickySymbol: line[14],
			Valuta:           valuta,
			DoplnujiciUdaj:   line[9],
			Mena:             line[4],
			DatumZauctovani:  datumZauctovani,
		})
	}

	data := export.Generate()
	writeFile(fmt.Sprintf("data/%v", destinationName), data)

	endBalance = export.GetEndBalance()

	fmt.Printf("month %v started with %.2f ended with: %.2f\n", exportFrom.Format("01"), pocatecniZustatek, endBalance)

	return
}

func writeFile(path string, data []byte) {

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(data)
}
