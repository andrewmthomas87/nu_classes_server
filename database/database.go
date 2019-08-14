package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/andrewmthomas87/northwestern/models"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(user, password, host string, port int, database string) (*Database, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database))
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) InsertTerms(ctx context.Context, terms []*models.Term) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("REPLACE INTO terms (id, name, start_date, end_date) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	for _, term := range terms {
		_, err = stmt.Exec(term.Id, term.Name, term.StartDate, term.EndDate)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Database) SelectAllTerms(ctx context.Context) ([]*models.Term, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, name, start_date, end_date FROM terms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terms []*models.Term
	for rows.Next() {
		term := &models.Term{}
		if err := rows.Scan(&term.Id, &term.Name, &term.StartDate, &term.EndDate); err != nil {
			return nil, err
		}
		terms = append(terms, term)
	}

	return terms, nil
}

func (d *Database) SelectTermByName(ctx context.Context, name string) (*models.Term, error) {
	row := d.db.QueryRowContext(ctx, "SELECT id, name, start_date, end_date FROM terms WHERE name=?", name)

	term := &models.Term{}
	if err := row.Scan(&term.Id, &term.Name, &term.StartDate, &term.EndDate); err != nil {
		return nil, err
	}

	return term, nil
}

func (d *Database) InsertSchools(ctx context.Context, schools []*models.School) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("REPLACE INTO schools (symbol, name) VALUES (?, ?)")
	if err != nil {
		return err
	}
	for _, school := range schools {
		_, err = stmt.Exec(school.Symbol, school.Name)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Database) SelectAllSchools(ctx context.Context) ([]*models.School, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT symbol, name FROM schools")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schools []*models.School
	for rows.Next() {
		school := &models.School{}
		if err := rows.Scan(&school.Symbol, &school.Name); err != nil {
			return nil, err
		}
		schools = append(schools, school)
	}

	return schools, nil
}

