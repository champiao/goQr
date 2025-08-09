package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/champiao/goQr/helper"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Uso: %s path/to/archive.xlsx\n", os.Args[0])
		os.Exit(1)
	}

	excelPath := os.Args[1]
	if _, err := os.Stat(excelPath); os.IsNotExist(err) {
		log.Fatalf("Arquivo não encontrado: %s", excelPath)
	}

	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		log.Fatalf("Erro ao ler arquivo Excel: %v", err)
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Erro ao ler linhas da planilha: %v", err)
	}

	outputDir := fmt.Sprintf("qrcodes_output/%s", sheetName)
	os.MkdirAll(outputDir, os.ModePerm)

	headers := rows[0]
	hIndex := map[string]int{}
	for i, h := range headers {
		hIndex[strings.ToLower(strings.TrimSpace(h))] = i
	}
	rowIterator, errRow := f.Rows(sheetName)
	if errRow != nil {
		fmt.Println("some problems on my ROW function")
	}
	defer func() {
		if errDefer := rowIterator.Close(); errDefer != nil {
			fmt.Println(errDefer)
		}
	}()

	for i, row := range rows[1:] {
		var content string
		for j, cell := range row {
			content += fmt.Sprintf("%s: %s\n", headers[j], cell)
		}

		to := helper.GetCellByHeader(row, hIndex, "e-mail to")
		if to == "" {
			log.Printf("Linha %d: campo 'to' vazio, ignorando linha.", i+1)
			continue
		}
		subject := helper.GetCellByHeader(row, hIndex, "Planta")
		repCode := helper.GetCellByHeader(row, hIndex, "REP")
		area := helper.GetCellByHeader(row, hIndex, "Área")
		loc := helper.GetCellByHeader(row, hIndex, "Localização")
		body := fmt.Sprintf("%s \n %s \n %s \n Informe o problema:", repCode, area, loc)
		mailto := helper.BuildMailtoURI(to, subject, body)

		// Criar QR code como PNG temporário
		qrPath := filepath.Join(outputDir, fmt.Sprintf("%s.png", repCode))
		qrCode, err := qrcode.New(mailto, qrcode.Low)
		if err != nil {
			log.Fatalf("Erro ao gerar QR code: %v", err)
		}

		qrFile, err := os.Create(qrPath)
		if err != nil {
			log.Fatalf("Erro ao criar arquivo PNG: %v", err)
		}
		png.Encode(qrFile, qrCode.Image(256))
		qrFile.Close()

		// Criar PDF e inserir QR code e texto
		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.ImageOptions(qrPath, 80, 50, 50, 50, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
		pdf.SetFont("Arial", "", 10)
		// y := 110.0
		// for _, line := range splitLines(content) {
		// 	pdf.Text(105, y, line)
		// 	y += 5
		// }
		pdfPath := filepath.Join(outputDir, fmt.Sprintf("%s_%s.pdf", subject, repCode))
		err = pdf.OutputFileAndClose(pdfPath)
		if err != nil {
			log.Fatalf("Erro ao salvar PDF: %v", err)
		}

		os.Remove(qrPath)
		fmt.Printf("PDF gerado: %s\n", pdfPath)
	}
}
