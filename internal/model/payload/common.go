package payload

import (
    "math"
)

type ResponseData struct {
    Data any `json:"data"`
}

type ResponseList struct {
    Total  int64          `json:"total"`
    Params any            `json:"params"`
    Meta   PaginationMeta `json:"meta"`
    Items  any            `json:"items"`
}

type PaginationMeta struct {
    CurrentPage int `json:"currentPage"`
    NextPage    int `json:"nextPage"`
    TotalPage   int `json:"totalPage"`
}

func Ok(data any) ResponseData {
    return ResponseData{Data: data}
}

func Paginated(data any, count int64) ResponseList {
    return ResponseList{
        Items: data,
    }
}

const (
    defaultLimit = 25
    maxLimit = 100
)

type PaginationFilter struct {
    Page   int    `json:"page" query:"page"`
    Limit  int    `json:"limit" query:"limit"`
    Search string `json:"search" query:"search"`
}

func (p *PaginationFilter) Normalize() {
    if p.Limit > maxLimit {
        p.Limit = maxLimit
    }

    if p.Limit == 0 {
        p.Limit = defaultLimit
    }

    if p.Page == 0 {
        p.Page = 1
    }
}

func (p *PaginationFilter) Paginate(items any, total int64) ResponseList {
    return ResponseList{
        Items: items,
        Total: total,
        Params: p,
        Meta: PaginationMeta{
            CurrentPage: p.Page,
            NextPage: p.Page + 1,
            TotalPage: func() int {
                if total == 0 {
                    return 1
                }

                return int(math.Ceil(float64(total) / float64(p.Limit)))
            }(),
        },
    }
}
