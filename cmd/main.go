package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/champiao/goQr/helper"
	"github.com/jung-kurt/gofpdf"

	"github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
)

func main() {
	a := app.New()
	w := a.NewWindow("QR Code Generator")
	w.Resize(fyne.NewSize(400, 300))

	widthEntry := widget.NewEntry()
	widthEntry.SetPlaceHolder("Largura (mm)")
	heightEntry := widget.NewEntry()
	heightEntry.SetPlaceHolder("Comprimento (mm)")

	btn := widget.NewButton("Select XLSX File", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				log.Println("No file selected.")
				return
			}
			defer reader.Close()

			widthStr := widthEntry.Text
			heightStr := heightEntry.Text
			width, err := strconv.ParseFloat(widthStr, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("largura invalida: %v", err), w)
				return
			}
			height, err := strconv.ParseFloat(heightStr, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("comprimento invalido: %v", err), w)
				return
			}

			excelPath := reader.URI().Path()
			if _, err := os.Stat(excelPath); os.IsNotExist(err) {
				dialog.ShowError(fmt.Errorf("file not found: %s", excelPath), w)
				return
			}

			f, err := excelize.OpenFile(excelPath)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			sheetName := f.GetSheetName(0)
			rows, err := f.GetRows(sheetName, excelize.Options{RawCellValue: true})
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			outputDir := fmt.Sprintf("qrcodes_output/%s", sheetName)
			os.MkdirAll(outputDir, os.ModePerm)

			headers := rows[0]
			hIndex := make(map[string]int)
			for i, h := range headers {
				hIndex[strings.ToLower(strings.TrimSpace(h))] = i
			}

			for i, row := range rows[1:] {
				to := helper.GetCellByHeader(row, hIndex, "e-mail to")
				if to == "" {
					log.Printf("Row %d: 'to' field empty, skipping row.", i+2)
					continue
				}
				subject := helper.GetCellByHeader(row, hIndex, "Planta")
				repCode := helper.GetCellByHeader(row, hIndex, "REP")
				area := helper.GetCellByHeader(row, hIndex, "Área")
				loc := helper.GetCellByHeader(row, hIndex, "Localização")
				body := fmt.Sprintf("%s \n %s \n %s \n Informe o problema:", repCode, area, loc)
				mailto := helper.BuildMailtoURI(to, subject, body)

				qrPath := filepath.Join(outputDir, fmt.Sprintf("%s.png", repCode))
				qrCode, err := qrcode.New(mailto, qrcode.Low)
				if err != nil {
					log.Printf("Error generating QR code for row %d: %v", i+2, err)
					continue
				}

				qrFile, err := os.Create(qrPath)
				if err != nil {
					log.Printf("Error creating PNG file for row %d: %v", i+2, err)
					continue
				}
				png.Encode(qrFile, qrCode.Image(256))
				qrFile.Close()

				pdf := gofpdf.NewCustom(&gofpdf.InitType{
					UnitStr: "mm",
					Size:    gofpdf.SizeType{Wd: width, Ht: height},
				})
				pdf.AddPage()
				pdf.ImageOptions(qrPath, (width-50)/2, (height-50)/2, 50, 50, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
				pdf.SetFont("Arial", "", 10)
				y := (height-50)/2 - 5
				var info []string
				Newloc := strings.Split(loc, ":")
				loc = Newloc[1]
				info = append(info, repCode, loc)
				for _, newInfo := range info {
					helper.NormalizeText(newInfo)
				}
				for _, line := range info {
					pdf.Text((width-pdf.GetStringWidth(line))/2, y, line)
					y += 5
				}
				pdfPath := filepath.Join(outputDir, fmt.Sprintf("%s_%s.pdf", subject, repCode))
				err = pdf.OutputFileAndClose(pdfPath)
				if err != nil {
					log.Printf("Error saving PDF for row %d: %v", i+2, err)
					continue
				}

				os.Remove(qrPath)
				fmt.Printf("PDF generated: %s\n", pdfPath)
			}
			dialog.ShowInformation("Success", "QR codes generated successfully!", w)
		}, w)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		fileDialog.Show()
	})

	w.SetContent(container.NewVBox(
		widthEntry,
		heightEntry,
		btn,
	))

	w.ShowAndRun()
}
