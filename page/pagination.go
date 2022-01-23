package page

import (
	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
	Total    int `json:"total" form:"total"`
	Offset   int `json:"-"`
}

// Parse pagination from gin ctx
func (p *Pagination) Parse(ctx *gin.Context) {
	p.Init()
	if err := ctx.ShouldBind(p); err != nil {
		return
	}
	p.Page = max(1, p.Page)
	p.PageSize = max(1, p.PageSize)
	p.Offset = (p.Page - 1) * p.PageSize
}

// init pagination
func (p *Pagination) Init() {
	p.Page = 1
	p.PageSize = 10
	p.Offset = 0
}

// Set Page
func (p *Pagination) SetPage(page int) {
	p.Page = max(1, page)
	p.PageSize = max(1, p.PageSize)
	p.Offset = (p.Page - 1) * p.PageSize
}

// Set Page Size
func (p *Pagination) SetPageSize(pageSize int) {
	p.Page = max(1, p.Page)
	p.PageSize = max(1, pageSize)
	p.Offset = (p.Page - 1) * p.PageSize
}

// Set Total Num
func (p *Pagination) SetTotal(total int) {
	p.Total = total
}

// Get offset
func (p *Pagination) GetOffset() int {
	if p.Page <= 0 {
		p.Init()
	}
	p.Offset = (p.Page - 1) * p.PageSize
	return p.Offset
}

// int max
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
