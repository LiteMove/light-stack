package service

import (
	"errors"
	"fmt"
	"github.com/LiteMove/light-stack/internal/modules/system/repository"

	"github.com/LiteMove/light-stack/internal/model"
)

// DictService 字典服务接口
type DictService interface {
	// 字典类型相关方法
	CreateType(dictType *model.DictType) error
	GetType(id uint64) (*model.DictType, error)
	GetTypeByCode(typeCode string) (*model.DictType, error)
	UpdateType(dictType *model.DictType) error
	DeleteType(id uint64) error
	GetTypeList(page, pageSize int, status int, name string) ([]*model.DictType, int64, error)

	// 字典数据相关方法
	CreateData(dictData *model.DictData) error
	GetData(id uint64) (*model.DictData, error)
	UpdateData(dictData *model.DictData) error
	DeleteData(id uint64) error
	GetDataList(dictType string, page, pageSize int, status int, label string) ([]*model.DictData, int64, error)
	GetDataByType(dictType string) ([]*model.DictData, error)
	GetEnabledDataByType(dictType string) ([]*model.DictData, error)

	// 批量操作
	BatchUpdateDataStatus(ids []uint64, status int) error
	BatchDeleteData(ids []uint64) error

	// 验证方法
	ValidateTypeCode(typeCode string, excludeID ...uint64) error
	ValidateDataValue(dictType, value string, excludeID ...uint64) error

	// 获取字典选项（用于前端下拉框）
	GetDictOptions(dictType string) ([]*DictOption, error)
}

// DictOption 字典选项结构（用于前端下拉框）
type DictOption struct {
	Label     string `json:"label"`
	Value     string `json:"value"`
	CssClass  string `json:"cssClass,omitempty"`
	ListClass string `json:"listClass,omitempty"`
	IsDefault bool   `json:"isDefault"`
}

// dictService 字典服务实现
type dictService struct {
	dictRepo repository.DictRepository
}

// NewDictService 创建字典服务
func NewDictService(dictRepo repository.DictRepository) DictService {
	return &dictService{
		dictRepo: dictRepo,
	}
}

// === 字典类型相关方法 ===

// CreateType 创建字典类型
func (s *dictService) CreateType(dictType *model.DictType) error {
	// 验证类型编码是否重复
	if err := s.ValidateTypeCode(dictType.Type); err != nil {
		return err
	}

	return s.dictRepo.CreateType(dictType)
}

// GetType 根据ID获取字典类型
func (s *dictService) GetType(id uint64) (*model.DictType, error) {
	return s.dictRepo.GetTypeByID(id)
}

// GetTypeByCode 根据类型编码获取字典类型
func (s *dictService) GetTypeByCode(typeCode string) (*model.DictType, error) {
	return s.dictRepo.GetTypeByType(typeCode)
}

// UpdateType 更新字典类型
func (s *dictService) UpdateType(dictType *model.DictType) error {
	// 验证类型编码是否重复（排除当前记录）
	if err := s.ValidateTypeCode(dictType.Type, dictType.ID); err != nil {
		return err
	}

	// 获取原有的字典类型
	oldType, err := s.dictRepo.GetTypeByID(dictType.ID)
	if err != nil {
		return fmt.Errorf("获取原字典类型失败: %w", err)
	}

	// 如果类型编码发生变化，需要更新相关的字典数据
	if oldType.Type != dictType.Type {
		// 获取该类型的所有字典数据
		dataList, err := s.dictRepo.GetDataByType(oldType.Type)
		if err != nil {
			return fmt.Errorf("获取字典数据失败: %w", err)
		}

		// 更新所有字典数据的类型
		for _, data := range dataList {
			data.DictType = dictType.Type
			if err := s.dictRepo.UpdateData(data); err != nil {
				return fmt.Errorf("更新字典数据类型失败: %w", err)
			}
		}
	}

	return s.dictRepo.UpdateType(dictType)
}

// DeleteType 删除字典类型
func (s *dictService) DeleteType(id uint64) error {
	// 检查是否存在关联的字典数据
	dictType, err := s.dictRepo.GetTypeByID(id)
	if err != nil {
		return fmt.Errorf("获取字典类型失败: %w", err)
	}

	dataList, err := s.dictRepo.GetDataByType(dictType.Type)
	if err != nil {
		return fmt.Errorf("检查关联字典数据失败: %w", err)
	}

	if len(dataList) > 0 {
		return errors.New("该字典类型下存在字典数据，无法删除")
	}

	return s.dictRepo.DeleteType(id)
}

// GetTypeList 获取字典类型列表
func (s *dictService) GetTypeList(page, pageSize int, status int, name string) ([]*model.DictType, int64, error) {
	return s.dictRepo.GetTypeList(page, pageSize, status, name)
}

// === 字典数据相关方法 ===

