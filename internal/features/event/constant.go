package event

type EventType string

const (
	EventTypeBirthday    EventType = "birthday"
	EventTypeAnniversary EventType = "anniversary"
	EventTypeFirstMeet   EventType = "first_meet"
	EventTypeFirstDate   EventType = "first_date"
	EventTypeWedding     EventType = "wedding"
	EventTypeTrip        EventType = "trip"
	EventTypeSpecial     EventType = "special"
	EventTypeCustom      EventType = "custom"
)
