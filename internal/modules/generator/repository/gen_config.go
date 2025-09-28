package repository

import (
	"fmt"

	"github.com/LiteMove/light-stack/internal/modules/generator/model"
	"gorm.io/gorm"
)

// GenConfigRepository 代码生成配置仓储
type GenConfigRepository struct {
	db *gorm.DB
}

// NewGenConfigRepository 创建代码生成配置仓储
func NewGenConfigRepository(db *gorm.DB) *GenConfigRepository {
	return &GenConfigRepository{
		db: db,
	}
}

// Create 创建配置
func (r *GenConfigRepository) Create(config *model.GenTableConfig) (*model.GenTableConfig, error) {
	if err := r.db.Create(config).Error; err != nil {
		return nil, fmt.Errorf("创建配置失败: %v", err)
	}
	return config, nil
}

// Update 更新配置
func (r *GenConfigRepository) Update(config *model.GenTableConfig) (*model.GenTableConfig, error) {
	if err := r.db.Save(config).Error; err != nil {
		return nil, fmt.Errorf("更新配置失败: %v", err)
	}
	return config, nil
}

// GetByID 根据ID获取配置
func (r *GenConfigRepository) GetByID(id int64) (*model.GenTableConfig, error) {
	var config model.GenTableConfig
	err := r.db.Preload("Columns", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort ASC")
	}).First(&config, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("配置不存在")
		}
		return nil, fmt.Errorf("查询配置失败: %v", err)
	}

	return &config, nil
}

// GetByTableName 根据表名获取配置
func (r *GenConfigRepository) GetByTableName(tableName string) (*model.GenTableConfig, error) {
	var config model.GenTableConfig
	err := r.db.Preload("Columns", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort ASC")
	}).Where("table_name = ?", tableName).First(&config).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("配置不存在")
		}
		return nil, fmt.Errorf("查询配置失败: %v", err)
	}

	return &config, nil
}

// GetList 获取配置列表
func (r *GenConfigRepository) GetList(page, size int, tableName, businessName string) ([]*model.GenTableConfig, int64, error) {
	var configs []*model.GenTableConfig
	var total int64

	query := r.db.Model(&model.GenTableConfig{})

	// 条件查询
	if tableName != "" {
		query = query.Where("table_name LIKE ?", "%"+tableName+"%")
	}
	if businessName != "" {
		query = query.Where("business_name LIKE ?", "%"+businessName+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * size
	err := query.Preload("Columns", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort ASC")
	}).Offset(offset).Limit(size).Order("created_at DESC").Find(&configs).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询配置列表失败: %v", err)
	}

	return configs, total, nil
}

// Delete 删除配置
func (r *GenConfigRepository) Delete(id int64) error {
	// 使用事务删除配置和关联的字段配置
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除字段配置
		if err := tx.Where("table_config_id = ?", id).Delete(&model.GenTableColumn{}).Error; err != nil {
			return fmt.Errorf("删除字段配置失败: %v", err)
		}

		// 删除表配置
		if err := tx.Delete(&model.GenTableConfig{}, id).Error; err != nil {
			return fmt.Errorf("删除表配置失败: %v", err)
		}

		return nil
	})
}

// DeleteColumns 删除字段配置
func (r *GenConfigRepository) DeleteColumns(tableConfigID int64) error {
	if err := r.db.Where("table_config_id = ?", tableConfigID).Delete(&model.GenTableColumn{}).Error; err != nil {
		return fmt.Errorf("删除字段配置失败: %v", err)
	}
	return nil
}

// BatchDelete 批量删除配置
func (r *GenConfigRepository) BatchDelete(ids []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除字段配置
		if err := tx.Where("table_config_id IN ?", ids).Delete(&model.GenTableColumn{}).Error; err != nil {
			return fmt.Errorf("批量删除字段配置失败: %v", err)
		}

		// 删除表配置
		if err := tx.Where("id IN ?", ids).Delete(&model.GenTableConfig{}).Error; err != nil {
			return fmt.Errorf("批量删除表配置失败: %v", err)
		}

		return nil
	})
}

// Exists 检查配置是否存在
func (r *GenConfigRepository) Exists(id int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.GenTableConfig{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查配置是否存在失败: %v", err)
	}
	return count > 0, nil
}

// ExistsByTableName 检查表名对应的配置是否存在
func (r *GenConfigRepository) ExistsByTableName(tableName string) (bool, error) {
	var count int64
	err := r.db.Model(&model.GenTableConfig{}).Where("table_name = ?", tableName).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查配置是否存在失败: %v", err)
	}
	return count > 0, nil
}

// GetAllConfigs 获取所有配置（不分页）
func (r *GenConfigRepository) GetAllConfigs() ([]*model.GenTableConfig, error) {
	var configs []*model.GenTableConfig
	err := r.db.Preload("Columns", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort ASC")
	}).Order("created_at DESC").Find(&configs).Error

	if err != nil {
		return nil, fmt.Errorf("查询所有配置失败: %v", err)
	}

	return configs, nil
}

// UpdateColumn 更新字段配置
func (r *GenConfigRepository) UpdateColumn(column *model.GenTableColumn) error {
	if err := r.db.Save(column).Error; err != nil {
		return fmt.Errorf("更新字段配置失败: %v", err)
	}
	return nil
}

// GetColumnByID 根据ID获取字段配置
func (r *GenConfigRepository) GetColumnByID(id int64) (*model.GenTableColumn, error) {
	var column model.GenTableColumn
	err := r.db.First(&column, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("字段配置不存在")
		}
		return nil, fmt.Errorf("查询字段配置失败: %v", err)
	}
	return &column, nil
}

// CreateHistory 创建历史记录
func (r *GenConfigRepository) CreateHistory(history *model.GenHistory) (*model.GenHistory, error) {
	if err := r.db.Create(history).Error; err != nil {
		return nil, fmt.Errorf("创建历史记录失败: %v", err)
	}
	return history, nil
}

// GetHistoryList 获取历史记录列表
func (r *GenConfigRepository) GetHistoryList(page, size int, tableName string) ([]*model.GenHistory, int64, error) {
	var histories []*model.GenHistory
	var total int64

	query := r.db.Model(&model.GenHistory{})

	// 条件查询
	if tableName != "" {
		query = query.Where("table_name LIKE ?", "%"+tableName+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询历史记录总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * size
	err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&histories).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询历史记录列表失败: %v", err)
	}

	return histories, total, nil
}

// UpdateHistoryDownloadCount 更新历史记录下载次数
func (r *GenConfigRepository) UpdateHistoryDownloadCount(id int64) error {
	if err := r.db.Model(&model.GenHistory{}).Where("id = ?", id).UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error; err != nil {
		return fmt.Errorf("更新下载次数失败: %v", err)
	}
	return nil
}

// GetHistoryByID 根据ID获取历史记录
func (r *GenConfigRepository) GetHistoryByID(id int64) (*model.GenHistory, error) {
	var history model.GenHistory
	err := r.db.First(&history, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("历史记录不存在")
		}
		return nil, fmt.Errorf("查询历史记录失败: %v", err)
	}
	return &history, nil
}
