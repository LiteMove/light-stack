package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Provider 提供者函数类型
type Provider func(c *DIContainer) interface{}

// DIContainer 简化版依赖注入容器
type DIContainer struct {
	mu        sync.RWMutex
	instances map[reflect.Type]interface{}
	providers map[reflect.Type]Provider
}

// NewDI 创建新的依赖注入容器
func NewDI() *DIContainer {
	return &DIContainer{
		instances: make(map[reflect.Type]interface{}),
		providers: make(map[reflect.Type]Provider),
	}
}

// RegisterSingleton 注册单例服务
func (c *DIContainer) RegisterSingleton(provider Provider, target interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	targetType := reflect.TypeOf(target)
	// 移除指针包装，直接使用类型
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	c.providers[targetType] = provider
}

// RegisterInstance 注册已有实例
func (c *DIContainer) RegisterInstance(instance interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	instanceType := reflect.TypeOf(instance)
	c.instances[instanceType] = instance
}

// Get 获取服务实例（类型安全）
func Get[T any](c *DIContainer) T {
	var zero T
	targetType := reflect.TypeOf(zero)

	instance := c.resolve(targetType)
	if instance == nil {
		return zero
	}

	if result, ok := instance.(T); ok {
		return result
	}

	return zero
}

// resolve 解析依赖
func (c *DIContainer) resolve(targetType reflect.Type) interface{} {
	c.mu.RLock()

	// 首先检查是否已有实例
	if instance, exists := c.instances[targetType]; exists {
		c.mu.RUnlock()
		return instance
	}

	// 检查是否有提供者
	provider, exists := c.providers[targetType]
	c.mu.RUnlock()

	if !exists {
		return nil
	}

	// 创建实例并缓存（单例模式）
	c.mu.Lock()
	defer c.mu.Unlock()

	// 双重检查
	if instance, exists := c.instances[targetType]; exists {
		return instance
	}

	instance := provider(c)
	c.instances[targetType] = instance
	return instance
}

// MustGet 获取服务实例，不存在会panic
func MustGet[T any](c *DIContainer) T {
	result := Get[T](c)
	var zero T

	if reflect.DeepEqual(result, zero) {
		panic(fmt.Sprintf("service of type %T is not registered", zero))
	}

	return result
}
