package memcap

import (
	"log"
	"sort"
)

type Stats struct {
	keys         map[string]*KeyStats
	keyStats     []*KeyStats
	commands     map[string]int
	commandCount int
	// keys map[string]int
}

func NewStats() *Stats {
	return &Stats{
		keys:     map[string]*KeyStats{},
		commands: map[string]int{},
	}
}

func (s *Stats) Calls() map[string]int {
	return s.commands
}

func (s *Stats) TotalCalls() int {
	return s.commandCount
}

func (s *Stats) KeyStats() []*KeyStats {
	return s.keyStats
}

// func (s *Stats) Add(key string) {
// 	_, ok := s.keys[key]
// 	if !ok {
// 		s.keys[key] = 1
// 		return
// 	}

// 	s.keys[key]++
// }

func (s *Stats) Add(key, command string) {
	ks, ok := s.keys[key]
	if !ok {
		ks = &KeyStats{key: key}
		s.keys[key] = ks
		s.keyStats = append(s.keyStats, ks)
	}

	if !ks.Add(command) {
		return
	}

	_, ok = s.commands[command]
	if !ok {
		s.commands[command] = 0
	}

	s.commands[command]++
	s.commandCount++
}

type KeyStats struct {
	key string

	gets    int
	casgets int

	incrs int
	decrs int

	adds    int
	sets    int
	deletes int

	replacements int
	appends      int
	prepends     int

	cas int

	total int
}

func (ks *KeyStats) Add(command string) bool {
	switch command {
	case "get":
		ks.gets++
	case "gets":
		ks.casgets++
	case "incr":
		ks.incrs++
	case "decr":
		ks.decrs++
	case "add":
		ks.adds++
	case "set":
		ks.sets++
	case "delete":
		ks.deletes++
	case "replace":
		ks.replacements++
	case "append":
		ks.appends++
	case "prepend":
		ks.prepends++
	case "cas":
		ks.cas++
	default:
		log.Printf("W Unknown command %s\n", command)
		return false
	}

	ks.total++

	return true
}

// Who doesn't love some boilerplate?
// Consider writing a generator for this; probably want to sort by any of the fields.

func (ks *KeyStats) Key() string  { return ks.key }
func (ks *KeyStats) Gets() int    { return ks.gets }
func (ks *KeyStats) Sets() int    { return ks.sets }
func (ks *KeyStats) Adds() int    { return ks.adds }
func (ks *KeyStats) Deletes() int { return ks.deletes }

func (ks *KeyStats) TotalCalls() int { return ks.total }

type KeyStatsList []*KeyStats

func (ksl KeyStatsList) Len() int           { return len(ksl) }
func (ksl KeyStatsList) Less(i, j int) bool { return ksl[i].total < ksl[j].total } // should consider breaking ties by key name
func (ksl KeyStatsList) Swap(i, j int)      { ksl[i], ksl[j] = ksl[j], ksl[i] }

type KeyStatsByGet []*KeyStats

func (ksl KeyStatsByGet) Len() int           { return len(ksl) }
func (ksl KeyStatsByGet) Less(i, j int) bool { return ksl[i].gets < ksl[j].gets }
func (ksl KeyStatsByGet) Swap(i, j int)      { ksl[i], ksl[j] = ksl[j], ksl[i] }

var _ sort.Interface = KeyStatsList([]*KeyStats{})
