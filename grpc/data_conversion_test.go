package grpc_test

import (
	"testing"
	"time"

	commonv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/common/v1"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/a-novel/golib/database"
	"github.com/a-novel/golib/grpc"
)

func TestTimestampOptional(t *testing.T) {
	timestamp := grpc.TimestampOptional(lo.ToPtr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
	require.NotNil(t, timestamp)
	require.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), timestamp.AsTime())

	timestamp = grpc.TimestampOptional(nil)
	require.Nil(t, timestamp)
}

func TestTimestampOptionalProto(t *testing.T) {
	timestamp := grpc.TimestampOptionalProto(timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
	require.NotNil(t, timestamp)
	require.Equal(t, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), *timestamp)

	timestamp = grpc.TimestampOptionalProto(nil)
	require.Nil(t, timestamp)
}

func TestDurationOptional(t *testing.T) {
	duration := grpc.DurationOptional(lo.ToPtr(5 * time.Second))
	require.NotNil(t, duration)
	require.Equal(t, 5*time.Second, duration.AsDuration())

	duration = grpc.DurationOptional(nil)
	require.Nil(t, duration)
}

func TestDurationOptionalProto(t *testing.T) {
	duration := grpc.DurationOptionalProto(durationpb.New(5 * time.Second))
	require.NotNil(t, duration)
	require.Equal(t, 5*time.Second, *duration)

	duration = grpc.DurationOptionalProto(nil)
	require.Nil(t, duration)
}

func TestStructOptional(t *testing.T) {
	structure, err := grpc.StructOptional(map[string]interface{}{
		"key": "value",
	})
	require.Nil(t, err)
	require.NotNil(t, structure)
	require.Equal(t, "value", structure.Fields["key"].GetStringValue())

	structure, err = grpc.StructOptional(nil)
	require.Nil(t, err)
	require.Nil(t, structure)
}

func TestStructOptionalProto(t *testing.T) {
	initial, err := grpc.StructOptional(map[string]interface{}{
		"key": "value",
	})
	require.Nil(t, err)

	structure := grpc.StructOptionalProto(initial)
	require.NotNil(t, structure)
	require.Equal(t, "value", structure["key"])

	structure = grpc.StructOptionalProto(nil)
	require.Nil(t, structure)
}

func TestProtoConverterFromProto(t *testing.T) {
	testCases := []struct {
		name string

		source string

		mapper        grpc.ProtoMapper[string, int]
		protoDefault  string
		entityDefault int

		expect int
	}{
		{
			name: "OK",

			source: "one",

			mapper: grpc.ProtoMapper[string, int]{
				"one":   1,
				"two":   2,
				"three": 3,
			},

			protoDefault:  "zero",
			entityDefault: 0,

			expect: 1,
		},
		{
			name: "Default",

			source: "four",

			mapper: grpc.ProtoMapper[string, int]{
				"one":   1,
				"two":   2,
				"three": 3,
			},

			protoDefault:  "zero",
			entityDefault: 0,

			expect: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			converter := grpc.NewProtoConverter(testCase.mapper, testCase.protoDefault, testCase.entityDefault)
			require.Equal(t, testCase.expect, converter.FromProto(testCase.source))
		})
	}
}

func TestProtoConverterToProto(t *testing.T) {
	testCases := []struct {
		name string

		source int

		mapper        grpc.ProtoMapper[string, int]
		protoDefault  string
		entityDefault int

		expect string
	}{
		{
			name: "OK",

			source: 1,

			mapper: grpc.ProtoMapper[string, int]{
				"one":   1,
				"two":   2,
				"three": 3,
			},

			protoDefault:  "zero",
			entityDefault: 0,

			expect: "one",
		},
		{
			name: "Default",

			source: 4,

			mapper: grpc.ProtoMapper[string, int]{
				"one":   1,
				"two":   2,
				"three": 3,
			},

			protoDefault:  "zero",
			entityDefault: 0,

			expect: "zero",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			converter := grpc.NewProtoConverter(testCase.mapper, testCase.protoDefault, testCase.entityDefault)
			require.Equal(t, testCase.expect, converter.ToProto(testCase.source))
		})
	}
}

func TestSortDirectionConverterFromProto(t *testing.T) {
	converter := grpc.SortDirectionConverter
	require.Equal(t, database.SortDirectionAsc, converter.FromProto(commonv1.SortDirection_SORT_DIRECTION_ASC))
	require.Equal(t, database.SortDirectionDesc, converter.FromProto(commonv1.SortDirection_SORT_DIRECTION_DESC))
	require.Equal(t, database.SortDirectionNone, converter.FromProto(commonv1.SortDirection_SORT_DIRECTION_UNSPECIFIED))
}
