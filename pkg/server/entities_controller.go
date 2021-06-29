package server

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/profiralex/go-bootstrap-redis/pkg/bl"
	"net/http"
)

// swagger:model
type entityResponse struct {
	UUID   string `json:"uuid"`
	Field1 string `json:"field_1"`
	Field2 int    `json:"field_2"`
	Field3 bool   `json:"field_3"`
	Field4 string `json:"field_4"`
}

func convertToEntityResponse(e bl.Entity) entityResponse {
	return entityResponse{
		UUID:   e.UUID,
		Field1: e.Field1,
		Field2: e.Field2,
		Field3: e.Field3,
		Field4: e.Field4,
	}
}

type entitiesController struct {
	entitiesRepo bl.EntitiesRepository
}

func newEntitiesController() *entitiesController {
	return &entitiesController{
		entitiesRepo: &bl.RedisEntitiesRepository{},
	}
}

// swagger:model
type createEntityRequest struct {
	Field1 string `json:"field_1" valid:"required"`
	Field2 int    `json:"field_2"`
	Field3 bool   `json:"field_3"`
	Field4 string `json:"field_4" valid:"required"`
}

func (d *createEntityRequest) Bind(*http.Request) error {
	_, err := govalidator.ValidateStruct(d)
	return err
}

// swagger:operation POST /entities createEntity
// Create Entity
//
// Create an Entity
// ---
// parameters:
// - { name: data, in: body, schema: { "$ref": "#/definitions/createEntityRequest" }, required: true, description: entity creation params}
// responses:
//  201: { schema: { "$ref": "#/definitions/entityResponse" } }
func (c *entitiesController) createEntity(w http.ResponseWriter, r *http.Request) {
	data := &createEntityRequest{}
	if err := render.Bind(r, data); err != nil {
		respondValidationErrors(w, r, err, http.StatusBadRequest)
		return
	}

	entity := bl.Entity{
		Field1: data.Field1,
		Field2: data.Field2,
		Field3: data.Field3,
		Field4: data.Field4,
	}
	err := c.entitiesRepo.Save(r.Context(), &entity)

	if err != nil {
		respondError(w, r, fmt.Errorf("failed to save entity: %w", err))
		return
	}

	respondSuccess(w, r, convertToEntityResponse(entity), http.StatusCreated)
}

// swagger:operation GET /entities/{uuid} getEntity
// Get Entity by uuid
//
// Get an Entity by uuid
// ---
// parameters:
// - { name: uuid, in: path, type: string, required: true, description: entity uuid}
// responses:
//  201: { schema: { "$ref": "#/definitions/entityResponse" } }
func (c *entitiesController) getEntity(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	entity, err := c.entitiesRepo.FindByUUID(r.Context(), uuid)

	if err != nil {
		respondError(w, r, fmt.Errorf("failed to get entity: %w", err))
		return
	}

	respondSuccess(w, r, convertToEntityResponse(entity), http.StatusOK)
}
