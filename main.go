package main

import (
	"fmt"
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
		type Data struct {
			Total         int
			MinorTotal    int
			MajorTotal    int
			ExtremeTotal  int
			TimeSinceLast string
			Records       []*Record
		}
		total, minorTotal, majorTotal, extremeTotal, timeSince, err := Stats(conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data := &Data{
			Total:         total,
			MinorTotal:    minorTotal,
			MajorTotal:    majorTotal,
			ExtremeTotal:  extremeTotal,
			TimeSinceLast: fmt.Sprintf("%.0f", timeSince.Hours()),
			Records:       []*Record{},
		}
		err = conn.Select(&data.Records, `SELECT id, category, notes, recorded_at FROM events WHERE archived=false ORDER BY recorded_at DESC`)
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
	log.Fatal(http.ListenAndServe(":"+*port, r))
}

func Stats(conn *sqlx.DB) (int, int, int, int, time.Duration, error) {
	data := []*Record{}
	err := conn.Select(&data, `SELECT id, category, notes, recorded_at FROM events WHERE archived=false ORDER BY recorded_at DESC`)
	if err != nil {
		return 0, 0, 0, 0, time.Duration(0), err
	}

	var total int
	conn.Get(&total, "SELECT COUNT(id) FROM events")
	var minorTotal int
	conn.Get(&minorTotal, "SELECT COUNT(id) FROM events WHERE category='minor'")
	var majorTotal int
	conn.Get(&majorTotal, "SELECT COUNT(id) FROM events WHERE category='major'")
	var extremeTotal int
	conn.Get(&extremeTotal, "SELECT COUNT(id) FROM events WHERE category='extreme'")
	var lastRecorded time.Time
	conn.Get(&lastRecorded, "SELECT recorded_at FROM events order BY recorded_at DESC limit 1")
	timeSinceLast := time.Since(lastRecorded)
	return total, minorTotal, majorTotal, extremeTotal, timeSinceLast, nil

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
			<div style="margin-left:auto; margin-right:auto;">
				<h1 style="text-align:center">Barkometer</h1>
				<input
					style="display:block; margin-left:auto; margin-right:auto;"
					type="button"
					onclick="location.href='/new';"
					value="Report new bark"
				/>
				<p style="text-align:center">Hours since last report: {{.TimeSinceLast}}</p>
				<p style="text-align:center">
<small>Total Reports: {{.Total}}</small><br/>
<small>Minor Reports: {{.MinorTotal}}</small><br/>
<small>Major Reports: {{.MajorTotal}}</small><br/>
<small>Extreme Reports: {{.ExtremeTotal}}</small>
				</p>

			</div>
			<hr />
			<p style="text-align:center">
				<a href="mailto:jtnguyen236@gmail.com?Subject=Barkometer%20Bug%20Report" target="_top">Report bugs here</a>
			</p>
			<hr />
			<div class="container">
				<div class="row">
					<div class="row-wrap">
						{{range .Records}}
						<div class="column">
							<h3>Incident #{{.ID}} - {{.Category}}</h3>
							<p>
								<em>Time: {{.RecordedAt.Format "15:04"}}</em>
								<br />
								<em>Day: {{.RecordedAt.Format "Monday"}}</em>
								<br />
								<em>Date: {{.RecordedAt.Format "02 Jan 2006"}}</em>
							</p>

							<p>{{.Notes}}</p>
							<form style="margin-bottom: 0" action="/delete?id={{.ID}}" method="post">
								<input class="button button-outline" type="submit" name="delete" value="Delete" />
							</form>
						</div>
						<hr />
						{{ end }}
					</div>
				</div>
			</div>
		</section>
	</body>
</html>


`
