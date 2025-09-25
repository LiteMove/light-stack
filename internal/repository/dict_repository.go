package repository

import (
	"errors"
	"github.com/LiteMove/light-stack/internal/model"
	"gorm.io/gorm"
)

// DictRepository 字典数据访问接口
type DictRepository interface {
	// 字典类型相关方法
	CreateType(dictType *model.DictType) error
	GetTypeByID(id uint64) (*model.DictType, error)
	GetTypeByType(typeCode string) (*model.DictType, error)
	UpdateType(dictType *model.DictType) error
	DeleteType(id uint64) error
	GetTypeList(page, pageSize int, status int, name string) ([]*model.DictType, int64, error)
	TypeExists(typeCode string) (bool, error)

	// 字典数据相关方法
	CreateData(dictData *model.DictData) error
	GetDataByID(id uint64) (*model.DictData, error)
	UpdateData(dictData *model.DictData) error
	DeleteData(id uint64) error
	GetDataList(dictType string, page, pageSize int, status int, label string) ([]*model.DictData, int64, error)
	GetDataByType(dictType string) ([]*model.DictData, error)
	GetEnabledDataByType(dictType string) ([]*model.DictData, error)
	DataExists(dictType, value string) (bool, error)
	GetDataByValue(dictType, value string) (*model.DictData, error)
}

// dictRepository 字典数据访问实现
type dictRepository struct {
	db *gorm.DB
}

// NewDictRepository 创建字典数据访问实例
func NewDictRepository(db *gorm.DB) DictRepository {
	return &dictRepository{
		db: db,
	}
}

// === 字典类型相关方法 ===

// CreateType 创建字典类型
func (r *dictRepository) CreateType(dictType *model.DictType) error {
	return r.db.Create(dictType).Error
}

// GetTypeByID 根据ID获取字典类型
func (r *dictRepository) GetTypeByID(id uint64) (*model.DictType, error) {
	var dictType model.DictType
	err := r.db.First(&dictType, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dict type not found")
		}
		return nil, err
	}
	return &dictType, nil
}

// GetTypeByType 根据类型编码获取字典类型
func (r *dictRepository) GetTypeByType(typeCode string) (*model.DictType, error) {
	var dictType model.DictType
	err := r.db.Where("type = ?", typeCode).First(&dictType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dict type not found")
		}
		return nil, err
	}
	return &dictType, nil
}

// UpdateType 更新字典类型
func (r *dictRepository) UpdateType(dictType *model.DictType) error {
	return r.db.Save(dictType).Error
}

// DeleteType 删除字典类型（软删除）
func (r *dictRepository) DeleteType(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除该类型下的所有字典数据
		var dictType model.DictType
		if err := tx.First(&dictType, id).Error; err != nil {
			return err
		}

		// 删除字典数据
		if err := tx.Where("dict_type = ?", dictType.Type).Delete(&model.DictData{}).Error; err != nil {
			return err
		}

		// 删除字典类型
		return tx.Delete(&model.DictType{}, id).Error
	})
}

// GetTypeList 获取字典类型列表（分页）
func (r *dictRepository) GetTypeList(page, pageSize int, status int, name string) ([]*model.DictType, int64, error) {
	var dictTypes []*model.DictType
	var total int64

	query := r.db.Model(&model.DictType{})

	// 状态筛选
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 名称模糊搜索
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&dictTypes).Error
	if err != nil {
		return nil, 0, err
	}

	return dictTypes, total, nil
}

// TypeExists 检查字典类型编码是否存在
func (r *dictRepository) TypeExists(typeCode string) (bool, error) {
	var count int64
	err := r.db.Model(&model.DictType{}).
		Where("type = ?", typeCode).
		Count(&count).Error
	return count > 0, err
}

// === 字典数据相关方法 ===

// CreateData 创建字典数据
func (r *dictRepository) CreateData(dictData *model.DictData) error {
	return r.db.Create(dictData).Error
}

// GetDataByID 根据ID获取字典数据
func (r *dictRepository) GetDataByID(id uint64) (*model.DictData, error) {
	var dictData model.DictData
	err := r.db.First(&dictData, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dict data not found")
		}
		return nil, err
	}
	return &dictData, nil
}

// UpdateData 更新字典数据
func (r *dictRepository) UpdateData(dictData *model.DictData) error {
	return r.db.Save(dictData).Error
}

// DeleteData 删除字典数据（软删除）
func (r *dictRepository) DeleteData(id uint64) error {
	return r.db.Delete(&model.DictData{}, id).Error
}

// GetDataList 获取字典数据列表（分页）
func (r *dictRepository) GetDataList(dictType string, page, pageSize int, status int, label string) ([]*model.DictData, int64, error) {
	var dictData []*model.DictData
	var total int64

	query := r.db.Model(&model.DictData{}).Where("dict_type = ?", dictType)

	// 状态筛选
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 标签模糊搜索
	if label != "" {
		query = query.Where("label LIKE ?", "%"+label+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("sort_order ASC, created_at DESC").
		Find(&dictData).Error
	if err != nil {
		return nil, 0, err
	}

	return dictData, total, nil
}

// GetDataByType 根据字典类型获取所有数据
func (r *dictRepository) GetDataByType(dictType string) ([]*model.DictData, error) {
	var dictData []*model.DictData
	err := r.db.Where("dict_type = ?", dictType).
		Order("sort_order ASC, created_at ASC").
		Find(&dictData).Error
	return dictData, err
}

// GetEnabledDataByType 根据字典类型获取启用的数据
func (r *dictRepository) GetEnabledDataByType(dictType string) ([]*model.DictData, error) {
	var dictData []*model.DictData
	err := r.db.Where("dict_type = ? AND status = ?", dictType, 1).
		Order("sort_order ASC, created_at ASC").
		Find(&dictData).Error
	return dictData, err
}

// DataExists 检查字典数据是否存在
func (r *dictRepository) DataExists(dictType, value string) (bool, error) {
	var count int64
	err := r.db.Model(&model.DictData{}).
		Where("dict_type = ? AND value = ?", dictType, value).
		Count(&count).Error
	return count > 0, err
}

// GetDataByValue 根据字典类型和值获取字典数据
func (r *dictRepository) GetDataByValue(dictType, value string) (*model.DictData, error) {
	var dictData model.DictData
	err := r.db.Where("dict_type = ? AND value = ?", dictType, value).First(&dictData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dict data not found")
		}
		return nil, err
	}
	return &dictData, nil
}
