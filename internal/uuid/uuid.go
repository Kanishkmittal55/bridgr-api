package uuid

import (
	"errors"
	"github.com/gofrs/uuid/v5"
	guuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	oapi "github.com/oapi-codegen/runtime/types"
)

// TODO: Use github.com/google/uuid instead.
// https://hassleskip.atlassian.net/browse/DPLAT-3534

func NewDBUuid() (uuid.UUID, error) {
	return uuid.NewV7()
}

func NewGUuid() guuid.UUID {
	return guuid.New()
}

func NewEventUuid() (uuid.UUID, error) {
	return uuid.NewV7()
}

func FromString(text string) (uuid.UUID, error) {
	return uuid.FromString(text)
}

func FromOpenapiUUID(oapiUuid oapi.UUID) (uuid.UUID, error) {
	return uuid.FromString(oapiUuid.String())
}

func ToOpenApiUUID(uuid uuid.UUID) oapi.UUID {
	return oapi.UUID(uuid)
}

// ConvertPgUUIDToOapiUUID converts pgtype.UUID to oapi-codegen UUID
func ConvertPgUUIDToOapiUUID(pgUUID pgtype.UUID) (oapi.UUID, error) {
	if !pgUUID.Valid {
		return oapi.UUID{}, errors.New("invalid UUID")
	}
	return oapi.UUID(pgUUID.Bytes), nil
}

// ToString converts a UUID byte array to a string.
func ToString(bytes [16]byte) (string, error) {
	uuidFromBytes, err := uuid.FromBytes(bytes[:])
	if err != nil {
		return "", err
	}
	return uuidFromBytes.String(), nil
}

// ToPgUuid casts a uuid.UUID to pgtype.UUID. If the UUID supplied if empty (i.e. nil)
// an empty pgtype.UUID instance is returned.
func ToPgUuid(uuid uuid.UUID) pgtype.UUID {
	if uuid.IsNil() {
		return pgtype.UUID{}
	}

	return pgtype.UUID{
		Bytes: [16]byte(uuid.Bytes()),
		Valid: true,
	}
}

// ConvertOapiUUIDToPgUUID converts oapi-codegen UUID to pgtype.UUID
func ConvertOapiUUIDToPgUUID(oapiUUID oapi.UUID) (pgtype.UUID, error) {
	uuidObj, err := FromOpenapiUUID(oapiUUID)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return ToPgUuid(uuidObj), nil
}
