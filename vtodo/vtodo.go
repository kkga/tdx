package vtodo

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

const (
	ToDoPriorityHigh   = 1
	ToDoPriorityMedium = 5
	ToDoPriorityLow    = 6
)

// generateUID returns a random string containing timestamp and hostname
func GenerateUID() string {
	sb := strings.Builder{}

	time := time.Now().UnixNano()

	randStr := func(n int) string {
		rs := rand.NewSource(time)
		r := rand.New(rs)
		var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyz")

		s := make([]rune, n)
		for i := range s {
			s[i] = letters[r.Intn(len(letters))]
		}
		return string(s)
	}

	sb.WriteString(fmt.Sprint(time))
	sb.WriteString(fmt.Sprintf("-%s", randStr(8)))
	if hostname, _ := os.Hostname(); hostname != "" {
		sb.WriteString(fmt.Sprintf("@%s", hostname))
	}

	return sb.String()
}

// encode adds vtodo into a new Calendar and returns a buffer ready for writing
func Encode(vtodo *ical.Component) (*bytes.Buffer, error) {
	cal := ical.NewCalendar()
	// TODO move this data somewhere
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN")
	cal.Children = append(cal.Children, vtodo)

	var buf bytes.Buffer
	if err := ical.NewEncoder(&buf).Encode(cal); err != nil {
		return &buf, err
	}
	return &buf, nil
}
