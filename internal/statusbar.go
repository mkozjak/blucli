package internal

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/mkozjak/tview"
)

// bottom bar - status
func (a *App) CreateStatusBar() {
	a.StatusBar = tview.NewTable().
		SetFixed(1, 3).
		SetSelectable(false, false).
		SetCell(0, 0, tview.NewTableCell("connecting").
			SetExpansion(1).
			SetTextColor(tcell.ColorDefault).
			SetAlign(tview.AlignLeft)).
		SetCell(0, 1, tview.NewTableCell("welcome to blutui =)").
			SetExpansion(2).
			SetTextColor(tcell.ColorDefault).
			SetAlign(tview.AlignCenter)).
		SetCell(0, 2, tview.NewTableCell("").
			SetExpansion(1).
			SetTextColor(tcell.ColorDefault).
			SetAlign(tview.AlignRight))

	a.StatusBar.SetBackgroundColor(tcell.ColorDefault).SetBorder(false).SetBorderPadding(0, 0, 1, 1)

	// channel for receiving player status updates
	a.sbMessages = make(chan Status)

	// start long-polling for updates
	go a.PollStatus()

	// start a goroutine for receiving the updates
	go a.listener()
}

func (a *App) listener() {
	for s := range a.sbMessages {
		var cpTitle string
		var cpFormat string
		var cpQuality string

		switch s.State {
		case "play":
			s.State = "playing"
			cpTitle = s.Artist + " - " + s.Track
			cpFormat = s.Format
			cpQuality = s.Quality
		case "stream":
			s.State = "streaming"
			cpTitle = s.Title2
			cpFormat = s.Format
			cpQuality = s.Quality
		case "stop":
			s.State = "stopped"
			cpTitle = ""
			cpFormat = ""
			cpQuality = ""
		case "pause":
			s.State = "paused"

			if s.Artist == "" && s.Track == "" {
				// streaming, set title to Title3 from /Status
				cpTitle = s.Title3
			} else {
				cpTitle = s.Artist + " - " + s.Track
			}

			cpFormat = s.Format
			cpQuality = s.Quality
		case "neterr":
			s.State = "network error"
			cpTitle = ""
			cpFormat = ""
			cpQuality = ""
		case "ctrlerr":
			s.State = "player control error"
			cpTitle = ""
			cpFormat = ""
			cpQuality = ""
		}

		a.StatusBar.GetCell(0, 0).SetText("vol: " + strconv.Itoa(s.Volume) + " | " + s.State)
		a.StatusBar.GetCell(0, 1).SetText(cpTitle)
		a.StatusBar.GetCell(0, 2).SetText(cpQuality + " " + cpFormat)
		a.Application.Draw()
	}
}