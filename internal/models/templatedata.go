package models

import "github/Atul-Ranjan12/booking/internal/forms"

// Template Data holds data set from templates to handlers
type TemplateData struct {
	StringMap map[string]string
	IntMap map[string]int 
	FloatMap map[string]float32
	Data map[string]interface{}
	CSRFToken string 
	Flash string
	Warning string
	Error string
	Form *forms.Form
}