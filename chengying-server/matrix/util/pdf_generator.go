// Licensed to Apache Software Foundation(ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation(ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package util

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
)

const (
	defaultOrientation = "P"
	defaultUnit        = "mm"
	defaultSize        = "A4"
	margin             = 10
)

type PdfGenerator struct {
	pdf        *gofpdf.Fpdf
	Id         int
	pageWidth  float64
	pageHeight float64
	left       float64
}

func NewPdfGenerator(id int) *PdfGenerator {
	pdfGenerator := &PdfGenerator{
		pdf: gofpdf.New(defaultOrientation, defaultUnit, defaultSize, ""),
	}
	pdfGenerator.pageWidth, pdfGenerator.pageHeight = pdfGenerator.pdf.GetPageSize()
	pdfGenerator.left = margin / 2
	pdfGenerator.Id = id
	return pdfGenerator
}

func (p *PdfGenerator) PageWidth() float64 {
	return p.pageWidth
}

func (p *PdfGenerator) PageHeight() float64 {
	return p.pageHeight
}

func (p *PdfGenerator) Left() float64 {
	return p.left
}

func (p *PdfGenerator) AddPage() {
	p.pdf.AddPage()
}

func (p *PdfGenerator) AddFont(fontFamily, fontFile string) {
	p.pdf.AddUTF8Font(fontFamily, "", fontFile)
}

func (p *PdfGenerator) SetFont(fontFamily, style string, size float64) {
	p.pdf.SetFont(fontFamily, style, size)
}

// 添加文本
func (p *PdfGenerator) AddText(x, fontSize, height float64, fontFamily, content string, r, g, b int) {
	if p.pdf.GetY()+height >= p.pageHeight-margin {
		p.AddPage()
	}
	p.pdf.SetX(x)
	p.pdf.SetFont(fontFamily, "", fontSize)
	p.pdf.SetTextColor(r, g, b)
	p.pdf.Cell(40, height, content)
	// 换行
	p.Ln(-1)
}

// 画线
func (p *PdfGenerator) AddLine(x float64) {
	p.pdf.SetDrawColor(30, 144, 255)
	p.pdf.SetLineWidth(.5)
	p.pdf.Line(x, p.pdf.GetY()+1, x, p.pdf.GetY()+4)
}

func (p *PdfGenerator) AddCircle(x, y, radius float64, r, g, b int) {
	p.pdf.SetFillColor(r, g, b)
	p.pdf.Circle(x, y, radius, "F")
}

// 添加表格
func (p *PdfGenerator) AddTable(datas [][]string, headers []string) {
	p.pdf.SetFillColor(175, 238, 238)
	p.pdf.SetDrawColor(188, 143, 143)
	p.pdf.SetLineWidth(.3)
	p.pdf.SetFont("Simhei", "", 0)
	// 表头
	p.pdf.SetX(p.left)
	if len(headers) > 0 {
		colCount := len(headers)
		widthSlice := CalTableHeaderWidth(colCount, p.pageWidth-margin)
		for j, str := range headers {
			p.pdf.CellFormat(widthSlice[j], 7, str, "1", 0, "C", true, 0, "")
		}
		p.pdf.Ln(-1)
		p.pdf.SetFillColor(224, 235, 255)
		p.pdf.SetTextColor(0, 0, 0)
		p.pdf.SetFont("Simhei", "", 0)
		fill := false
		for _, data := range datas {
			if len(data) == len(headers) {
				p.pdf.SetX(p.left)
				for i, column := range data {
					if column == "正常" {
						p.AddCircle(p.pdf.GetX()+widthSlice[i]/2-5, p.pdf.GetY()+3, 1, 50, 205, 50)
					}
					if column == "alerting" {
						p.AddCircle(p.pdf.GetX()+widthSlice[i]/2-8, p.pdf.GetY()+3, 1, 255, 0, 0)
					}
					if column == "no_data" {
						p.AddCircle(p.pdf.GetX()+widthSlice[i]/2-8, p.pdf.GetY()+3, 1, 225, 215, 0)
					}
					if column == "NORMAL" {
						column = "正常"
					} else if column == "ABNORMAL" {
						column = "异常"
					}
					p.pdf.CellFormat(widthSlice[i], 6, column, "1", 0, "C", false, 0, "")
				}
			}
			p.pdf.Ln(-1)
			fill = !fill
		}
		p.pdf.SetX(p.left)
		p.pdf.CellFormat(p.pageWidth-margin, 0, "", "1", 2, "", false, 0, "")
		p.Ln(4)
	}
}

func (p *PdfGenerator) Ln(h float64) {
	p.pdf.Ln(h)
}

func (p *PdfGenerator) AddFooter() {
	p.pdf.SetFooterFunc(func() {
		p.pdf.SetY(-15)
		p.pdf.SetFont("Simhei", "", 8)
		p.pdf.SetTextColor(128, 128, 128)
		p.pdf.CellFormat(0, 10, fmt.Sprintf("第 %d 页", p.pdf.PageNo()), "", 0, "C", false, 0, "")
	})
}

func (p *PdfGenerator) AddLineChart(imageType, imageStr string) {
	p.pdf.SetX(p.left)
	var opt gofpdf.ImageOptions
	opt.ImageType = imageType
	opt.AllowNegativePosition = false
	p.pdf.ImageOptions(imageStr, p.left, p.pdf.GetY(), p.pageWidth-margin, 60, true, opt, 0, "")
	p.Ln(4)
}

func (p *PdfGenerator) OutputFileAndClose(fileStr string) error {
	return p.pdf.OutputFileAndClose(fileStr)
}

// 根据表格列数确定每一列宽度
func CalTableHeaderWidth(count int, width float64) []float64 {
	var headerWidthSlice []float64
	interval := width / float64(count)
	for i := 0; i < count; i++ {
		headerWidthSlice = append(headerWidthSlice, interval)
	}
	return headerWidthSlice
}
