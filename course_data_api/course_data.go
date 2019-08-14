package course_data_api

import (
	"encoding/json"
	"fmt"
	"github.com/andrewmthomas87/northwestern/models"
	"io/ioutil"
	"net/http"
	"strings"
)

var badResponseError = fmt.Errorf("request returned an error")

type Client struct {
	baseUrl         string
	apiKey          string
	apiKeyParameter string

	httpClient *http.Client
}

func NewClient(baseUrl, apiKey, apiKeyParameter string) *Client {
	return &Client{
		baseUrl:         baseUrl,
		apiKey:          apiKey,
		apiKeyParameter: apiKeyParameter,
		httpClient:      &http.Client{},
	}
}

func (c *Client) newRequest(endpoint string, parameters []string) (*http.Request, error) {
	parameters = append(parameters, fmt.Sprintf("%s=%s", c.apiKeyParameter, c.apiKey))
	url := fmt.Sprintf("%s%s?%s", c.baseUrl, endpoint, strings.Join(parameters, "&"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created request: %s\n", url)

	return req, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, badResponseError
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type apiTerm struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (c *Client) Terms() ([]*models.Term, error) {
	req, err := c.newRequest("terms", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var apiTerms []apiTerm
	if err := json.Unmarshal(body, &apiTerms); err != nil {
		return nil, err
	}

	terms := make([]*models.Term, len(apiTerms))
	for i, apiTerm := range apiTerms {
		terms[i] = &models.Term{
			Id:        apiTerm.Id,
			Name:      apiTerm.Name,
			StartDate: apiTerm.StartDate,
			EndDate:   apiTerm.EndDate,
		}
	}

	return terms, nil
}

type apiSchool struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func (c *Client) Schools() ([]*models.School, error) {
	req, err := c.newRequest("schools", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var apiSchools []apiSchool
	if err := json.Unmarshal(body, &apiSchools); err != nil {
		return nil, err
	}

	schools := make([]*models.School, len(apiSchools))
	for i, apiSchool := range apiSchools {
		schools[i] = &models.School{
			Symbol: apiSchool.Symbol,
			Name:   apiSchool.Name,
		}
	}

	return schools, nil
}

type apiSubject struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func (c *Client) Subjects(term int, school string) ([]*models.Subject, error) {
	var parameters []string
	if term != -1 {
		parameters = append(parameters, fmt.Sprintf("term=%d", term))
	}
	if len(school) > 0 {
		parameters = append(parameters, fmt.Sprintf("school=%s", school))
	}

	req, err := c.newRequest("subjects", parameters)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var apiSubjects []apiSubject
	if err := json.Unmarshal(body, &apiSubjects); err != nil {
		return nil, err
	}

	subjects := make([]*models.Subject, len(apiSubjects))
	for i, apiSubject := range apiSubjects {
		subjects[i] = &models.Subject{
			Symbol: apiSubject.Symbol,
			Name:   apiSubject.Name,
		}
	}

	return subjects, nil
}

type apiInstructor struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (c *Client) Instructors(subject string) ([]*models.Instructor, error) {
	parameters := []string{fmt.Sprintf("subject=%s", subject)}

	req, err := c.newRequest("instructors", parameters)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var apiInstructors []apiInstructor
	if err := json.Unmarshal(body, &apiInstructors); err != nil {
		return nil, err
	}

	instructors := make([]*models.Instructor, len(apiInstructors))
	for i, apiInstructor := range apiInstructors {
		instructors[i] = &models.Instructor{
			Id:    apiInstructor.Id,
			Name:  apiInstructor.Name,
			Phone: apiInstructor.Phone,
		}
	}

	return instructors, nil
}

type apiBuilding struct {
	Id   int     `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

func (c *Client) Buildings() ([]*models.Building, error) {
	req, err := c.newRequest("buildings", nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var apiBuildings []apiBuilding
	if err := json.Unmarshal(body, &apiBuildings); err != nil {
		return nil, err
	}

	buildings := make([]*models.Building, len(apiBuildings))
	for i, apiBuilding := range apiBuildings {
		buildings[i] = &models.Building{
			Id:   apiBuilding.Id,
			Name: apiBuilding.Name,
			Lat:  apiBuilding.Lat,
			Lon:  apiBuilding.Lon,
		}
	}

	return buildings, nil
}

type apiRoom struct {
	Id         int    `json:"id"`
	BuildingId int    `json:"building_id"`
	Name       string `json:"name"`
}

func (c *Client) Rooms(building int) ([]*models.Room, error) {
	parameters := []string{fmt.Sprintf("building=%d", building)}

	req, err := c.newRequest("rooms", parameters)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var apiRooms []apiRoom
	if err := json.Unmarshal(body, &apiRooms); err != nil {
		return nil, err
	}

	rooms := make([]*models.Room, len(apiRooms))
	for i, apiRoom := range apiRooms {
		rooms[i] = &models.Room{
			Id:         apiRoom.Id,
			BuildingId: apiRoom.BuildingId,
			Name:       apiRoom.Name,
		}
	}

	return rooms, nil
}

type apiCourse struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Term       string `json:"term"`
	School     string `json:"school"`
	Instructor struct {
		Name        string `json:"name"`
		Bio         string `json:"bio"`
		Address     string `json:"address"`
		Phone       string `json:"phone"`
		OfficeHours string `json:"office_hours"`
	} `json:"instructor"`
	Subject    string `json:"subject"`
	CatalogNum string `json:"catalog_num"`
	Section    string `json:"section"`
	Room       struct {
		Id           int    `json:"id"`
		BuildingId   int    `json:"building_id"`
		BuildingName string `json:"building_name"`
		Name         string `json:"name"`
	} `json:"room"`
	MeetingDays        string `json:"meeting_days"`
	StartTime          string `json:"start_time"`
	EndTime            string `json:"end_time"`
	StartDate          string `json:"start_date"`
	EndDate            string `json:"end_date"`
	Seats              int    `json:"seats"`
	Overview           string `json:"overview"`
	Topic              string `json:"topic"`
	Attributes         string `json:"attributes"`
	Requirements       string `json:"requirements"`
	Component          string `json:"component"`
	ClassNum           int    `json:"class_num"`
	CourseId           int    `json:"course_id"`
	CourseDescriptions []struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	} `json:"course_descriptions"`
	CourseComponents []struct {
		Component   string `json:"component"`
		MeetingDays string `json:"meeting_days"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
		Section     string `json:"section"`
		Room        string `json:"room"`
	} `json:"course_components"`
}

func (c *Client) Courses(term int, subject string, instructors map[string]int) ([]*models.Course, []*models.CourseDescription, []*models.CourseComponent, error) {
	parameters := []string{
		fmt.Sprintf("term=%d", term),
		fmt.Sprintf("subject=%s", subject),
	}

	req, err := c.newRequest("courses/details", parameters)
	if err != nil {
		return nil, nil, nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, nil, nil, err
	}

	var apiCourses []apiCourse
	if err := json.Unmarshal(body, &apiCourses); err != nil {
		return nil, nil, nil, err
	}

	courses := make([]*models.Course, len(apiCourses))
	courseDescriptions := make([]*models.CourseDescription, 0)
	courseComponents := make([]*models.CourseComponent, 0)
	for i, apiCourse := range apiCourses {
		courses[i] = &models.Course{
			Id:           apiCourse.Id,
			Title:        apiCourse.Title,
			Term:         term,
			School:       apiCourse.School,
			Instructor:   instructors[apiCourse.Instructor.Name],
			Subject:      apiCourse.Subject,
			CatalogNum:   apiCourse.CatalogNum,
			Section:      apiCourse.Section,
			Room:         apiCourse.Room.Id,
			MeetingDays:  apiCourse.MeetingDays,
			StartTime:    apiCourse.StartTime,
			EndTime:      apiCourse.EndTime,
			StartDate:    apiCourse.StartDate,
			EndDate:      apiCourse.EndDate,
			Seats:        apiCourse.Seats,
			Overview:     apiCourse.Overview,
			Topic:        apiCourse.Topic,
			Attributes:   apiCourse.Attributes,
			Requirements: apiCourse.Requirements,
			Component:    apiCourse.Component,
			ClassNum:     apiCourse.ClassNum,
			CourseId:     apiCourse.CourseId,
		}

		for _, apiCourseDescription := range apiCourse.CourseDescriptions {
			courseDescriptions = append(courseDescriptions, &models.CourseDescription{
				Course: apiCourse.Id,
				Name:   apiCourseDescription.Name,
				Desc:   apiCourseDescription.Desc,
			})
		}
		for _, apiCourseComponent := range apiCourse.CourseComponents {
			courseComponents = append(courseComponents, &models.CourseComponent{
				Course:      apiCourse.Id,
				Component:   apiCourseComponent.Component,
				MeetingDays: apiCourseComponent.MeetingDays,
				StartTime:   apiCourseComponent.StartTime,
				EndTime:     apiCourseComponent.EndTime,
				Section:     apiCourseComponent.Section,
				Room:        apiCourseComponent.Room,
			})
		}
	}

	return courses, courseDescriptions, courseComponents, nil
}
