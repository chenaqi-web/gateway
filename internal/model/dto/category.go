package dto

import "backend/gateway/internal/client/rpc/core-rpc/categorypb"

// ---------- 实体 ----------

type CategoryType struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID       uint64 `json:"id"`
	ParentID uint64 `json:"parentID"`
	Name     string `json:"name"`
}

// ---------- Type 请求 / 响应 ----------

type CreateTypeRequest struct {
	Name string `json:"name" binding:"required"`
}

type DeleteTypeRequest struct {
	ID uint64 `json:"id" binding:"required"`
}

type ListTypesResponse struct {
	Types []*CategoryType `json:"types"`
}

// ---------- Category 请求 / 响应 ----------

type CreateCategoryRequest struct {
	ParentID uint64 `json:"parentID" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type DeleteCategoryRequest struct {
	ID uint64 `json:"id" binding:"required"`
}

type ListCategoriesRequest struct {
	ParentID uint64 `json:"parentID" binding:"required"`
}

type ListCategoriesResponse struct {
	Categories []*Category `json:"categories"`
}

type CategoryBoolResponse struct {
	Success bool `json:"success"`
}

// ---------- 转换 ----------

func ToCategoryType(item *categorypb.CategoryType) *CategoryType {
	if item == nil {
		return nil
	}
	return &CategoryType{ID: item.GetId(), Name: item.GetName()}
}

func ToCategoryTypes(items []*categorypb.CategoryType) []*CategoryType {
	list := make([]*CategoryType, 0, len(items))
	for _, item := range items {
		list = append(list, ToCategoryType(item))
	}
	return list
}

func ToCategory(item *categorypb.Category) *Category {
	if item == nil {
		return nil
	}
	return &Category{
		ID:       item.GetId(),
		ParentID: item.GetParentID(),
		Name:     item.GetName(),
	}
}

func ToCategories(items []*categorypb.Category) []*Category {
	list := make([]*Category, 0, len(items))
	for _, item := range items {
		list = append(list, ToCategory(item))
	}
	return list
}

func ToListTypesResponse(resp *categorypb.ListTypesResponse) *ListTypesResponse {
	if resp == nil {
		return &ListTypesResponse{Types: []*CategoryType{}}
	}
	return &ListTypesResponse{Types: ToCategoryTypes(resp.GetTypes())}
}

func ToListCategoriesResponse(resp *categorypb.ListCategoriesResponse) *ListCategoriesResponse {
	if resp == nil {
		return &ListCategoriesResponse{Categories: []*Category{}}
	}
	return &ListCategoriesResponse{Categories: ToCategories(resp.GetCategories())}
}

func ToCategoryBoolResponse(success bool) *CategoryBoolResponse {
	return &CategoryBoolResponse{Success: success}
}
