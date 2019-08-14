package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/andrewmthomas87/northwestern/course_data_api"
	"github.com/andrewmthomas87/northwestern/database"
	"github.com/andrewmthomas87/northwestern/models"
	"github.com/spf13/viper"
	"log"
)

func terms(ctx context.Context, db *database.Database, apiClient *course_data_api.Client, term *models.Term) {
	if term != nil {
		fmt.Println("Skipping terms")

		return
	}

	fmt.Println("Fetching terms")

	terms, err := apiClient.Terms()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Storing terms")

	err = db.InsertTerms(ctx, terms)
	if err != nil {
		log.Fatal(err)
	}
}

func schools(ctx context.Context, db *database.Database, apiClient *course_data_api.Client, term *models.Term) {
	if term != nil {
		fmt.Println("Skipping schools")

		return
	}

	fmt.Println("Fetching schools")

	schools, err := apiClient.Schools()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Storing schools")

	err = db.InsertSchools(ctx, schools)
	if err != nil {
		log.Fatal(err)
	}
}

func subjects(ctx context.Context, db *database.Database, apiClient *course_data_api.Client, term *models.Term) {
	var terms []*models.Term
	if term != nil {
		fmt.Printf("Deleting subject availabilities for term %s\n", term.Name)

		err := db.DeleteSubjectAvailabilitiesByTerm(ctx, term.Id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Fetching subjects for term %s\n", term.Name)

		terms = []*models.Term{term}
	} else {
		fmt.Println("Fetching subjects")

		var err error
		terms, err = db.SelectAllTerms(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	schools, err := db.SelectAllSchools(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var subjects []*models.Subject
	var subjectAvailabilities []*models.SubjectAvailability
	for _, term := range terms {
		for _, school := range schools {
			filteredSubjects, err := apiClient.Subjects(term.Id, school.Symbol)
			if err != nil {
				log.Fatal(err)
			}

			subjects = append(subjects, filteredSubjects...)
			for _, subject := range filteredSubjects {
				subjectAvailabilities = append(subjectAvailabilities, &models.SubjectAvailability{
					Term:    term.Id,
					School:  school.Symbol,
					Subject: subject.Symbol,
				})
			}
		}
	}

	fmt.Println("Storing subjects")

	err = db.InsertSubjects(ctx, subjects)
	if err != nil {
		log.Fatal(err)
	}
	err = db.InsertSubjectAvailabilities(ctx, subjectAvailabilities)
	if err != nil {
		log.Fatal(err)
	}
}

func instructors(ctx context.Context, db *database.Database, apiClient *course_data_api.Client, term *models.Term) {
	if term != nil {
		fmt.Println("Skipping instructors")

		return
	}

	fmt.Println("Fetching instructors")

	subjects, err := db.SelectAllSubjects(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var instructors []*models.Instructor
	var instructorSubjects []*models.InstructorSubject
	for _, subject := range subjects {
		filteredInstructors, err := apiClient.Instructors(subject.Symbol)
		if err != nil {
			log.Fatal(err)
		}

		instructors = append(instructors, filteredInstructors...)
		for _, instructor := range filteredInstructors {
			instructorSubjects = append(instructorSubjects, &models.InstructorSubject{
				Instructor: instructor.Id,
				Subject:    subject.Symbol,
			})
		}
	}

	fmt.Println("Storing instructors")

	err = db.InsertInstructors(ctx, instructors)
	if err != nil {
		log.Fatal(err)
	}
	err = db.InsertInstructorSubjects(ctx, instructorSubjects)
	if err != nil {
		log.Fatal(err)
	}
}

func buildings(ctx context.Context, db *database.Database, apiClient *course_data_api.Client, term *models.Term) {
	if term != nil {
		fmt.Println("Skipping buildings")

		return
	}

	fmt.Println("Fetching buildings")

	buildings, err := apiClient.Buildings()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Storing buildings")

	err = db.InsertBuildings(ctx, buildings)
	if err != nil {
		log.Fatal(err)
	}
}

func rooms(ctx context.Context, db *database.Database, apiClient *course_data_api.Client, term *models.Term) {
	if term != nil {
		fmt.Println("Skipping rooms")

		return
	}

	fmt.Println("Fetching rooms")

	buildings, err := db.SelectAllBuildings(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var rooms []*models.Room
	for _, building := range buildings {
		filteredRooms, err := apiClient.Rooms(building.Id)
		if err != nil {
			log.Fatal(err)
		}

		rooms = append(rooms, filteredRooms...)
	}

	fmt.Println("Storing rooms")

	err = db.InsertRooms(ctx, rooms)
	if err != nil {
		log.Fatal(err)
	}
}

func courses(ctx context.Context, db *database.Database, apiClient *course_data_api.Client, term *models.Term) {
	fmt.Printf("Fetching courses for term %s\n", term.Name)

	subjects, err := db.SelectSubjectsByTerm(ctx, term.Id)
	if err != nil {
		log.Fatal(err)
	}

	instructors, err := db.SelectAllInstructors(ctx)
	if err != nil {
		log.Fatal(err)
	}

	instructorsMap := make(map[string]int, len(instructors))
	for _, instructor := range instructors {
		instructorsMap[instructor.Name] = instructor.Id
	}

	var courses []*models.Course
	var courseDescriptions []*models.CourseDescription
	var courseComponents []*models.CourseComponent
	for _, subject := range subjects {
		filteredCourses, filteredCourseDescriptions, filteredCourseComponents, err := apiClient.Courses(term.Id, subject.Symbol, instructorsMap)
		if err != nil {
			log.Fatal(err)
		}

		courses = append(courses, filteredCourses...)
		courseDescriptions = append(courseDescriptions, filteredCourseDescriptions...)
		courseComponents = append(courseComponents, filteredCourseComponents...)
	}

	fmt.Println("Storing courses")

	err = db.InsertCourses(ctx, courses)
	if err != nil {
		log.Fatal(err)
	}

	err = db.InsertCourseDescriptions(ctx, courseDescriptions)
	if err != nil {
		log.Fatal(err)
	}

	err = db.InsertCourseComponents(ctx, courseComponents)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx := context.Background()

	termName := flag.String("term", "", "fetch data for a specific term")
	courseTermName := flag.String("courses", "", "fetch course data for a specific term")
	flag.Parse()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	db, err := database.NewDatabase(viper.GetString("database.user"), viper.GetString("database.password"), viper.GetString("database.host"), viper.GetInt("database.port"), viper.GetString("database.database"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	apiClient := course_data_api.NewClient(viper.GetString("courseDataAPI.baseUrl"), viper.GetString("courseDataAPI.apiKey"), viper.GetString("courseDataAPI.apiKeyParameter"))

	var term *models.Term
	if len(*courseTermName) > 0 {
		term, err = db.SelectTermByName(ctx, *courseTermName)
		if err != nil {
			log.Fatal(err)
		}

		courses(ctx, db, apiClient, term)
		return
	} else if len(*termName) > 0 {
		term, err = db.SelectTermByName(ctx, *termName)
		if err != nil {
			log.Fatal(err)
		}
	}

	terms(ctx, db, apiClient, term)
	schools(ctx, db, apiClient, term)
	subjects(ctx, db, apiClient, term)
	instructors(ctx, db, apiClient, term)
	buildings(ctx, db, apiClient, term)
	rooms(ctx, db, apiClient, term)
}
