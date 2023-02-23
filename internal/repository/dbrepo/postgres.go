package dbrepo

import (
	"context"
	"github/Atul-Ranjan12/booking/internal/models"
	"log"
	"time"
)

func (m *PostgresDBRepo) AllUsers() bool {
	return true
}

// Inserts a reservation in the database
func (m *PostgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// Take care of timeouts in browsers
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Write query
	var newID int

	stmt := `INSERT INTO reservations (first_name, last_name, email, phone, start_date, 
		end_date, room_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	err := m.DB.QueryRowContext(ctx,
		stmt,
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

// Inserts a room Restriction into the database
func (m *PostgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	// Take care of timeouts in browsers
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO room_restrictions (start_date, end_date, room_id, reservation_id,
		restriction_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.QueryContext(ctx,
		stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		log.Println("Error while trying to execute the query")
		return err
	}

	return nil
}

// Returns true if availability exists for roomID :: else it returns false
func (m *PostgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT 
			COUNT(id)
		FROM 
			room_restrictions
		WHERE
			room_id = $1 and
			$2 < end_date and $3 > start_date
	`
	var numRows int

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		log.Println("ERROR: Error while quering for reservation")
		return false, err
	}

	log.Println("Value of numRows is: ", numRows)

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// Returns a slice of available rooms if any for given date range
func (m *PostgresDBRepo) SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `
	SELECT 
		r.id, r.room_name
	FROM 
		rooms r
	WHERE
		r.id NOT IN (
			SELECT 
				rr.room_id
			FROM 
				room_restrictions rr
			WHERE 
				$1 < rr.end_date and $2 > rr.start_date
		)
	`
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		log.Println("Error while executing the query")
		return rooms, err
	}

	// Append to the rooms slice
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

	if err := rows.Err(); err != nil {
		log.Println("Error while reading and storing values into rooms")
		return rooms, err
	}

	return rooms, err
}

// Get Room information from the room ID
func (m *PostgresDBRepo) GetRoomByID(id int) (models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, room_name, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`
	var room models.Room

	row := m.DB.QueryRowContext(ctx, query, id)
	if err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt); err != nil {
		log.Println("Error scanning the room details")
		return room, err
	}

	return room, nil
}
