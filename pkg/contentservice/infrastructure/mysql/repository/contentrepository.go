package repository

import (
	"database/sql"

	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"contentservice/pkg/contentservice/domain"
)

func NewContentRepository(client mysql.Client) domain.ContentRepository {
	return &contentRepository{
		client: client,
	}
}

type contentRepository struct {
	client mysql.Client
}

func (repo *contentRepository) NewID() domain.ContentID {
	return domain.ContentID(uuid.New())
}

func (repo *contentRepository) Find(contentID domain.ContentID) (domain.Content, error) {
	const selectSql = `SELECT * from content WHERE content_id = ?`

	binaryUUID, err := uuid.UUID(contentID).MarshalBinary()
	if err != nil {
		return domain.Content{}, err
	}

	var content sqlxContent

	err = repo.client.Get(&content, selectSql, binaryUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Content{}, domain.ErrContentNotFound
		}
		return domain.Content{}, errors.WithStack(err)
	}

	return domain.Content{
		ID:               domain.ContentID(content.ContentID),
		Title:            content.Title,
		AuthorID:         domain.AuthorID(content.AuthorID),
		ContentType:      domain.ContentType(content.Type),
		AvailabilityType: domain.ContentAvailabilityType(content.AvailabilityType),
	}, nil
}

func (repo *contentRepository) Store(content domain.Content) error {
	const insertSql = `
		INSERT INTO content (content_id, title, author_id, type, availability_type) VALUES(?, ?, ?, ?, ?)
		ON DUPLICATE KEY 
		UPDATE content_id=VALUES(content_id), title=VALUES(title), author_id=VALUES(author_id), type=VALUES(type), availability_type=VALUES(availability_type)
	`

	binaryUUID, err := uuid.UUID(content.ID).MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	authorBinaryUUID, err := uuid.UUID(content.AuthorID).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = repo.client.Exec(insertSql, binaryUUID, content.Title, authorBinaryUUID, content.ContentType, content.AvailabilityType)
	return err
}

func (repo *contentRepository) Remove(contentID domain.ContentID) error {
	const deleteSql = `DELETE FROM content WHERE content_id = ?`

	binaryUUID, err := uuid.UUID(contentID).MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = repo.client.Exec(deleteSql, binaryUUID)
	return err

}

type sqlxContent struct {
	ContentID        uuid.UUID `db:"content_id"`
	Title            string    `db:"title"`
	AuthorID         uuid.UUID `db:"author_id"`
	Type             int       `db:"type"`
	AvailabilityType int       `db:"availability_type"`
}
