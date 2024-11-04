package grpc

import (
	"time"

	commonv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/common/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/a-novel/golib/database"
)

func TimestampOptional(src *time.Time) *timestamppb.Timestamp {
	if src == nil {
		return nil
	}

	return timestamppb.New(*src)
}

type ProtoMapper[Proto comparable, Entity comparable] map[Proto]Entity

type ProtoConverter[Proto comparable, Entity comparable] interface {
	ToProto(src Entity) Proto
	FromProto(src Proto) Entity
}

type protoConverterImpl[Proto comparable, Entity comparable] struct {
	mapper        ProtoMapper[Proto, Entity]
	protoDefault  Proto
	entityDefault Entity
}

func (c *protoConverterImpl[Proto, Entity]) ToProto(src Entity) Proto {
	for proto, entity := range c.mapper {
		if entity == src {
			return proto
		}
	}

	return c.protoDefault
}

func (c *protoConverterImpl[Proto, Entity]) FromProto(src Proto) Entity {
	if entity, ok := c.mapper[src]; ok {
		return entity
	}

	return c.entityDefault
}

func NewProtoConverter[Proto comparable, Entity comparable](
	mapper ProtoMapper[Proto, Entity],
	protoDefault Proto,
	entityDefault Entity,
) ProtoConverter[Proto, Entity] {
	return &protoConverterImpl[Proto, Entity]{
		mapper:        mapper,
		protoDefault:  protoDefault,
		entityDefault: entityDefault,
	}
}

var SortDirectionConverter = NewProtoConverter[commonv1.SortDirection, database.SortDirection](
	ProtoMapper[commonv1.SortDirection, database.SortDirection]{
		commonv1.SortDirection_SORT_DIRECTION_ASC:  database.SortDirectionAsc,
		commonv1.SortDirection_SORT_DIRECTION_DESC: database.SortDirectionDesc,
	},
	commonv1.SortDirection_SORT_DIRECTION_UNSPECIFIED,
	database.SortDirectionNone,
)
