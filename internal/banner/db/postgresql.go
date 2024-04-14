package db

import (
	"avito_banners/internal/banner"
	"avito_banners/pkg/client/postgresql"
	"avito_banners/pkg/logging"
	"avito_banners/pkg/utils"
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

type storage struct {
	client postgresql.Client
	logger *logging.Logger
}

var connect pgx.Conn

func (s *storage) GetBanner(tagID int, featureID int, banner *banner.Banner) error {
	return connect.QueryRow(context.Background(), `
				SELECT b.id, b.data_id, b.is_active
				FROM banners b
				INNER JOIN banner_tags bt ON b.id = bt.banner_id
				WHERE b.feature_id = $1 AND bt.tag_id = $2
			`, featureID, tagID).Scan(&banner.ID, &banner.DataID, &banner.IsActive)
}


func (s *storage) GetAllBanners(tagID int, featureID int, limit int, offset int) ([]banner.Banner, error) {
	rows, err := connect.Query(context.Background(), `
	SELECT DISTINCT b.id, b.data_id, b.is_active
	FROM banners b
	INNER JOIN banner_tags bt ON b.id = bt.banner_id
	WHERE ($1 = -1 or b.feature_id = $1) AND ($2 = -1 OR bt.tag_id = $2)
	LIMIT $3
	OFFSET $4
`, featureID, tagID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []banner.Banner
	for rows.Next() {
		var banner banner.Banner
		if err := rows.Scan(&banner.ID, &banner.DataID, &banner.IsActive); err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}
	return banners, nil
}


func (s *storage) CreateBanner(requestBody *banner.CreateBannerDTO) (int, error) {
	for _, tag := range requestBody.TagIds {
		var banner banner.Banner
		err := s.GetBanner(tag, requestBody.FeatureId, &banner)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				continue
			}
		}
		return 0, errors.New("banner already exist")
	}

	nextId := utils.GenerateNextId()
	var insertedID int
	err := connect.QueryRow(context.Background(), `
					INSERT INTO banners (feature_id, data_id, is_active)
					VALUES ($1, $2, $3)
					RETURNING id;
				`, requestBody.FeatureId, strconv.Itoa(nextId), requestBody.IsActive).Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	log.Printf("Вставлена новая строка %v в таблицу banners", insertedID)

	for ind, tag := range requestBody.TagIds {
		log.Printf("ind = ", ind)
		_, err = connect.Exec(context.Background(), `
					INSERT INTO banner_tags (banner_id, tag_id)
					VALUES ($1, $2);
				`, insertedID, tag)
		if err != nil {
			return 0, err
		}
	}
	log.Printf("Вставлены новые строки в таблицу banner_tags")

	return nextId, nil
}



func (s *storage) DeleteBanner(id int) (int, error) {
	result, err := connect.Exec(context.Background(), `
					DELETE FROM banner_tags
					WHERE banner_id = $1
				`, id)
	if err != nil {
		return 0, err
	}
	count := result.RowsAffected()

	return int(count), nil
}

func (s *storage) UpdateBanner(id int, requestBody *banner.CreateBannerDTO) (int, error) {
	for _, tag := range requestBody.TagIds {
		var banner banner.Banner
		err := s.GetBanner(tag, requestBody.FeatureId, &banner)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				continue
			}
			return 1, err
		}
		if banner.ID == id {
			continue
		}
		
		return 0, errors.New("banner already exist")
	}

	var dataIdStr string
	err := connect.QueryRow(context.Background(), `
					SELECT data_id
					FROM banners
					WHERE id = $1
				`, id).Scan(&dataIdStr)
	if err != nil {
		return 1, err
	}
	

	_, err = connect.Exec(context.Background(), `
					UPDATE banners
					SET feature_id = $2, is_active = $3
					WHERE id = $1;
				`, id, requestBody.FeatureId, requestBody.IsActive)
	if err != nil {
		log.Printf("id = %v[%T], featureId = %v[%T], isActicve = %v[%T]", id, id, requestBody.FeatureId, requestBody.FeatureId, requestBody.IsActive, requestBody.IsActive)
		log.Printf("err = %+v", err.Error())
		return 1, err
	}


	_, err = connect.Exec(context.Background(), `
					DELETE FROM banner_tags
					WHERE banner_id = $1;
				`, id)
	if err != nil {
		return 1, err
	}

	
	for _, tag := range requestBody.TagIds {
		_, err = connect.Exec(context.Background(), `
			INSERT INTO banner_tags (banner_id, tag_id)
			VALUES ($1, $2);
		`, id, tag)
		if err != nil {
			return 1, err
		}
	}

	return 0 , nil
}

func NewStorage(client postgresql.Client, logger *logging.Logger) banner.Storage {
	return &storage{
		client: client,
		logger: logger,
	}
}
