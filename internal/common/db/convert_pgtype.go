package db

import "github.com/jackc/pgx/v5/pgtype"

func ToPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}

	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}

func ToPgInt8(value int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64: value,
		Valid: true,
	}
}
