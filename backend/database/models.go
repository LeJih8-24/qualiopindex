package database

// import "gorm.io/gorm"
import "sync"


type Format struct {
    ID string `json:"id"`
    Client string `json:"client"`
    Description string `json:"description"`
    Formation string `json:"formation"`
    Date string `json:"date"`
    Formateur string `json:"formateur"`
    Financer string `json:"financer"`
    Durer int `json:"durer"`
    Satis int `json:"satis"`
    SatisLater int `json:"satisLater"`
}

type Day struct {
	ID string `json:"id"`
	FormatID string `json:"format_id"`
	Which int `json:"which"`
	Index int `json:"index"`
}

type Present struct {
	ID string `json:"id"`
	DayID string `json:"day_id"`
	Present string `json:"present"`
}

type Store struct {
	Mu sync.RWMutex
	Days map[string]Day
	Presents map[string]Present
}

var Planner = &Store{
	Days: map[string]Day{},
	Presents: map[string]Present{},
}
