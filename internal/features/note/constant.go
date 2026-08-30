package note

const (
	MoodHappy    = "happy"
	MoodSad      = "sad"
	MoodExcited  = "excited"
	MoodCalm     = "calm"
	MoodAngry    = "angry"
	MoodTired    = "tired"
	MoodGrateful = "grateful"
	MoodNeutral  = "neutral"
)

var ValidMoods = map[string]struct{}{
	MoodHappy:    {},
	MoodSad:      {},
	MoodExcited:  {},
	MoodCalm:     {},
	MoodAngry:    {},
	MoodTired:    {},
	MoodGrateful: {},
	MoodNeutral:  {},
}
