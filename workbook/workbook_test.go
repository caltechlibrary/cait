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
	"bytes"
	"path"
	"testing"
)

func TestNew(t *testing.T) {
	wb := New()
	if wb == nil {
		t.Errorf("workbook.New() failed")
		t.FailNow()
	}
	expected := []byte(`<Workbook></Workbook>`)

	result := []byte(wb.String())
	if bytes.Equal(expected, result) != true {
		t.Errorf("\n%s\n!=\n%s", expected, result)
	}
}

func TestNewFromExcelFilename(t *testing.T) {
	fname := path.Join("testdata", "example1.xlsx")
	wb, err := NewFromExcelFilename(fname)
	if err != nil || wb == nil {
		t.Errorf("workbook.FromExcelFilename(%q) failed", fname)
		t.FailNow()
	}
	expected := []byte(`<Workbook><filename>example1.xlsx</filename><sheets sheet_name="Sheet1" sheet_no="1"><row row_no="1"><col row_no="1" col_no="1">One</col><col row_no="1" col_no="2">Two </col><col row_no="1" col_no="3">Three</col></row><row row_no="2"><col row_no="2" col_no="1">1</col><col row_no="2" col_no="2">2</col><col row_no="2" col_no="3">3</col></row></sheets></Workbook>`)

	result := []byte(wb.String())
	if bytes.Equal(expected, result) != true {
		t.Errorf("\n%s\n!=\n%s", expected, result)
	}
}

func TestToContainer(t *testing.T) {
	fname := path.Join("testdata", "example1.xlsx")
	wb, err := NewFromExcelFilename(fname)
	if err != nil || wb == nil {
		t.Errorf("workbook.FromExcelFilename(%q) failed", fname)
		t.FailNow()
	}

	container, err := wb.ToContainer()
	if err != nil {
		t.Errorf("workbook.ToContainer() failed")
		t.FailNow()
	}
	if container == nil {
		t.Errorf("workbook.ToContainer(), container is nil, failed")
		t.FailNow()
	}
}
