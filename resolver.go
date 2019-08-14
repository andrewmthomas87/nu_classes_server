package northwestern

import (
	"context"
	"github.com/andrewmthomas87/northwestern/database"
	"github.com/andrewmthomas87/northwestern/generated"
	"github.com/andrewmthomas87/northwestern/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	Db *database.Database
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Room() generated.RoomResolver {
	return &roomResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Terms(ctx context.Context) ([]*models.Term, error) {
	terms, err := r.Db.SelectAllTerms(ctx)
	if err != nil {
		return nil, err
	}

	return terms, nil
}

func (r *queryResolver) Schools(ctx context.Context) ([]*models.School, error) {
	schools, err := r.Db.SelectAllSchools(ctx)
	if err != nil {
		return nil, err
	}

	return schools, nil
}

func (r *queryResolver) Subjects(ctx context.Context) ([]*models.Subject, error) {
	subjects, err := r.Db.SelectAllSubjects(ctx)
	if err != nil {
		return nil, err
	}

	return subjects, nil
}

func (r *queryResolver) SubjectsByTerm(ctx context.Context, term int) ([]*models.Subject, error) {
	subjects, err := r.Db.SelectSubjectsByTerm(ctx, term)
	if err != nil {
		return nil, err
	}

	return subjects, nil
}

func (r *queryResolver) Buildings(ctx context.Context) ([]*models.Building, error) {
	buildings, err := r.Db.SelectAllBuildings(ctx)
	if err != nil {
		return nil, err
	}

	return buildings, nil
}

func (r *queryResolver) Rooms(ctx context.Context) ([]*models.Room, error) {
	rooms, err := r.Db.SelectAllRooms(ctx)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *queryResolver) RoomsByBuilding(ctx context.Context, building int) ([]*models.Room, error) {
	rooms, err := r.Db.SelectRoomsByBuilding(ctx, building)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

type roomResolver struct{ *Resolver }

func (r *roomResolver) Building(ctx context.Context, obj *models.Room) (*models.Building, error) {
	building, err := r.Db.SelectBuilding(ctx, obj.BuildingId)
	if err != nil {
		return nil, err
	}

	return building, nil
}
