package models

import (
	"github.com/jinzhu/gorm"
)

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// ExistTagByName checks if there is a tag with the same name
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ? AND deleted_on = ? ", name, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

// AddTag Add a Tag
func AddTag(name string, state int, createdBy string) error {
	tag := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}
	if err := db.Create(&tag).Error; err != nil {
		return err
	}

	return nil
}

// GetTags gets a list of tags based on paging and constraints
func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var (
		tags []Tag
		err  error
	)

	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Find(&tags).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

// GetTagTotal counts the total number of tags based on the constraint
func GetTagTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// ExistTagByID determines whether a Tag exists based on the ID
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

// DeleteTag delete a tag
func DeleteTag(id int) error {
	if err := db.Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		return err
	}

	return nil
}

// EditTag modify a single tag
func EditTag(id int, data interface{}) error {
	if err := db.Model(&Tag{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// CleanAllTag clear all tag
func CleanAllTag() (bool, error) {
	if err := db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{}).Error; err != nil {
		return false, err
	}

	return true, nil
}

//type Tag struct {
//	Model
//
//	Name       string `json:"name"`
//	CreatedBy  string `json:"created_by"`
//	ModifiedBy string `json:"modified_by"`
//	State      int    `json:"state"`
//}
//
//// GetTags 我们创建了一个Tag struct{}，用于Gorm的使用。并给予了附属属性json，这样子在c.JSON的时候就会自动转换格式，非常的便利
//// 编写标签列表的models逻辑
//func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
//	// maps 里面有state name 是多少的信息, 通过c.query得到
//	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
//
//	return
//}
//
//func GetTagTotal(maps interface{}) (count int) {
//	db.Model(&Tag{}).Where(maps).Count(&count)
//
//	return
//}
//
//// ExistTagByName AddTag 编写新增标签的models逻辑
//func ExistTagByName(name string) bool {
//	var tag Tag
//	db.Select("id").Where("name = ?", name).First(&tag)
//	if tag.ID > 0 {
//		return true
//	}
//
//	return false
//}
//
//func AddTag(name string, state int, createdBy string) bool {
//	db.Create(&Tag{
//		Name:      name,
//		State:     state,
//		CreatedBy: createdBy,
//	})
//
//	return true
//}
//
//func ExistTagByID(id int) bool {
//	var tag Tag
//	db.Select("id").Where("id = ?", id).First(&tag)
//	if tag.ID > 0 {
//		return true
//	}
//
//	return false
//}
//
//func DeleteTag(id int) bool {
//	db.Where("id = ?", id).Delete(&Tag{})
//
//	return true
//}
//
//func EditTag(id int, data interface{}) bool {
//	db.Model(&Tag{}).Where("id = ?", id).Updates(data)
//
//	return true
//}

// BeforeCreate BeforeUpdate 这属于gorm的Callbacks，可以将回调方法定义为模型结构的指针，在创建、更新、查询、删除时将被调用，如果任何回调返回错误，gorm将停止未来操作并回滚所有更改。
// gorm所支持的回调方法：
// 创建：BeforeSave、BeforeCreate、AfterCreate、AfterSave
// 更新：BeforeSave、BeforeUpdate、AfterUpdate、AfterSave
// 删除：BeforeDelete、AfterDelete
// 查询：AfterFind
// ***已经有了回调函数 所以不再需要了***
//func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
//	scope.SetColumn("CreatedOn", time.Now().Unix())
//
//	return nil
//}
//
//func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
//	scope.SetColumn("ModifiedOn", time.Now().Unix())
//
//	return nil
//}

// CleanAllTag 用于corn包相关操作
//func CleanAllTag() bool {
//	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{})
//
//	return true
//}
