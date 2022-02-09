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
	return false, nil
}

// SearchAvailabilityForAllRooms returns the available rooms by dates
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID returns a room based on its ID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("test error")
	}
	return room, nil
}
