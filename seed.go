package main

import (
	"time"

	"github.com/jmoiron/sqlx"
)

func seed() {
	log.Info("seeding...")
	conn, err := sqlx.Connect("sqlite3", "barkometer.db")
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := conn.Prepare(`INSERT INTO events (category, notes, recorded_at) VALUES ($1, $2, $3)`)
	if err != nil {
		log.Fatal(err)
	}
	var t time.Time
	// "2006-01-02T15:04:05Z07:00"
	t, err = time.Parse(time.RFC3339, "2019-04-27T19:08:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-04-28T00:45:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("minor", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-04-28T12:07:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("minor", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-04-28T15:46:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("minor", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-04-29T08:23:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-04-29T08:40:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-04-29T11:30:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("minor", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-04-29T20:40:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-04-30T08:00:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-04-30T08:45:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-04-30T21:08:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-04-30T21:40:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-01T11:40:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("minor", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-03T07:20:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-04T08:22:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("minor", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-05T09:02:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("minor", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-05T09:20:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-05T13:44:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-05T23:05:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-06T07:00:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-07T21:30:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-08T07:00:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("extreme", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-08T21:45:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("minor", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-09T08:18:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-09T21:33:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-10T07:00:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-10T21:40:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-10T22:09:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("minor", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-10T22:30:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("minor", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-11T07:00:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("extreme", "constant whining 1 hour unprecedented", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-11T08:15:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-11T23:40:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-12T09:40:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-12T17:09:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-13T06:45:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("extreme", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	t, err = time.Parse(time.RFC3339, "2019-05-13T07:30:00+08:00")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec("major", "", t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
	log.Info("seeding complete")
}
