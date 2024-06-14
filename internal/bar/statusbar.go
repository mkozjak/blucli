package bar

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/mkozjak/blutui/internal/app"
	"github.com/mkozjak/blutui/internal/library"
	"github.com/mkozjak/blutui/internal/player"
	"github.com/mkozjak/tview"
)

type Command interface {
	SetCurrentPage(name string)
}

// injection target
type StatusBar struct {
	container *tview.Table
	app       app.Command
	library   library.Command
}

func newStatusBar(a app.Command, l library.Command) *StatusBar {
	return &StatusBar{
		app:     a,
		library: l,
	}
}

func (sb *StatusBar) createContainer() (*tview.Table, error) {
	sb.container = tview.NewTable().
		SetFixed(1, 3).
		SetSelectable(false, false).
		SetCell(0, 0, tview.NewTableCell("processing").
			SetExpansion(1).
			SetTextColor(tcell.ColorDefault).
			SetAlign(tview.AlignLeft)).
		SetCell(0, 1, tview.NewTableCell("welcome to blutui =)").
			SetExpansion(2).
			SetTextColor(tcell.ColorDefault).
			SetAlign(tview.AlignCenter).
			SetMaxWidth(40)).
		SetCell(0, 2, tview.NewTableCell("").
			SetExpansion(1).
			SetTextColor(tcell.ColorDefault).
			SetAlign(tview.AlignRight))

	sb.container.SetBackgroundColor(tcell.ColorDefault).SetBorder(false).SetBorderPadding(0, 0, 1, 1)

	return sb.container, nil
}

func (sb *StatusBar) listen(ch <-chan player.Status) {
	for s := range ch {
		var cpTitle string
		var cpFormat string
		var cpQuality string

		switch s.State {
		case "play":
			s.State = "playing"
			cpTitle = s.Artist + " - " + s.Track
			cpFormat = s.Format
			cpQuality = s.Quality

			// TODO: should probably be done elsewhere
			if sb.app.GetCurrentPage() == "library" {
				if s.Service == "LocalMusic" {
					sb.library.HighlightCpArtist(s.Artist)
					sb.library.SetCpTrackName(s.Track)
				} else {
					sb.library.HighlightCpArtist("")
					sb.library.SetCpTrackName("")
				}
			}
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
			sb.library.HighlightCpArtist("")
			sb.library.SetCpTrackName("")
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

		currPage := sb.app.GetCurrentPage()
		format := ""
		if cpQuality != "" || cpFormat != "" {
			format = " | " + cpQuality + " " + cpFormat
		}

		sb.container.GetCell(0, 0).SetText("vol: " + strconv.Itoa(s.Volume) +
			" | " + s.State + format)
		sb.container.GetCell(0, 1).SetText(cpTitle)
		sb.container.GetCell(0, 2).SetText(currPage)
		sb.app.Draw()
	}
}

func (sb *StatusBar) setCurrentPage(name string) {
	sb.container.GetCell(0, 2).SetText(name)
}