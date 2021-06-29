package bl

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/profiralex/go-bootstrap-redis/pkg/db"
)

var _ encoding.BinaryMarshaler = (*Entity)(nil)
var _ encoding.BinaryUnmarshaler = (*Entity)(nil)

type Entity struct {
	UUID   string `json:"uuid"`
	Field1 string `json:"field_1"`
	Field2 int    `json:"field_2"`
	Field3 bool   `json:"field_3"`
	Field4 string `json:"field_4"`
}

func (e Entity) MarshalBinary() (data []byte, err error) {
	return json.Marshal(e)
}

func (e *Entity) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}

type EntitiesRepository interface {
	FindByUUID(ctx context.Context, uuid string) (Entity, error)
	Save(ctx context.Context, entity *Entity) error
}

var _ EntitiesRepository = (*RedisEntitiesRepository)(nil)

type RedisEntitiesRepository struct {
}

func (r *RedisEntitiesRepository) FindByUUID(ctx context.Context, uuid string) (Entity, error) {
	rdb, err := db.GetRedisClientFromContext(ctx)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to get rdb client: %w", err)
	}

	var e Entity
	err = rdb.Get(ctx, uuid).Scan(&e)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to get record: %w", err)
	}

	return e, nil
}

func (r *RedisEntitiesRepository) Save(ctx context.Context, entity *Entity) error {
	rdb, err := db.GetRedisClientFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get rdb client: %w", err)
	}

	if len(entity.UUID) == 0 {
		entity.UUID = uuid.NewString()
	}

	err = rdb.Set(ctx, entity.UUID, entity, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save record: %w", err)
	}

	return err
}
