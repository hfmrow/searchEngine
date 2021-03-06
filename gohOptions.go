// gohOptions.go

// Source file auto-generated on Wed, 02 Oct 2019 17:33:10 using Gotk3ObjHandler v1.3.8 ©2018-19 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"

	gltsbh "github.com/hfmrow/genLib/tools/bench"
	gltses "github.com/hfmrow/genLib/tools/errors"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gidgat "github.com/hfmrow/gtk3Import/dialog/about"
	gimc "github.com/hfmrow/gtk3Import/misc"
	gitw "github.com/hfmrow/gtk3Import/treeview"
	// gidgcr "github.com/hfmrow/gtk3Import/dialog/chooser"
)

// App infos
var Name = "SearchEngine"
var Vers = "v1.8.5"
var Descr = "This program is designed to search files over directory,\nsubdirectory, and retrieving information based on\ndate, type, patterns contained in name."
var Creat = "H.F.M"
var YearCreat = "2018-19"
var LicenseShort = "This program comes with absolutely no warranty.\nSee the The MIT License (MIT) for details:\nhttps://opensource.org/licenses/mit-license.php"
var LicenseAbrv = "License (MIT)"
var Repository = "github.com/hfmrow/searchEngine"

// Vars declarations
var absoluteRealPath, optFilename = getAbsRealPath()
var mainOptions *MainOpt
var devMode bool
var tempDir string
var doTempDir bool
var namingWidget bool

var errNoSelection = errors.New("There is no selection currently.")
var tvs *gitw.TreeViewStructure
var columnsNames = [][]string{{"Name", "text"}, {"Size", "text"}, {"Type", "text"}, {"Time", "text"}, {"Path", "text"}}
var timer = gltsbh.BenchNew(false)
var statusbar = gimc.StatusBar{}
var clipboard = gimc.Clipboard{}
var titlebar *gimc.TitleBar
var maintitle string
var filesScanned int

// Functions mapping
var Check = gltses.Check
var DialogMessage = gidg.DialogMessage

// var FileChooser = gidgcr.FileChooser

// To store original label content for newer than and older than buttons.
var origLabelNT, origLabelOT string

type searchTimeCal struct {
	y, m, d uint
	H, M, S float64
	Ready   bool
}

type searchList struct {
	And []string
	Or  []string
	Not []string
}

type MainOpt struct {
	/* Public, will be saved and restored */
	AboutOptions                *gidgat.AboutInfos
	MainWinWidth, MainWinHeight int
	LanguageFilename            string
	LastDirectory               string
	CaseSensitive               bool
	CharClass                   bool
	CharClasStrict              bool
	WildCard                    bool
	Regex                       bool
	FollowSymlinks              bool
	FileType                    int
	DateType                    int
	Depth                       int
	WordAnd                     bool
	WordOr                      bool
	WordNot                     bool
	SplitAnd                    bool
	SplitOr                     bool
	SplitNot                    bool
	UpdateOnChanges             bool
	SearchList                  searchList
	AppLauncher                 string
	WebSearchEngine             string
	FileExplorer                string

	/* Private, will NOT be saved */
	searchNewerThan  searchTimeCal
	searchOlderThan  searchTimeCal
	foundFilesList   []string
	displayFilesList [][]string
}

// Main options initialisation
func (opt *MainOpt) Init() {
	opt.AboutOptions = new(gidgat.AboutInfos)

	opt.LanguageFilename = "assets/lang/eng.lang"

	opt.WebSearchEngine = `https://www.google.com/search?q=`
	opt.FileExplorer = "thunar"
	opt.AppLauncher = "xdg-open"

	opt.MainWinWidth = 800
	opt.MainWinHeight = 600
}

// Variables -> Objects.
func (opt *MainOpt) UpdateObjects() {
	mainObjects.MainWindow.Resize(opt.MainWinWidth, opt.MainWinHeight)

	mainObjects.SearchFilechooserbutton.SetFilename(opt.LastDirectory)
	mainObjects.SearchCheckbuttonCaseSensitive.SetActive(opt.CaseSensitive)
	mainObjects.SearchCheckbuttonWildCard.SetActive(opt.WildCard)
	mainObjects.SearchCheckbuttonRegex.SetActive(opt.Regex)
	mainObjects.SearchCheckbuttonCharClasses.SetActive(opt.CharClass)
	mainObjects.SearchCheckbuttonCharClassesStrictMode.SetActive(opt.CharClasStrict)

	mainObjects.SearchSpinbuttonDepth.SetValue(float64(opt.Depth))
	mainObjects.SearchCheckbuttonWordAnd.SetActive(opt.WordAnd)
	mainObjects.SearchCheckbuttonWordOr.SetActive(opt.WordOr)
	mainObjects.SearchCheckbuttonWordNot.SetActive(opt.WordNot)
	mainObjects.SearchCheckbuttonSplitedAnd.SetActive(opt.SplitAnd)
	mainObjects.SearchCheckbuttonSplitedOr.SetActive(opt.SplitOr)
	mainObjects.SearchCheckbuttonSplitedNot.SetActive(opt.SplitNot)

	mainObjects.SearchComboboxTextType.SetActive(opt.FileType)
	mainObjects.SearchComboboxTextDateType.SetActive(opt.DateType)

	mainObjects.SearchCheckbuttonFollowSL.SetActive(opt.FollowSymlinks)
}

// Objects -> Variables.
func (opt *MainOpt) UpdateOptions() {

	opt.MainWinWidth, opt.MainWinHeight = mainObjects.MainWindow.GetSize()

	opt.LastDirectory = mainObjects.SearchFilechooserbutton.GetFilename()
	opt.CaseSensitive = mainObjects.SearchCheckbuttonCaseSensitive.GetActive()
	opt.WildCard = mainObjects.SearchCheckbuttonWildCard.GetActive()
	opt.Regex = mainObjects.SearchCheckbuttonRegex.GetActive()
	opt.CharClass = mainObjects.SearchCheckbuttonCharClasses.GetActive()
	opt.CharClasStrict = mainObjects.SearchCheckbuttonCharClassesStrictMode.GetActive()

	opt.Depth = int(mainObjects.SearchSpinbuttonDepth.GetValue())
	opt.WordAnd = mainObjects.SearchCheckbuttonWordAnd.GetActive()
	opt.WordOr = mainObjects.SearchCheckbuttonWordOr.GetActive()
	opt.WordNot = mainObjects.SearchCheckbuttonWordNot.GetActive()
	opt.SplitAnd = mainObjects.SearchCheckbuttonSplitedAnd.GetActive()
	opt.SplitOr = mainObjects.SearchCheckbuttonSplitedOr.GetActive()
	opt.SplitNot = mainObjects.SearchCheckbuttonSplitedNot.GetActive()

	opt.FileType = mainObjects.SearchComboboxTextType.GetActive()
	opt.DateType = mainObjects.SearchComboboxTextDateType.GetActive()

	opt.FollowSymlinks = mainObjects.SearchCheckbuttonFollowSL.GetActive()
}

// Read Options from file
func (opt *MainOpt) Read() (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(optFilename); err == nil {
		err = json.Unmarshal(textFileBytes, &opt)
	}
	return err
}

// Write Options to file
func (opt *MainOpt) Write() (err error) {
	var out bytes.Buffer
	var jsonData []byte
	opt.UpdateOptions()
	if jsonData, err = json.Marshal(&opt); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(optFilename, out.Bytes(), 0644)
		}
	}
	return err
}
