package appcontext_test

import (
	"testing"

	"github.com/cvartan/appcontext"
)

type SimpleObject struct {
	App appcontext.ApplicationContext `component:""`
}

func TestAppcontextCreate(t *testing.T) {
	app := appcontext.Application
	if app == nil {
		t.Error("application context is not defined")
	}

	o := &SimpleObject{}
	app.Components().Bind(o)
	if o.App == nil {
		t.Error("application is not filed")
	}
}
