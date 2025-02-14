import { dayjs } from "element-plus"

export const formatTime = (time: string | number) => {
    let timestamp = Number(time)
    if (!isNaN(timestamp)) {
        // 判断时间戳是否是以毫秒为单位
        if (timestamp.toString().length <= 10) {
            // 如果长度小于等于10，说明是以秒为单位的时间戳
            timestamp *= 1000
        }
    } else {
        console.warn('无效的时间戳:', time)
        return '' // 返回空字符串或其他默认值
    }
    return dayjs(timestamp).format('YYYY-MM-DD')
} 