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
		end_date, room_id, created_at, updated_at, processed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 0) RETURNING id`

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

// Returns a slice of Reservations from the database
func (m *PostgresDBRepo) AllReservations(showNew bool) ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	var query string

	if showNew {
		query = `
		SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name, r.processed
		FROM reservations r
		LEFT JOIN rooms rm ON (r.room_id =  rm.id)
		WHERE r.processed = 0
		ORDER BY r.start_date asc
	`
	} else {
		query = `
		SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name, r.processed
		FROM reservations r
		LEFT JOIN rooms rm ON (r.room_id =  rm.id)
		ORDER BY r.start_date asc
	`
	}

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println("Cannot execute this query")
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
			&i.Processed,
		)
		if err != nil {
			log.Println("error scanning the rows into variables")
			return reservations, err
		}

		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// Get one Reservatinn from the database
func (m *PostgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `
		SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
		FROM reservations r 
		LEFT JOIN rooms rm ON (r.room_id = rm.id)
		WHERE r.id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)

	if err != nil {
		log.Println("Counld not find the row in the databse")
		return res, err
	}

	return res, nil
}

// Updates a reservation in the database
func (m *PostgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE reservations
		SET first_name = $1, last_name = $2, email = $3, phone=$4, updated_at=$5
		WHERE id = $6
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM reservations WHERE id=$1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// Updates Processed for Resrvation by ID
func (m *PostgresDBRepo) UpdateProcessedReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE reservations
		SET processed = $1
		WHERE id = $2
	`

	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}
	return nil
}
