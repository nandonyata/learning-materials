package scalars

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// // MarshalDate serializes time.Time to string (YYYY-MM-DD format)
// func MarshalDate(t time.Time) graphql.Marshaler {
// 	if t.IsZero() {
// 		return graphql.Null
// 	}

// 	return graphql.WriterFunc(func(w io.Writer) {
// 		date := t.Format("2006-01-02")
// 		io.WriteString(w, strconv.Quote(date))
// 	})
// }

// // UnmarshalDate deserializes string to time.Time
// func UnmarshalDate(v interface{}) (time.Time, error) {
// 	if tmpStr, ok := v.(string); ok {
// 		return time.Parse("2006-01-02", tmpStr)
// 	}
// 	return time.Time{}, fmt.Errorf("date must be a string in YYYY-MM-DD format")
// }

// Date is a custom scalar type for dates
type Date struct {
	time.Time
}

// NewDate creates a Date from time.Time
func NewDate(t time.Time) Date {
	return Date{Time: t}
}

// NewDatePtr creates a *Date from *time.Time
func NewDatePtr(t *time.Time) *Date {
	if t == nil {
		return nil
	}
	date := Date{Time: *t}
	return &date
}

// ToTime converts Date to time.Time
func (d Date) ToTime() time.Time {
	return d.Time
}

// ToTimePtr converts *Date to *time.Time
func ToTimePtr(d *Date) *time.Time {
	if d == nil {
		return nil
	}
	t := d.Time
	return &t
}

// MarshalDate serializes Date to string (YYYY-MM-DD format)
func MarshalDate(t Date) graphql.Marshaler {
	if t.IsZero() {
		return graphql.Null
	}

	return graphql.WriterFunc(func(w io.Writer) {
		date := t.Format("2006-01-02")
		io.WriteString(w, strconv.Quote(date))
	})
}

// UnmarshalDate deserializes string to Date
func UnmarshalDate(v interface{}) (Date, error) {
	if tmpStr, ok := v.(string); ok {
		t, err := time.Parse("2006-01-02", tmpStr)
		if err != nil {
			return Date{}, err
		}
		return Date{Time: t}, nil
	}
	return Date{}, fmt.Errorf("date must be a string in YYYY-MM-DD format")
}
