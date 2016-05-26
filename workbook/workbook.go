//
// Package workbook supports creating EADs from Excel Workbooks that can be imported into ArchivesSpace
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package workbook

import (
	"encoding/xml"
	"fmt"
	"path"

	// Caltech Library packages

	// 3rd party packages
	"github.com/tealeg/xlsx"
)

type Workbook struct {
	XMLName  xml.Name `json:"-"`
	Filename string   `xml:"filename,omitempty" json:"filename,omitempty"`
	Sheets   []Sheet  `xml:"sheets,omitempty" json:"sheets,omitempty"`
}

type Sheet struct {
	XMLName xml.Name `json:"-"`
	Name    string   `xml:"sheet_name,attr,omitempty" json:"sheet_name,omitempty"`
	No      int      `xml:"sheet_no,attr,omitempty" json:"sheet_no,int,omitempty"`
	Rows    []Row    `xml:"row,omitempty" json:"row,omitempty"`
}

type Row struct {
	XMLName xml.Name `json:"-"`
	No      int      `xml:"row_no,attr,omitempty" json:"row_no,int,omitempty"`
	Cols    []Cell   `xml:"col,omitempty" json:"col,omitempty"`
}

type Cell struct {
	XMLName xml.Name `json:"-"`
	RowNo   int      `xml:"row_no,attr,omitempty" json:"row_no,int,omitempty"`
	ColNo   int      `xml:"col_no,attr,omitempty" json:"col_no,int,omitempty"`
	Value   string   `xml:",innerxml" json:"value,string"`
}

// New creates an empty Workbook object
func New() *Workbook {
	return new(Workbook)
}

// NewFromExcelFile creates a Workbook from an xlsx.File object
func NewFromExcelFile(xlFile *xlsx.File) (*Workbook, error) {
	wb := new(Workbook)
	for i, sheet := range xlFile.Sheets {
		newSheet := new(Sheet)
		newSheet.No = i + 1
		newSheet.Name = sheet.Name
		for j, row := range sheet.Rows {
			newRow := new(Row)
			newRow.No = j + 1
			for k, cell := range row.Cells {
				newCell := new(Cell)
				newCell.RowNo = j + 1
				newCell.ColNo = k + 1
				newCell.Value = cell.Value
				newRow.Cols = append(newRow.Cols, *newCell)
			}
			newSheet.Rows = append(newSheet.Rows, *newRow)
		}
		wb.Sheets = append(wb.Sheets, *newSheet)
	}
	return wb, nil
}

// NewFromExcelFilename creates a Workbook by opening a excel file given a filename
func NewFromExcelFilename(fname string) (*Workbook, error) {
	xlFile, err := xlsx.OpenFile(fname)
	if err != nil {
		return nil, err
	}
	//defer xlFile.Close()
	wb, err := NewFromExcelFile(xlFile)
	if err != nil {
		return nil, err
	}
	wb.Filename = path.Base(fname)
	return wb, nil
}

func (wb *Workbook) String() string {
	src, _ := xml.Marshal(wb)
	return string(src)
}

func (wb *Workbook) ToContainer() ([]string, error) {
	return nil, fmt.Errorf("ToContainer() not implemented.")
}
