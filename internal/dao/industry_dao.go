package dao

import (
	"context"

	"github.com/dqixuan/stock_info/internal/data"
	"github.com/dqixuan/stock_info/internal/model"
	"gorm.io/gorm"
)

// ==========================================
// 1. Industry DAO (行业信息表)
// ==========================================

type IndustryDao struct {
	db *gorm.DB
}

func NewIndustryDao(data *data.Data) *IndustryDao {
	return &IndustryDao{db: data.DB()}
}

// Create 创建单个行业
func (d *IndustryDao) Create(ctx context.Context, industry *model.Industry) error {
	return d.db.WithContext(ctx).Create(industry).Error
}

// BatchCreate 批量创建行业
func (d *IndustryDao) BatchCreate(ctx context.Context, industries []model.Industry) error {
	if len(industries) == 0 {
		return nil
	}
	// 每批次插入 100 条，防止 SQL 过长
	return d.db.WithContext(ctx).CreateInBatches(industries, 100).Error
}

// Update 根据主键 ID 更新行业信息 (注意：Save 会更新所有字段，包括零值)
func (d *IndustryDao) Update(ctx context.Context, industry *model.Industry) error {
	return d.db.WithContext(ctx).Save(industry).Error
}

// UpdateByIndustryID 根据业务行业代码更新名称和数量
func (d *IndustryDao) UpdateByIndustryID(ctx context.Context, industryID string, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&model.Industry{}).
		Where("INDUSTRY_ID = ?", industryID).
		Updates(updates).Error
}

// DeleteByID 根据主键 ID 删除
func (d *IndustryDao) DeleteByID(ctx context.Context, id uint64) error {
	return d.db.WithContext(ctx).Delete(&model.Industry{}, id).Error
}

// DeleteByIndustryID 根据业务行业代码删除
func (d *IndustryDao) DeleteByIndustryID(ctx context.Context, industryID string) error {
	return d.db.WithContext(ctx).Where("INDUSTRY_ID = ?", industryID).Delete(&model.Industry{}).Error
}

// GetByID 根据主键 ID 查询
func (d *IndustryDao) GetByID(ctx context.Context, id uint64) (*model.Industry, error) {
	var industry model.Industry
	err := d.db.WithContext(ctx).First(&industry, id).Error
	if err != nil {
		return nil, err
	}
	return &industry, nil
}

// GetByIndustryID 根据业务行业代码查询
func (d *IndustryDao) GetByIndustryID(ctx context.Context, industryID string) (*model.Industry, error) {
	var industry model.Industry
	err := d.db.WithContext(ctx).Where("INDUSTRY_ID = ?", industryID).First(&industry).Error
	if err != nil {
		return nil, err
	}
	return &industry, nil
}

// List 分页查询行业列表
func (d *IndustryDao) List(ctx context.Context, page, pageSize int) ([]model.Industry, int64, error) {
	var industries []model.Industry
	var total int64

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// 1. 查询总数
	if err := d.db.WithContext(ctx).Model(&model.Industry{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 2. 查询分页数据
	err := d.db.WithContext(ctx).Offset(offset).Limit(pageSize).Order("ID DESC").Find(&industries).Error
	return industries, total, err
}

func (d *IndustryDao) BatchInsert(ctx context.Context, industries []*model.Industry) error {
	if len(industries) == 0 {
		return nil
	}
	// 每批次插入 100 条，防止 SQL 过长
	return d.db.WithContext(ctx).CreateInBatches(industries, 100).Error
}
