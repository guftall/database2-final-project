package db

import (
	"database/sql"
	"encoding/json"
	msdb "github.com/denisenkom/go-mssqldb"
	"os"

	"log"
)

type Olympic struct {
	Olympic_name    string `json:"name"`
	Olympic_website string `json:"website"`
	Olympic_logo    string `json:"logo"`
	Olympic_year    int    `json:"year"`
}

type Athlete struct {
	Id     msdb.UniqueIdentifier
	Name   string
	Image  sql.NullString
	Gender int8
	Sport  string
}

type AthleteApiModel struct {
	Id     string         `json:"id"`
	Name   string         `json:"name"`
	Image  sql.NullString `json:"image"`
	Gender int8           `json:"gender"`
	Sport  string         `json:"experienced"`
}

func connect() *sql.DB {
	connStr, ok := os.LookupEnv("DB_CONNECTION")
	if !ok {
		connStr = "sqlserver://SA:Rr12345678@192.168.1.9?database=UNI_Database2"
	}

	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		log.Fatalf("can't open sql connection to %s\n", connStr)
	}

	return db
}

func ReadOlympics() []Olympic {

	db := connect()

	rows, err := db.Query("select olympic_name, olympic_website, olympic_logo, olympic_year from En_Olympic")
	if err != nil {
		log.Fatalf("query failed: %s", err)
	}

	olympics := make([]Olympic, 0)

	for rows.Next() {

		var o Olympic
		if err := rows.Scan(&o.Olympic_name, &o.Olympic_website, &o.Olympic_logo, &o.Olympic_year); err != nil {
			log.Fatal(err)
		}

		olympics = append(olympics, o)
	}

	log.Printf("read %d olympic\n", len(olympics))

	return olympics
}

func ReadAthletes(name, year, country string) []*AthleteApiModel {

	db := connect()

	query := "select En_Athlete.athlete_id, athlete_name, athlete_image, athlete_gender, sport from En_Athlete"
	query += " join Experienced_Athlete_Sport on (En_Athlete.athlete_id = Experienced_Athlete_Sport.athlete_id)"
	query += " where athlete_name like '%" + name + "%'"

	if year != "" {
		query += " and " + year + " = (select olympic_year from En_Olympic" +
			" where olympic_name = (select has_a_olympic from En_Tournament" +
			" where tournament_id = (select tournament_id from Attended_Athlete_Tournament" +
			" where Attended_Athlete_Tournament.athlete_id = En_Athlete.athlete_id)))"
	}

	if country != "" {
		query += " and '" + country + "' = (select from_country from En_Team" +
			" where team_id = En_Athlete.athlete_team_id)"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("query failed: %s", err)
	}

	athletes := make([]*AthleteApiModel, 0)

	for rows.Next() {

		var o Athlete
		if err := rows.Scan(&o.Id, &o.Name, &o.Image, &o.Gender, &o.Sport); err != nil {
			log.Fatal(err)
		}
		log.Printf("id00: %v", o.Id)
		// uidStr := o.Id.String()
		// uid, _ := uuidLib.Parse(uidStr)
		// log.Printf("id01: %s", uid)
		athletes = append(athletes, mapAthlete(&o))
	}

	athletesMap := make(map[string]*AthleteApiModel)

	for _, a := range athletes {
		_, prs := athletesMap[a.Id]
		if prs {
			continue
		}

		athletesMap[a.Id] = a
	}

	log.Printf("read %d athlete\n", len(athletes))

	result := make([]*AthleteApiModel, 0)
	for _, val := range athletesMap {
		result = append(result, val)
	}

	return result
}

func ReadAthlete(id string) *AthleteApiModel {

	db := connect()

	query := "select athlete_id, athlete_name, athlete_image, athlete_gender from En_Athlete "
	query += "where athlete_id='" + id + "'"
	println(query)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("query failed: %s", err)
	}

	if rows.Next() {

		var athlete Athlete
		if err := rows.Scan(&athlete.Id, &athlete.Name, &athlete.Image, &athlete.Gender); err != nil {
			log.Fatal(err)
		}

		return mapAthlete(&athlete)
	}

	return nil
}

func mapAthlete(athlete *Athlete) *AthleteApiModel {
	return &AthleteApiModel{
		Id:     athlete.Id.String(),
		Name:   athlete.Name,
		Image:  athlete.Image,
		Gender: athlete.Gender,
		Sport:  athlete.Sport,
	}
}

func DeleteAthletes(id []byte) {

	type s struct {
		Id []byte `json:"id"`
	}

	var me s
	err := json.Unmarshal(id, &me)
	if err != nil {
		log.Fatalf("error unmarshal: %s\n", err)
	}

	db := connect()

	log.Printf("return id is:%v\n", me.Id)

	var uid msdb.UniqueIdentifier
	e := uid.Scan(me.Id)
	if e != nil {
		panic(e)
	}

	log.Println("uid: ", uid)

	res, err := db.Exec("delete from En_Athlete where athlete_id='" + uid.String() + "'")
	if err != nil {
		log.Fatalf("query failed: %s", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("counted", count)
}
