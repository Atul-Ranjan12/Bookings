package repository

import (
	"github/Atul-Ranjan12/booking/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error)
}