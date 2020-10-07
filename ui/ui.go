package ui

import (
	"fmt"
	"log"
	"sort"

	"github.com/ajm188/go-mctop/memcap"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	isInit   = false
	isClosed = false

	callsTable *widgets.Table = nil
	keysTable  *widgets.Table = nil
)

func Init(quit chan bool) error {
	if isInit {
		return nil
	}

	if err := termui.Init(); err != nil {
		return err
	}

	isInit = true
	isClosed = false

	go func() {
		for e := range termui.PollEvents() {
			switch e.ID {
			case "q", "<C-c>":
				quit <- true
			}
		}
	}()

	return nil
}

func Close() {
	if isClosed {
		return
	}

	isInit = false
	isClosed = true

	termui.Close()
}

func UpdateCalls(stats *memcap.Stats) {
	if callsTable == nil {
		callsTable = widgets.NewTable()
		callsTable.Title = "Calls"
		callsTable.Rows = [][]string{
			{"Op", "Count"},
		}

		callsTable.SetRect(0, 0, 30, 15)
	}

	callsTable.Rows = callsTable.Rows[:1] // reset everything but the header row

	calls := memcap.CallsList(stats.Calls())
	sort.Sort(sort.Reverse(calls))

	for _, callStat := range calls {
		callsTable.Rows = append(callsTable.Rows, []string{callStat.Name(), fmt.Sprint(callStat.Count())})
	}
}

func UpdateKeys(stats *memcap.Stats) {
	if keysTable == nil {
		keysTable = widgets.NewTable()
		keysTable.Title = "Keys"
		keysTable.Rows = [][]string{
			{"Key", "Reads", "Writes"},
		}

		keysTable.SetRect(0, 15, 80, 50)
	}

	keysTable.Rows = keysTable.Rows[:1]

	keys := memcap.KeyStatsList(stats.KeyStats())
	sort.Sort(sort.Reverse(keys))

	for i, ks := range keys {
		if i >= 10 {
			break
		}

		keysTable.Rows = append(keysTable.Rows, []string{ks.Key(), fmt.Sprint(ks.Gets()), fmt.Sprint(ks.Writes())})
	}

	log.Printf("I KeysTable %s\n", keysTable.Rows)
}

func Render() {
	termui.Render(callsTable)
	termui.Render(keysTable)
}
