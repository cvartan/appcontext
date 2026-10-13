package appcontext

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/cvartan/appcontext/pkg/objectrepo"
	"github.com/cvartan/goconfig"
)

type Component interface {
	Init(ctx context.Context)
	Destroy(ctx context.Context)
}

type Service interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

type ApplicationContext interface {
	Context() context.Context
	Configuration() *goconfig.Configuration
	Configure(*goconfig.Options)
	Run() error
}

type application struct {
	ctx           context.Context
	configuration *goconfig.Configuration
	components    *objectrepo.ObjectRepository
	logger        *log.Logger
}

var Application *application

type Options struct {
	Configuration goconfig.Options
}

func (app *application) Configure(options *goconfig.Options) {
	app.configuration.Init(*options)
}

func (app *application) Context() context.Context {
	return app.ctx
}

func (app *application) Configuration() *goconfig.Configuration {
	return app.configuration
}

func (app *application) Components() *objectrepo.ObjectRepository {
	return app.components
}

func (app *application) start(ctx context.Context) {

}

func (app *application) close(ctx context.Context) {

}

func (app *application) Run() {

	c, cancel := signal.NotifyContext(app.ctx, os.Interrupt)
	defer cancel()

	app.start(c)

	<-c.Done()

	app.close(c)
}

func init() {
	Application = &application{
		ctx:           context.Background(),
		components:    objectrepo.New(),
		configuration: goconfig.NewConfiguration(nil),
	}

	Application.components.Bind(Application)
}
