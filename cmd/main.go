package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"

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
	rowIterator, errRow := f.Rows(sheetName)
	if errRow != nil {
		fmt.Println("some problems on my ROW function")
	}
	defer func() {
		if errDefer := rowIterator.Close(); errDefer != nil {
			fmt.Println(errDefer)
		}
	}()
	var firstColumnValue []string
	var column []string
	var errColumn error
	var lines []string
	for rowIterator.Next() {
		column, errColumn = rowIterator.Columns()
		if errColumn != nil {
			fmt.Print("error to get columns")
		}
		if len(column) == 0 || column[0] == "" {
			continue
		} else {
			lines = splitLines(column[0])
		}
		firstColumnValue = append(firstColumnValue, lines...)

	}
	fmt.Println(firstColumnValue)

	for i, row := range rows[1:] {
		var content string
		for j, cell := range row {
			content += fmt.Sprintf("%s: %s\n", headers[j], cell)
		}

		// Criar QR code como PNG temporário
		qrPath := filepath.Join(outputDir, fmt.Sprintf("%s.png", firstColumnValue[i+1]))
		qrCode, err := qrcode.New(content, qrcode.Low)
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
		y := 110.0
		for _, line := range splitLines(content) {
			pdf.Text(105, y, line)
			y += 5
		}
		pdfPath := filepath.Join(outputDir, fmt.Sprintf("%s.pdf", firstColumnValue[i+1]))
		err = pdf.OutputFileAndClose(pdfPath)
		if err != nil {
			log.Fatalf("Erro ao salvar PDF: %v", err)
		}

		os.Remove(qrPath)
		fmt.Printf("PDF gerado: %s\n", pdfPath)
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
