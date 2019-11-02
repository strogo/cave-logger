package gui

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/idlephysicist/cave-logger/internal/model"
)

type trips struct {
	*tview.Table
	trips			 chan *model.Log
	filterWord string
}

func newTrips(g *Gui) *trips {
	trips := &trips{
		Table: tview.NewTable().SetSelectable(true, false).Select(0,0).SetFixed(1,1),
		trips: make(chan *model.Log),
	}

	trips.SetTitle(` Logs `).SetTitleAlign(tview.AlignLeft)
	trips.SetBorder(true)
	trips.setEntries(g)
	trips.setKeybinding(g)
	return trips
}

func (t *trips) name() string {
	return `trips`
}

func (t *trips) setKeybinding(g *Gui) {
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)


		return event
	})
}

func (t *trips) entries(g *Gui) {
	trips, err := g.db.GetAllLogs()
	if err != nil {
		return
	}

	g.state.resources.trips = trips	
}

func (t *trips) setEntries(g *Gui) {
	t.entries(g)
	table := t.Clear()

	headers := []string{
		"Date",
		"Cave",
		"Names",
		"Notes",
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold,
		})
	}

	for i, trip := range g.state.resources.trips {
		table.SetCell(i+1, 0, tview.NewTableCell(trip.Date).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(30).
			SetExpansion(0))

		table.SetCell(i+1, 1, tview.NewTableCell(trip.Cave).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(30).
			SetExpansion(0))

		table.SetCell(i+1, 2, tview.NewTableCell(trip.Names).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(0).
			SetExpansion(2))

		table.SetCell(i+1, 3, tview.NewTableCell(trip.Notes).
			SetTextColor(tcell.ColorLightGreen).
			SetMaxWidth(0).
			SetExpansion(1))
	}
}

func (t *trips) updateEntries(g *Gui) {}

func (t *trips) focus(g *Gui) {
	t.SetSelectable(true, false)
	g.app.SetFocus(t)
}

func (t *trips) unfocus() {
	t.SetSelectable(false, false)
}

func (t *trips) setFilterWord(word string) {
	t.filterWord = word
}

func (t *trips) monitoringTrips(g *Gui) {
	ticker := time.NewTicker(5 * time.Second)

LOOP:
	for {
		select {
		case <-ticker.C:
			t.updateEntries(g)
		case <-g.state.stopChans["trips"]:
			ticker.Stop()
			break LOOP
		}
	}
}