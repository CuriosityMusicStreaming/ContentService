package query

import (
	"contentservice/pkg/contentservice/app/query"
	"contentservice/pkg/contentservice/app/service"
	"fmt"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewContentQueryService(client mysql.Client) query.ContentQueryService {
	return &contentQueryService{client: client}
}

type contentQueryService struct {
	client mysql.Client
}

func (service *contentQueryService) ContentList(spec query.ContentSpecification) ([]query.ContentView, error) {
	selectSql := `SELECT * from content`
	selectSql, args, err := conditionsBySpec(selectSql, spec)
	if err != nil {
		return nil, err
	}

	var contents []sqlxContent

	err = service.client.Select(&contents, selectSql, args...)
	if err != nil {
		return nil, err
	}

	return convertContents(contents), nil
}

func conditionsBySpec(query string, spec query.ContentSpecification) (string, []interface{}, error) {
	if len(spec.ContentIDs) != 0 {
		queryString, args, err := sqlx.In(fmt.Sprintf(`%s WHERE content_id IN (?)`, query), marshalUUIDS(spec.ContentIDs))
		if err != nil {
			return "", nil, err
		}

		return queryString, args, nil
	}

	return query, nil, nil
}

func convertContents(contents []sqlxContent) []query.ContentView {
	res := make([]query.ContentView, len(contents))
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
