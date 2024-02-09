package data

import (
	"context"
	"database/sql"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"kratos-k8s-job/internal/biz"
	"kratos-k8s-job/internal/conf"

	"kratos-k8s-job/internal/data/mysql"
)

type (
	mAdapter struct {
		MySqlDB *sql.DB
		log     *log.Helper
	}
)

func NewMySqlAdapter(c *conf.Data, logger log.Logger) (biz.MySqlAdapter, func(), error) {

	l := log.NewHelper(logger)

	dbCf := c.Database
	l.Debug("mysql source: ", dbCf.GetSource())
	db, err := sql.Open(dbCf.GetDriver(), dbCf.GetSource())
	if err != nil {
		l.Error("Fail on connect to MySql")
		return nil, nil, err
	}

	cleanup := func() {
		l.Info("closing mysql connection")
		db.Close()
	}

	return &mAdapter{
		MySqlDB: db,
		log:     log.NewHelper(logger),
	}, cleanup, nil
}

func (r *mAdapter) QueryMySqlDB(ctx context.Context) ([]biz.Message, error) {

	queries := mysql.New(r.MySqlDB)

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
