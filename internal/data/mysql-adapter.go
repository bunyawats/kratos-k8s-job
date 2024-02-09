package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-k8s-job/internal/biz"

	"kratos-k8s-job/internal/data/mysql"
)

type (
	M struct {
		data *Data
		log  *log.Helper
	}
)

func NewMySqlAdapter(data *Data, logger log.Logger) biz.MySqlAdapter {
	return &M{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *M) QueryMySqlDB(ctx context.Context) ([]biz.Message, error) {

	queries := mysql.New(r.data.MySqlDB)

	var lastTempalteId int64
	// get current template
	currentTemplate, err := queries.GetCurrentTemplate(ctx)
	if err != nil {
		log.Warnf(err.Error())
		lastTempalteId = -1
	} else {
		lastTempalteId = currentTemplate.ConsentTemplateID
	}
	log.Infof("Found Last Tempalte ID: %v", lastTempalteId)

	lastUpdatedTemplateList, err := queries.ListAllLastUpdatedTemplate(ctx, lastTempalteId)
	if err != nil {
		log.Warnf(err.Error())
	}

	messageList := make([]biz.Message, 0)

	for _, consentTemplate := range lastUpdatedTemplateList {
		result, err := queries.CreateCurrentTemplate(ctx, mysql.CreateCurrentTemplateParams{
			ConsentTemplateID: consentTemplate.ID,
			TemplateName:      consentTemplate.TemplateName,
			Version:           consentTemplate.Version,
		})
		if err != nil {
			return nil, err
		}
		insertedCurrentTemplateID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		log.Info("insertedCurrentTemplateID", insertedCurrentTemplateID)

		messageList = append(
			messageList,
			biz.Message{
				TemPlateName: consentTemplate.TemplateName,
				Version:      consentTemplate.Version,
			},
		)
	}
	return messageList, nil
}
