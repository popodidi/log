package log

// Level defines the Severity Level type
type Level int

// Level Enums
const (
	first Level = iota
	Critical
	Error
	Warn
	Notice
	Info
	Debug
	last
)

var (
	levelName = []string{
		"",
		"Crit",
		"Error",
		"Warn",
		"Note",
		"Info",
		"Debug",
		"",
	}
)

// IsValid returns if the l is valid.
func (l Level) IsValid() bool {
	return l < last && l > first
}

// String return the string description of l.
func (l Level) String() string {
	return levelName[l]
}
