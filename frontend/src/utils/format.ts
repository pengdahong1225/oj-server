import { dayjs } from "element-plus"

export const formatTime = (time: string | number) => {
  if (time === null || time === undefined || time === '') return ''

  // 纯数字（或数字字符串）按时间戳处理
  if (typeof time === 'number' || /^\d+$/.test(String(time))) {
    let timestamp = Number(time)
    // 秒级时间戳转毫秒
    if (String(timestamp).length <= 10) {
      timestamp *= 1000
    }
    return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
  }

  // 其他情况交给 dayjs 解析日期字符串
  const d = dayjs(time)
  return d.isValid() ? d.format('YYYY-MM-DD HH:mm:ss') : ''
}