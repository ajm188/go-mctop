package memcap

import "log"

type Stats struct {
	keys         map[string]*KeyStats
	commands     map[string]int
	commandCount int
	// keys map[string]int
}

func NewStats() *Stats {
	return &Stats{
		map[string]*KeyStats{},
		map[string]int{},
	}
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
		ks = &KeyStats{}
		s.keys[key] = ks
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
