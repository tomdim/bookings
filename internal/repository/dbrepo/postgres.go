package dbrepo

import (
	"context"
	"github.com/tomdim/bookings/internal/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a new reservation to the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	stmt := `
		INSERT INTO reservations
		(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id
	`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a new room restriction to the database
func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO room_restrictions
		(start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := m.DB.ExecContext(ctx, stmt,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		res.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID returns if availability exists for a given room, otherwise false
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	query := `
		SELECT 
			count(id) 
		FROM 
			room_restrictions
		WHERE
		    room_id = $1 and
			$2 < end_date and $3 > start_date;
	`
	err := m.DB.QueryRowContext(ctx, query,
		roomID,
		start,
		end,
	).Scan(&numRows)
	if err != nil {
		return false, err
	}

	return numRows == 0, nil
}

// SearchAvailabilityForAllRooms returns the available rooms by dates
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `
		SELECT 
			r.id, r.room_name
		FROM 
			rooms r
		WHERE
		    r.id not in (
		        SELECT 
		            rr.room_id 
		        FROM
		        	room_restrictions rr
		        WHERE 
		        	$1 < end_date and $2 > start_date
		    )
	`
	rows, err := m.DB.QueryContext(ctx, query,
		start,
		end,
	)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// GetRoomByID returns a room based on its ID
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	query := `
		SELECT 
			id, room_name, created_at, updated_at
		FROM 
			rooms
		WHERE
		    id = $1;
	`
	row := m.DB.QueryRowContext(ctx, query,
		id,
	)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}

	return room, nil
}