// CreateData 创建字典数据
func (s *dictService) CreateData(dictData *model.DictData) error {
	// 验证字典类型是否存在
	_, err := s.dictRepo.GetTypeByType(dictData.DictType)
	if err != nil {
		return errors.New("字典类型不存在")
	}

	// 验证字典值是否重复
	if err := s.ValidateDataValue(dictData.DictType, dictData.Value); err != nil {
		return err
	}

	return s.dictRepo.CreateData(dictData)
}

// GetData 根据ID获取字典数据
func (s *dictService) GetData(id uint64) (*model.DictData, error) {
	return s.dictRepo.GetDataByID(id)
}

// UpdateData 更新字典数据
func (s *dictService) UpdateData(dictData *model.DictData) error {
	// 验证字典类型是否存在
	_, err := s.dictRepo.GetTypeByType(dictData.DictType)
	if err != nil {
		return errors.New("字典类型不存在")
	}

	// 验证字典值是否重复（排除当前记录）
	if err := s.ValidateDataValue(dictData.DictType, dictData.Value, dictData.ID); err != nil {
		return err
	}

	return s.dictRepo.UpdateData(dictData)
}

// DeleteData 删除字典数据
func (s *dictService) DeleteData(id uint64) error {
	return s.dictRepo.DeleteData(id)
}

// GetDataList 获取字典数据列表
func (s *dictService) GetDataList(dictType string, page, pageSize int, status int, label string) ([]*model.DictData, int64, error) {
	return s.dictRepo.GetDataList(dictType, page, pageSize, status, label)
}

// GetDataByType 根据字典类型获取所有数据
func (s *dictService) GetDataByType(dictType string) ([]*model.DictData, error) {
	return s.dictRepo.GetDataByType(dictType)
}

// GetEnabledDataByType 根据字典类型获取启用的数据
func (s *dictService) GetEnabledDataByType(dictType string) ([]*model.DictData, error) {
	return s.dictRepo.GetEnabledDataByType(dictType)
}

// === 批量操作方法 ===

// BatchUpdateDataStatus 批量更新字典数据状态
func (s *dictService) BatchUpdateDataStatus(ids []uint64, status int) error {
	for _, id := range ids {
		dictData, err := s.dictRepo.GetDataByID(id)
		if err != nil {
			return fmt.Errorf("获取字典数据失败 ID=%d: %w", id, err)
		}
		dictData.Status = status
		if err := s.dictRepo.UpdateData(dictData); err != nil {
			return fmt.Errorf("更新字典数据状态失败 ID=%d: %w", id, err)
		}
	}
	return nil
}

// BatchDeleteData 批量删除字典数据
func (s *dictService) BatchDeleteData(ids []uint64) error {
	for _, id := range ids {
		if err := s.dictRepo.DeleteData(id); err != nil {
			return fmt.Errorf("删除字典数据失败 ID=%d: %w", id, err)
		}
	}
	return nil
}

// === 验证方法 ===

// ValidateTypeCode 验证字典类型编码是否重复
func (s *dictService) ValidateTypeCode(typeCode string, excludeID ...uint64) error {
	exists, err := s.dictRepo.TypeExists(typeCode)
	if err != nil {
		return fmt.Errorf("检查字典类型编码是否存在失败: %w", err)
	}

	if exists {
		// 如果提供了excludeID，需要检查是否是当前记录
		if len(excludeID) > 0 {
			existingType, err := s.dictRepo.GetTypeByType(typeCode)
			if err == nil && existingType.ID == excludeID[0] {
				return nil // 是当前记录，允许
			}
		}
		return errors.New("字典类型编码已存在")
	}

	return nil
}

// ValidateDataValue 验证字典数据值是否重复
func (s *dictService) ValidateDataValue(dictType, value string, excludeID ...uint64) error {
	exists, err := s.dictRepo.DataExists(dictType, value)
	if err != nil {
		return fmt.Errorf("检查字典数据值是否存在失败: %w", err)
	}

	if exists {
		// 如果提供了excludeID，需要检查是否是当前记录
		if len(excludeID) > 0 {
			existingData, err := s.dictRepo.GetDataByValue(dictType, value)
			if err == nil && existingData.ID == excludeID[0] {
				return nil // 是当前记录，允许
			}
		}
		return errors.New("字典数据值已存在")
	}

	return nil
}

// === 前端下拉框相关方法 ===

// GetDictOptions 获取字典选项（用于前端下拉框）
func (s *dictService) GetDictOptions(dictType string) ([]*DictOption, error) {
	dataList, err := s.GetEnabledDataByType(dictType)
	if err != nil {
		return nil, err
	}

	options := make([]*DictOption, 0, len(dataList))
	for _, data := range dataList {
		options = append(options, &DictOption{
			Label:     data.Label,
			Value:     data.Value,
			CssClass:  data.CssClass,
			ListClass: data.ListClass,
			IsDefault: data.IsDefault,
		})
	}

	return options, nil
}
