package main

import (
	"github/Atul-Ranjan12/booking/internal/config"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T){
	var app config.AppConfig

	mux := routes(&app)
	switch v := mux.(type){
	case *chi.Mux:
		//test passed
	default:
		t.Errorf("Type is not chi.Mux but it is of type: %T", v)
	} 
}