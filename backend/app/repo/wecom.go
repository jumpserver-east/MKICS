package repo

import (
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"
	"context"

	"gorm.io/gorm"
)

type IWecomRepo interface {
	Get(opts ...DBOption) (model.WecomConfig, error)
	Update(ctx context.Context, conf *model.WecomConfig) error
	WithType(typ string) DBOption
}

type WecomRepo struct{}

func NewIWecomRepo() IWecomRepo {
	return &WecomRepo{}
}

func (r *WecomRepo) Get(opts ...DBOption) (model.WecomConfig, error) {
	var config model.WecomConfig

	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}

	if err := db.First(&config).Error; err != nil {
		return config, err
	}

	return config, nil
}

func (r *WecomRepo) Update(ctx context.Context, conf *model.WecomConfig) error {
	return getTx(ctx).Updates(&conf).Error
}

func (r *WecomRepo) WithType(typ string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("type = ?", typ)
	}
}
