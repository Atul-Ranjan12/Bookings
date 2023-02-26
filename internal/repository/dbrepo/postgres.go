package dbrepo

import (
	"context"
	"errors"
	"github/Atul-Ranjan12/booking/internal/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
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

// Returns a user by ID
func (m *PostgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, first_name, last_name, email, password, access_level, created_at, updated_at
		FROM users where id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Function to update a user in the database
func (m *PostgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, access_level=$4, updated_at=$5
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

// Function to authenticate a user into the databases
func (m *PostgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "SELECT id, password from users WHERE email = $1", email)

	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	// COmpare password with my password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}
