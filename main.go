package main

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app  = kingpin.New("barkometer", "Track nuisance barks")
	dev  = app.Flag("dev", "Enable dev mode.").Bool()
	port = app.Flag("port", "Set server port").Default("8082").String()
)

var log *zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	log = logger.Sugar()
}

const schema = `
CREATE TABLE IF NOT EXISTS events (
	id	    			INTEGER UNIQUE NOT NULL PRIMARY KEY,
	category			VARCHAR(64) NOT NULL,
	notes				VARCHAR(200) NOT NULL,
	recorded_at			TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	archived			BOOL NOT NULL DEFAULT false,
	archived_at			TIMESTAMP
);
`

const drop = `
DROP TABLE IF EXISTS events;
`

type Record struct {
	ID         string     `db:"id"`
	Category   string     `db:"category"`
	Notes      string     `db:"notes"`
	RecordedAt *time.Time `db:"recorded_at"`
}

func main() {
	app.Parse(os.Args[1:])
	conn, err := sqlx.Connect("sqlite3", "barkometer.db")
	if err != nil {
		log.Fatal(err)
	}
	if *dev {
		log.Info("DEV MODE")
		conn.MustExec(drop)
	}
	conn.MustExec(schema)
	if *dev {
		seed()
	}
	app.Parse(os.Args[1:])

	templateIndex, err := template.New("index").Parse(HTMLIndex)
	if err != nil {
		log.Fatal(err)
	}
	templateNew, err := template.New("new").Parse(HTMLNew)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/delete", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "no id provided", http.StatusInternalServerError)
			return
		}

		_, err := conn.Exec(`UPDATE events SET archived=true, archived_at=$1 WHERE id=$2`, time.Now(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := []*Record{}
		err := conn.Select(&data, `SELECT id, category, notes, recorded_at FROM events WHERE archived=false ORDER BY recorded_at DESC`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		templateIndex.Execute(w, data)
	})
	r.Get("/new", func(w http.ResponseWriter, r *http.Request) {
		templateNew.Execute(w, nil)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		category := r.FormValue("category")
		notes := r.FormValue("notes")
		if category == "" {
			http.Error(w, "no category provided", http.StatusInternalServerError)
			return
		}

		_, err = conn.Exec("INSERT INTO events (category, notes, recorded_at) VALUES ($1, $2, $3)", category, notes, time.Now())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Infow("Running server", "port", *port)
	http.ListenAndServe(":"+*port, r)
}

const HTMLNew = `
<!DOCTYPE html>
<html lang="en">
	<title>Create new</title>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<link rel="stylesheet" href="//fonts.googleapis.com/css?family=Roboto:300,300italic,700,700italic" />
	<link rel="stylesheet" href="//cdn.rawgit.com/necolas/normalize.css/master/normalize.css" />
	<link rel="stylesheet" href="//cdn.rawgit.com/milligram/milligram/master/dist/milligram.min.css" />
	<body>
		<section class="container">
			<h1>Report bark</h1>
			<form method="POST" action="/" id="bark-form">
				<input type="radio" name="category" value="minor" checked /> Minor<br />
				<input type="radio" name="category" value="major" /> Major<br />
				<input type="radio" name="category" value="extreme" /> Extreme<br />
				<textarea name="notes" form="bark-form" placeholder="Enter notes here..."></textarea>

				<button type="submit">Submit</button>
			</form>
		</section>
	</body>
</html>

`

const HTMLIndex = `
<!DOCTYPE html>
<html lang="en">
	<title>Create new</title>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<link rel="stylesheet" href="//fonts.googleapis.com/css?family=Roboto:300,300italic,700,700italic" />
	<link rel="stylesheet" href="//cdn.rawgit.com/necolas/normalize.css/master/normalize.css" />
	<link rel="stylesheet" href="//cdn.rawgit.com/milligram/milligram/master/dist/milligram.min.css" />
	<body>
		<section class="container">
		<h1>Barkometer</h1>
		<input type="button" onclick="location.href='/new';" value="Report new bark" />

		<table>
		<thead>
			<tr>
				<th>Name</th>
				<th>Category</th>
				<th>Notes</th>
				<th>Recorded At</th>
				<th>Actions</th>
			</tr>
			</thead>
			<tbody>
			{{range .}}
			<tr>
			<td>Incident #{{.ID}}</td>
			<td>
				<em>{{.Category}}</em>
			</td>
			<td>{{.Notes}}</td>
			<td>{{.RecordedAt.Format "02 Jan 2006 - Mon - 15:04"}}</td>
			<td>
				<form style="margin-bottom: 0" action="/delete?id={{.ID}}" method="post">
					<input class="button button-outline" type="submit" name="delete" value="Delete" />
				</form>
			</td>
			</tr>
			{{end}}
			</tbody>
		</table>
		</section>
	</body>
</html>
`