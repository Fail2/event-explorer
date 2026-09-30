package models

type UIEvent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	TimeZone    string `json:"timeZone"`
	Venue       string `json:"venue"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Address     string `json:"address"`
	Description string `json:"description"`
	TicketURL   string `json:"ticketUrl"`
	Category    string `json:"category"`
	Genre       string `json:"genre"`
	GeneralRule string `json:"genralRule"`
	ChildRule   string `json:"childRule"`
}

type ConcurrentEventPayload struct {
	Events []UIEvent
	Error  error
}

type TMImage struct {
	URL string `json:"url"`
}

type TMDateStart struct {
	LocalDate string `json:"localDate"`
	LocalTime string `json:"localTime"`
}

type TMDates struct {
	Start    TMDateStart `json:"start"`
	TimeZone string      `json:"timeZone"`
}

type TMCategory struct {
	Name string `json:"name"`
}

type TMGenre struct {
	Name string `json:"name"`
}

type TMCity struct {
	Name string `json:"name"`
}

type TMState struct {
	Name string `json:"name"`
}

type TMCountry struct {
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
}

type TMAddress struct {
	Adress string `json:"line1"`
}

type TMGeneralInfo struct {
	GeneralRule string `json:"generalRule"`
	ChildRule   string `json:"childRule"`
}

type TMClassification struct {
	Category TMCategory `json:"segment"`
	Genre    TMGenre    `json:"genre"`
}

type TMVenue struct {
	Name        string        `json:"name"`
	City        TMCity        `json:"city"`
	State       TMState       `json:"state"`
	Country     TMCountry     `json:"country"`
	Address     TMAddress     `json:"address"`
	GeneralInfo TMGeneralInfo `json:"generalInfo"`
}

type TMEventEmbedded struct {
	Venues []TMVenue `json:"venues"`
}

type TMEvent struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	URL             string             `json:"url"`
	Description     string             `json:"pleaseNote"`
	Images          []TMImage          `json:"images"`
	Dates           TMDates            `json:"dates"`
	Classifications []TMClassification `json:"classifications"`
	Embedded        TMEventEmbedded    `json:"_embedded"`
}

type TMResponseEmbedded struct {
	Events []TMEvent `json:"events"`
}

type TicketmasterResponse struct {
	Embedded TMResponseEmbedded `json:"_embedded"`
}
