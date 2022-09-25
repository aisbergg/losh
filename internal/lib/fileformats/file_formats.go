// Copyright 2022 André Lehmann
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package fileformats contains lists of common file formats used for OSH
// projects.
//
// The list of file formats is taken from:
//  - https://gitlab.com/OSEGermany/osh-tool/-/tree/master/data
//  - https://github.com/hoijui/file-extension-list
package fileformats

import "strings"

type fileExtensionCategory struct {
	category string
}

const (
	SourceFile uint8 = iota
	ExportFile
)

var (
	TextFileExtensions = map[string]struct{}{
		"doc":   {},
		"docx":  {},
		"ebook": {},
		"log":   {},
		"md":    {},
		"msg":   {},
		"odt":   {},
		"org":   {},
		"pages": {},
		"pdf":   {},
		"rtf":   {},
		"rst":   {},
		"tex":   {},
		"txt":   {},
		"wpd":   {},
		"wps":   {},
	}

	ImageFileExtensions = map[string]struct{}{
		"jpg":  {},
		"jpeg": {},
		"png":  {},
		"3dm":  {},
		"3ds":  {},
		"max":  {},
		"bmp":  {},
		"dds":  {},
		"gif":  {},
		"psd":  {},
		"xcf":  {},
		"tga":  {},
		"thm":  {},
		"tif":  {},
		"tiff": {},
		"yuv,": {},
		"ai":   {},
		"eps":  {},
		"ps":   {},
		"svg":  {},
		"dwg":  {},
		"dxf":  {},
		"gpx":  {},
		"kml":  {},
		"kmz":  {},
		"webp": {},
	}

	SheetFileExtensions = map[string]struct{}{
		"ods":  {},
		"xls":  {},
		"xlsx": {},
		"csv":  {},
		"ics":  {},
		"vcf":  {},
	}

	CodeFileExtensions = map[string]struct{}{
		"1.ada":     {},
		"2.ada":     {},
		"ada":       {},
		"adb":       {},
		"ads":       {},
		"asm":       {},
		"bas":       {},
		"bash":      {},
		"bat":       {},
		"c":         {},
		"c++":       {},
		"cbl":       {},
		"cc":        {},
		"class":     {},
		"clj":       {},
		"cob":       {},
		"cpp":       {},
		"cs":        {},
		"csh":       {},
		"cxx":       {},
		"d":         {},
		"diff":      {},
		"e":         {},
		"el":        {},
		"f":         {},
		"f77":       {},
		"f90":       {},
		"fish":      {},
		"for":       {},
		"fth":       {},
		"ftn":       {},
		"go":        {},
		"groovy":    {},
		"h":         {},
		"hh":        {},
		"hpp":       {},
		"hs":        {},
		"html":      {},
		"htm":       {},
		"hxx":       {},
		"java":      {},
		"js":        {},
		"jsx":       {},
		"jsp":       {},
		"ksh":       {},
		"kt":        {},
		"lhs":       {},
		"lisp":      {},
		"lua":       {},
		"m":         {},
		"m4":        {},
		"nim":       {},
		"patch":     {},
		"php":       {},
		"pl":        {},
		"po":        {},
		"pp":        {},
		"py":        {},
		"r":         {},
		"rb":        {},
		"rs":        {},
		"s":         {},
		"scala":     {},
		"sh":        {},
		"swg":       {},
		"swift":     {},
		"v":         {},
		"vb":        {},
		"vcxproj":   {},
		"xcodeproj": {},
		"xml":       {},
		"zsh":       {},
	}

	CADFileExtensions = map[string]uint8{
		"3dm":        SourceFile,
		"3dxml":      ExportFile,
		"3ko":        ExportFile,
		"3mf":        ExportFile,
		"amf":        ExportFile,
		"asab":       ExportFile,
		"asat":       ExportFile,
		"asm":        SourceFile,
		"catpart":    SourceFile,
		"catproduct": SourceFile,
		"cgr":        SourceFile,
		"csg":        ExportFile,
		"dae":        ExportFile,
		"dgn":        SourceFile,
		"dwg":        SourceFile,
		"dxf":        SourceFile,
		"fcstd":      SourceFile,
		"html":       ExportFile,
		"iam":        SourceFile,
		"iges":       ExportFile,
		"igs":        ExportFile,
		"ipt":        SourceFile,
		"iwb":        ExportFile,
		"iwp":        ExportFile,
		"jt":         ExportFile,
		"j_t":        ExportFile,
		"model":      SourceFile,
		"obj":        ExportFile,
		"off":        ExportFile,
		"par":        SourceFile,
		"pdf":        ExportFile,
		"ply":        ExportFile,
		"pod":        ExportFile,
		"prc":        ExportFile,
		"prt":        SourceFile,
		"psm":        SourceFile,
		"sab":        ExportFile,
		"sat":        ExportFile,
		"scad":       SourceFile,
		"sldasm":     SourceFile,
		"sldprt":     SourceFile,
		"sms":        ExportFile,
		"step":       ExportFile,
		"stl":        ExportFile,
		"stp":        ExportFile,
		"svg":        ExportFile,
		"u3d":        ExportFile,
		"vda":        ExportFile,
		"wrl":        ExportFile,
		"x_t":        ExportFile,
		"xcgm":       ExportFile,
	}

	PCBFileExtensions = map[string]uint8{
		"brd":       SourceFile,
		"pro":       SourceFile,
		"sch":       SourceFile,
		"kicad_pcb": SourceFile,
	}
)

// IsTextFile returns true if the file extension is a text file.
func IsTextFile(extension string) bool {
	return isFileExtension(extension, TextFileExtensions)
}

// IsImageFile returns true if the file extension is an image file extension.
func IsImageFile(extension string) bool {
	return isFileExtension(extension, ImageFileExtensions)
}

// IsSheetFile returns true if the file extension is a sheet file extension.
func IsSheetFile(extension string) bool {
	return isFileExtension(extension, SheetFileExtensions)
}

// IsCodeFile returns true if the file extension is a code file extension.
func IsCodeFile(extension string) bool {
	return isFileExtension(extension, CodeFileExtensions)
}

// IsCADFile returns true if the file extension is a CAD file extension.
func IsCADFile(extension string) bool {
	return isFileExtension2(extension, CADFileExtensions)
}

// IsPCBFile returns true if the file extension is a PCB file extension.
func IsPCBFile(extension string) bool {
	return isFileExtension2(extension, PCBFileExtensions)
}

func isFileExtension(extension string, extensions map[string]struct{}) bool {
	extension = strings.ToLower(extension)
	if strings.HasPrefix(extension, ".") {
		extension = extension[1:]
	}
	_, ok := extensions[extension]
	return ok
}

func isFileExtension2(extension string, extensions map[string]uint8) bool {
	extension = strings.ToLower(extension)
	if strings.HasPrefix(extension, ".") {
		extension = extension[1:]
	}
	_, ok := extensions[extension]
	return ok
}

func IsSourceFile(extension string) bool {
	c, ok := CADFileExtensions[extension]
	if ok {
		return c == SourceFile
	}
	c, ok = PCBFileExtensions[extension]
	if ok {
		return c == SourceFile
	}
	return false
}
