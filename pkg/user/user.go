package user

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"
)

const Timestamp = time.RFC3339

type User struct {
	ID         int       `storm:"id,increment"`
	TelegramID int64     `json:"telegram_id"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	IsBanned   bool      `json:"is_banned"`
	BanDate    time.Time `json:"ban_date"`
	JoinDate   time.Time `json:"join_date"`
}

func (user User) String() string {
	var buf = &bytes.Buffer{}
	if user.Name != "" {
		fmt.Fprintf(buf, " %s", user.Name)
	}
	if user.Username != "" {
		fmt.Fprintf(buf, "@%s ", user.Username)
	}
	buf.WriteString(user.JoinDate.Format(time.RFC822))
	switch user.IsBanned {
	case true:
		buf.WriteString("is banned ")
	case false:
		buf.WriteString("is alive ")
	}
	return buf.String()
}

func (User) CSVHeaders() []string {
	return []string{
		"id",
		"name",
		"username",
		"is_banned",
		"join_date",
	}
}

func (user User) Values() []interface{} {
	return []interface{}{
		user.ID,
		user.Name,
		user.Username,
		user.IsBanned,
		user.JoinDate,
	}
}

func (user User) CSV() []string {
	return []string{
		strconv.Itoa(user.ID),
		user.Name,
		user.Username,
		strconv.FormatBool(user.IsBanned),
		user.JoinDate.Format(Timestamp),
	}
}

type Users []User

func (users Users) Copy() Users {
	return append(Users{}, users...)
}

func (users Users) Len() int {
	return len(users)
}

func (users Users) New() Users {
	return make(Users, 0, users.Len())
}

func (users Users) Filter(pred func(user User) bool) Users {
	var filtered = users.New()
	for _, user := range users {
		if pred(user) {
			filtered = append(filtered, user)
		}
	}
	return filtered
}

func (users Users) SortedByLess(less func(a, b User) bool) Users {
	var sorted = users.Copy()
	sort.Slice(sorted, func(i, j int) bool {
		return less(sorted[i], sorted[j])
	})
	return sorted
}

func (users Users) SortedByNewest() Users {
	return users.SortedByLess(func(a, b User) bool {
		return a.JoinDate.Before(b.JoinDate)
	})
}

func (users Users) Banned() Users {
	return users.Filter(func(user User) bool {
		return user.IsBanned
	})
}

func (users Users) NotBanned() Users {
	return users.Filter(func(user User) bool {
		return !user.IsBanned
	})
}

func (Users) CSVHeaders() []string {
	return User{}.CSVHeaders()
}

func (users Users) WriteCSV(wr io.Writer) error {
	var csvWr = csv.NewWriter(wr)
	if err := csvWr.Write(users.CSVHeaders()); err != nil {
		return err
	}
	for _, user := range users {
		if err := csvWr.Write(user.CSV()); err != nil {
			return err
		}
	}
	defer csvWr.Flush()
	return nil
}