func (d *Database) InsertSubjects(ctx context.Context, subjects []*models.Subject) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("REPLACE  INTO subjects (symbol, name) VALUES (?, ?)")
	if err != nil {
		return err
	}
	for _, subject := range subjects {
		_, err = stmt.Exec(subject.Symbol, subject.Name)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Database) SelectAllSubjects(ctx context.Context) ([]*models.Subject, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT symbol, name FROM subjects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []*models.Subject
	for rows.Next() {
		subject := &models.Subject{}
		if err := rows.Scan(&subject.Symbol, &subject.Name); err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (d *Database) SelectSubjectsByTerm(ctx context.Context, term int) ([]*models.Subject, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT symbol, name FROM subjects, subject_availabilities WHERE symbol=subject AND term=?", term)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []*models.Subject
	for rows.Next() {
		subject := &models.Subject{}
		if err := rows.Scan(&subject.Symbol, &subject.Name); err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (d *Database) InsertSubjectAvailabilities(ctx context.Context, subjectAvailabilites []*models.SubjectAvailability) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT IGNORE INTO subject_availabilities (term, school, subject) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	for _, subjectAvailability := range subjectAvailabilites {
		_, err = stmt.Exec(subjectAvailability.Term, subjectAvailability.School, subjectAvailability.Subject)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Database) DeleteSubjectAvailabilitiesByTerm(ctx context.Context, termId int) error {
	_, err := d.db.ExecContext(ctx, "DELETE FROM subject_availabilities WHERE term=?", termId)
	return err
}

func (d *Database) InsertInstructors(ctx context.Context, instructors []*models.Instructor) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("REPLACE INTO instructors (id, name, phone) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	for _, instructor := range instructors {
		_, err := stmt.Exec(instructor.Id, instructor.Name, instructor.Phone)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Database) SelectAllInstructors(ctx context.Context) ([]*models.Instructor, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, name, phone FROM instructors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instructors []*models.Instructor
	for rows.Next() {
		instructor := &models.Instructor{}
		if err := rows.Scan(&instructor.Id, &instructor.Name, &instructor.Phone); err != nil {
			return nil, err
		}
		instructors = append(instructors, instructor)
	}

	return instructors, nil
}

func (d *Database) InsertInstructorSubjects(ctx context.Context, instructorSubjects []*models.InstructorSubject) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT IGNORE INTO instructor_subjects (instructor, subject) VALUES (?, ?)")
	if err != nil {
		return err
	}
	for _, instructorSubject := range instructorSubjects {
		_, err := stmt.Exec(instructorSubject.Instructor, instructorSubject.Subject)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Database) InsertBuildings(ctx context.Context, buildings []*models.Building) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("REPLACE  INTO buildings (id, name, lat, lon) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	for _, building := range buildings {
		_, err := stmt.Exec(building.Id, building.Name, building.Lat, building.Lon)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Database) SelectAllBuildings(ctx context.Context) ([]*models.Building, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, name, lat, lon FROM buildings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []*models.Building
	for rows.Next() {
		building := &models.Building{}
		if err := rows.Scan(&building.Id, &building.Name, &building.Lat, &building.Lon); err != nil {
			return nil, err
		}
		buildings = append(buildings, building)
	}

	return buildings, nil
}

func (d *Database) SelectBuilding(ctx context.Context, id int) (*models.Building, error) {
	row := d.db.QueryRowContext(ctx, "SELECT id, name, lat, lon FROM buildings WHERE id=?", id)

	building := &models.Building{}
	if err := row.Scan(&building.Id, &building.Name, &building.Lat, &building.Lon); err != nil {
		return nil, err
	}

	return building, nil
}

func (d *Database) InsertRooms(ctx context.Context, rooms []*models.Room) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("REPLACE INTO rooms (id, building_id, name) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	for _, room := range rooms {
		_, err := stmt.Exec(room.Id, room.BuildingId, room.Name)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Database) SelectAllRooms(ctx context.Context) ([]*models.Room, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, building_id, name FROM rooms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		if err := rows.Scan(&room.Id, &room.BuildingId, &room.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (d *Database) SelectRoomsByBuilding(ctx context.Context, building int) ([]*models.Room, error) {
	rows, err := d.db.QueryContext(ctx, "SELECT id, building_id, name FROM rooms WHERE building_id=?", building)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		if err := rows.Scan(&room.Id, &room.BuildingId, &room.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (d *Database) InsertCourses(ctx context.Context, courses []*models.Course) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("REPLACE INTO courses (id, title, term, school, instructor, subject, catalog_num, section, room, meeting_days, start_time, end_time, start_date, end_date, seats, overview, topic, attributes, requirements, component, class_num, course_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	for _, course := range courses {
		_, err := stmt.Exec(course.Id, course.Title, course.Term, course.School, course.Instructor, course.Subject, course.CatalogNum, course.Section, course.Room, course.MeetingDays, course.StartTime, course.EndTime, course.StartDate, course.EndDate, course.Seats, course.Overview, course.Topic, course.Attributes, course.Requirements, course.Component, course.ClassNum, course.CourseId)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Database) InsertCourseDescriptions(ctx context.Context, courseDescriptions []*models.CourseDescription) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("REPLACE INTO course_descriptions (course, name, description) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	for _, courseDescription := range courseDescriptions {
		_, err := stmt.Exec(courseDescription.Course, courseDescription.Name, courseDescription.Desc)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Database) InsertCourseComponents(ctx context.Context, courseComponents []*models.CourseComponent) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("REPLACE INTO course_components (course, component, meeting_days, start_time, end_time, section, room) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	for _, courseComponent := range courseComponents {
		_, err := stmt.Exec(courseComponent.Course, courseComponent.Component, courseComponent.MeetingDays, courseComponent.StartTime, courseComponent.EndTime, courseComponent.Section, courseComponent.Room)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
