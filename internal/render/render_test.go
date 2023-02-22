package render

import (
	"github/Atul-Ranjan12/booking/internal/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Fatal(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("Value of the flash was not the same")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error("Error creating template cache")
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Fatal(err)
	}

	var ww myTestWriter

	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error(err)
	}

	err = Template(&ww, r, "non-existant.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered template that it was not supposedd to")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
