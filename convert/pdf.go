package convert

import (
	"fmt"
	"image/png"
	"os"

	"github.com/signintech/gopdf"
)

func PDF(filename string) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	imgFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening image file:", err)
		return
	}
	defer imgFile.Close()

	// 이미지 디코딩
	img, err := png.Decode(imgFile)
	if err != nil {
		fmt.Println("Error decoding image file:", err)
		return
	}

	// 이미지 크기 가져오기
	bounds := img.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	// PDF에 이미지 추가
	pdf.Image(filename, 0, 0, &gopdf.Rect{W: width, H: height})

	// PDF 저장
	err = pdf.WritePdf(filename + ".pdf")
	if err != nil {
		fmt.Println("Error saving PDF file:", err)
		return
	}
}
