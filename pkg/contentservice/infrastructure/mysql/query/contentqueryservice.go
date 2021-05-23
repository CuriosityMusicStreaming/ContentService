package query

import (
	"contentservice/pkg/contentservice/app/query"
	"contentservice/pkg/contentservice/app/service"
	"fmt"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"strings"
)

func NewContentQueryService(client mysql.Client) query.ContentQueryService {
	return &contentQueryService{client: client}
}

type contentQueryService struct {
	client mysql.Client
}

func (service *contentQueryService) ContentList(spec query.ContentSpecification) ([]query.ContentView, error) {
	selectSql := `SELECT * from content`
	conditions, args, err := getWhereConditionsBySpec(spec)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if conditions != "" {
		selectSql += fmt.Sprintf(` WHERE %s`, conditions)
	}

	var contents []sqlxContent

	err = service.client.Select(&contents, selectSql, args...)
	if err != nil {
		return nil, err
	}

	return convertContents(contents), nil
}

func getWhereConditionsBySpec(spec query.ContentSpecification) (string, []interface{}, error) {
	var conditions []string
	var params []interface{}

	if len(spec.ContentIDs) != 0 {
		ids := marshalUUIDS(spec.ContentIDs)
		sqlQuery, args, err := sqlx.In(`content_id IN (?)`, ids)
		if err != nil {
			return "", nil, errors.WithStack(err)
		}
		conditions = append(conditions, sqlQuery)
		for _, arg := range args {
			params = append(params, arg)
		}
	}

	if len(spec.AuthorIDs) != 0 {
		ids := marshalUUIDS(spec.AuthorIDs)
		sqlQuery, args, err := sqlx.In(`author_id IN (?)`, ids)
		if err != nil {
			return "", nil, errors.WithStack(err)
		}
		conditions = append(conditions, sqlQuery)
		for _, arg := range args {
			params = append(params, arg)
		}
	}

	return strings.Join(conditions, " AND "), params, nil
}

func convertContents(contents []sqlxContent) []query.ContentView {
	res := make([]query.ContentView, 0, len(contents))
	for _, content := range contents {
		res = append(res, convertContent(content))
	}
	return res
}

func convertContent(content sqlxContent) query.ContentView {
	return query.ContentView{
		ID:               content.ID,
		Name:             content.Name,
		AuthorID:         content.AuthorID,
		ContentType:      service.ContentType(content.ContentType),
		AvailabilityType: service.ContentAvailabilityType(content.AvailabilityType),
	}
}

func marshalUUIDS(uuids []uuid.UUID) [][]byte {
	res := make([][]byte, len(uuids))
	for _, id := range uuids {
		marshaled, _ := id.MarshalBinary()
		res = append(res, marshaled)
	}
	return res
}

type sqlxContent struct {
	ID               uuid.UUID `db:"content_id"`
	Name             string    `db:"name"`
	AuthorID         uuid.UUID `db:"author_id"`
	ContentType      int       `db:"type"`
	AvailabilityType int       `db:"availability_type"`
}
