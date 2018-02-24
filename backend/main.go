package backend

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/favclip/ucon"
	"github.com/favclip/ucon/swagger"
	"go.mercari.io/datastore"
	"go.mercari.io/datastore/aedatastore"
	"google.golang.org/appengine"
)

func init() {
	ucon.Middleware(UseAppengineContext)
	ucon.Orthodox()
	ucon.Middleware(swagger.RequestValidator())

	swPlugin := swagger.NewPlugin(&swagger.Options{
		Object: &swagger.Object{
			Info: &swagger.Info{
				Title:   "GCPUG",
				Version: "v1",
			},
			Schemes: []string{"http", "https"},
		},
		DefinitionNameModifier: func(refT reflect.Type, defName string) string {
			if strings.HasSuffix(defName, "JSON") {
				return defName[:len(defName)-4]
			}
			return defName
		},
	})
	ucon.Plugin(swPlugin)

	setupOrderP1(swPlugin)
	setupConnpassAPI(swPlugin)

	ucon.DefaultMux.Prepare()
	http.Handle("/api/", ucon.DefaultMux)
}

// UseAppengineContext is UseAppengineContext
func UseAppengineContext(b *ucon.Bubble) error {
	if b.Context == nil {
		b.Context = appengine.NewContext(b.R)
	} else {
		b.Context = appengine.WithContext(b.Context, b.R)
	}

	return b.Next()
}

// FromContext is Create Datastore Client from Context
func FromContext(ctx context.Context) (datastore.Client, error) {
	return aedatastore.FromContext(ctx)
}
