package resource

import (
	"net/http"

	"github.com/hectane/hectane/db"
	"github.com/manyminds/api2go"
)

// Create attempts to save a new model instance to the database.
func (r *Resource) Create(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	if r.AllHook != nil {
		if err := r.AllHook(ActionCreate, req); err != nil {
			return nil, err
		}
	}
	if r.SetHook != nil {
		r.SetHook(obj, req)
	}
	if err := db.C.Create(obj).Error; err != nil {
		return nil, err
	}
	return &api2go.Response{
		Res:  obj,
		Code: http.StatusCreated,
	}, nil
}
