package services

import (
	"event-explorer/models"
)

func GetConcurrentEvents(city, countryCode string) (map[string][]models.UIEvent, map[string]error) {
	musicChan := make(chan models.ConcurrentEventPayload)
	sportsChan := make(chan models.ConcurrentEventPayload)

	go func() {
		resp, err := FetchTicketmasterEvents(city, countryCode, "Music", "6")
		if err != nil {
			musicChan <- models.ConcurrentEventPayload{Error: err}
			return
		}
		musicChan <- models.ConcurrentEventPayload{Events: mapTMToUIEvents(resp)}
	}()

	go func() {
		resp, err := FetchTicketmasterEvents(city, countryCode, "Sports", "6")
		if err != nil {
			sportsChan <- models.ConcurrentEventPayload{Error: err}
			return
		}
		sportsChan <- models.ConcurrentEventPayload{Events: mapTMToUIEvents(resp)}
	}()

	musicResult := <-musicChan
	sportsResult := <-sportsChan

	results := make(map[string][]models.UIEvent)
	errors := make(map[string]error)

	if musicResult.Error != nil {
		errors["music"] = musicResult.Error
	} else {
		results["music"] = musicResult.Events
	}

	if sportsResult.Error != nil {
		errors["sports"] = sportsResult.Error
	} else {
		results["sports"] = sportsResult.Events
	}

	return results, errors
}

func mapTMToUIEvents(resp models.TicketmasterResponse) []models.UIEvent {
	var uiEvents []models.UIEvent

	if resp.Embedded.Events == nil {
		return uiEvents
	}

	for _, tmEvent := range resp.Embedded.Events {
		event := models.UIEvent{
			ID:          tmEvent.ID,
			Name:        tmEvent.Name,
			TicketURL:   tmEvent.URL,
			Description: tmEvent.Description,
			Date:        tmEvent.Dates.Start.LocalDate,
			Time:        tmEvent.Dates.Start.LocalTime,
			TimeZone:    tmEvent.Dates.TimeZone,
		}

		if len(tmEvent.Images) > 0 {
			event.Image = tmEvent.Images[0].URL
		}

		if len(tmEvent.Embedded.Venues) > 0 {
			v := tmEvent.Embedded.Venues[0]
			event.Venue = v.Name
			event.City = v.City.Name
			event.State = v.State.Name
			event.Country = v.Country.Name
			event.Address = v.Address.Adress
		}

		uiEvents = append(uiEvents, event)
	}

	return uiEvents
}
