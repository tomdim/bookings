package dbrepo

import (
	"errors"
	"github.com/tomdim/bookings/internal/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a new reservation to the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail, otherwise pass
	if res.RoomID == 2 {
		return 0, errors.New("test error")
	}

	return 1, nil
}

// InsertRoomRestriction inserts a new room restriction to the database
func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	// if the room id is 1000, then fail, otherwise pass
	if res.RoomID == 1000 {
		return errors.New("test error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns if availability exists for a given room, otherwise false
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	// if the room id is 1000, then return error
	if roomID == 1000 {
		return false, errors.New("test error")
	}
	// if the room id is 2, then return true
	if roomID == 2 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns the available rooms by dates
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	startDate := start.Format("2006-01-02")
	if startDate == "2050-01-01" {
		room := models.Room{
			ID:        1,
			RoomName:  "test room",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		rooms = append(rooms, room)
	} else if startDate == "2000-01-01" {
		return rooms, errors.New("test error")
	}

	return rooms, nil
}

// GetRoomByID returns a room based on its ID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id == 2 {
		return models.Room{
			ID:        2,
			RoomName:  "General's Quarters",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	} else if id >= 3 {
		return room, errors.New("test error")
	}
	return room, nil
}
