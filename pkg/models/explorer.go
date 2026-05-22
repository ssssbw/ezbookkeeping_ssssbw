package models

import "encoding/json"

// InsightsExplorer represents a saved insights explorer configuration
type InsightsExplorer struct {
	ExplorerId      int64  `xorm:"PK comment('探索器ID')"`
	Uid             int64  `xorm:"INDEX(IDX_insights_explorer_uid_deleted_order) NOT NULL comment('用户ID')"`
	Deleted         bool   `xorm:"INDEX(IDX_insights_explorer_uid_deleted_order) NOT NULL comment('是否删除')"`
	Name            string `xorm:"VARCHAR(64) NOT NULL comment('探索器名称')"`
	DisplayOrder    int32  `xorm:"INDEX(IDX_insights_explorer_uid_deleted_order) NOT NULL comment('显示排序')"`
	Data            string `xorm:"MEDIUMBLOB comment('配置数据JSON')"`
	Hidden          bool   `xorm:"NOT NULL comment('是否隐藏')"`
	CreatedUnixTime int64  `comment('创建时间')"`
	UpdatedUnixTime int64  `comment('更新时间')"`
	DeletedUnixTime int64  `comment('删除时间')"`
}

// InsightsExplorerCreateRequest represents all parameters of insights explorer creation request
type InsightsExplorerCreateRequest struct {
	Name            string         `json:"name" binding:"required,notBlank,max=64"`
	Data            map[string]any `json:"data" binding:"required"`
	ClientSessionId string         `json:"clientSessionId"`
}

// InsightsExplorerModifyRequest represents all parameters of insights explorer modification request
type InsightsExplorerModifyRequest struct {
	Id              int64          `json:"id,string" binding:"required,min=0"`
	Name            string         `json:"name" binding:"required,notBlank,max=64"`
	Data            map[string]any `json:"data" binding:"required"`
	Hidden          bool           `json:"hidden"`
	ClientSessionId string         `json:"clientSessionId"`
}

// InsightsExplorerGetRequest represents all parameters of insights explorer getting request
type InsightsExplorerGetRequest struct {
	Id int64 `form:"id,string" binding:"required,min=1"`
}

// InsightsExplorerHideRequest represents all parameters of insights explorer hiding request
type InsightsExplorerHideRequest struct {
	Id     int64 `json:"id,string" binding:"required,min=1"`
	Hidden bool  `json:"hidden"`
}

// InsightsExplorerMoveRequest represents all parameters of insights explorer moving request
type InsightsExplorerMoveRequest struct {
	NewDisplayOrders []*InsightsExplorerNewDisplayOrderRequest `json:"newDisplayOrders" binding:"required,min=1"`
}

// InsightsExplorerNewDisplayOrderRequest represents a data pair of id and display order
type InsightsExplorerNewDisplayOrderRequest struct {
	Id           int64 `json:"id,string" binding:"required,min=1"`
	DisplayOrder int32 `json:"displayOrder"`
}

// InsightsExplorerDeleteRequest represents all parameters of insights explorer deleting request
type InsightsExplorerDeleteRequest struct {
	Id int64 `json:"id,string" binding:"required,min=1"`
}

// InsightsExplorerInfoResponse represents a view-object of insights explorer info
type InsightsExplorerInfoResponse struct {
	Id           int64          `json:"id,string"`
	Name         string         `json:"name"`
	DisplayOrder int32          `json:"displayOrder"`
	Hidden       bool           `json:"hidden"`
	Data         map[string]any `json:"data,omitempty"`
}

// ToInsightsExplorerInfoResponse returns a view-object according to database model
func (a *InsightsExplorer) ToInsightsExplorerInfoResponse() (*InsightsExplorerInfoResponse, error) {
	var data map[string]any = nil

	if a.Data != "" {
		err := json.Unmarshal([]byte(a.Data), &data)

		if err != nil {
			return nil, err
		}
	}

	return &InsightsExplorerInfoResponse{
		Id:           a.ExplorerId,
		Name:         a.Name,
		DisplayOrder: a.DisplayOrder,
		Hidden:       a.Hidden,
		Data:         data,
	}, nil
}

// InsightsExplorerInfoResponseSlice represents the slice data structure of InsightsExplorerInfoResponse
type InsightsExplorerInfoResponseSlice []*InsightsExplorerInfoResponse

// Len returns the count of items
func (s InsightsExplorerInfoResponseSlice) Len() int {
	return len(s)
}

// Swap swaps two items
func (s InsightsExplorerInfoResponseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the first item is less than the second one
func (s InsightsExplorerInfoResponseSlice) Less(i, j int) bool {
	return s[i].DisplayOrder < s[j].DisplayOrder
}
