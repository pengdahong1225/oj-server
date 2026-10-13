import request from '@/utils/request'

// 查询公告列表
export async function queryNoticeListService(params: API.QueryNoticeListParams) {
    return request('/notice/list', {
        method: 'GET',
        params: {
            ...params
        }
    })
}

// 查询公告详情
export async function queryNoticeDetailService(id: number) {
    return request('/notice/detail', {
        method: 'GET',
        params: { id }
    })
}

// 新增公告
export async function addNoticeService(form: API.AddNoticeForm) {
    return request('/notice/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        data: form
    })
}

// 更新公告
export async function updateNoticeService(form: API.UpdateNoticeForm) {
    return request('/notice/update', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        data: form
    })
}

// 删除公告
export async function deleteNoticeService(id: number) {
    return request('/notice/delete', {
        method: 'DELETE',
        params: { id }
    })
}

// 发布公告
export async function publishNoticeService(id: number) {
    return request('/notice/publish', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        data: { id }
    })
}

// 隐藏公告
export async function hideNoticeService(id: number) {
    return request('/notice/hide', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        data: { id }
    })
}
