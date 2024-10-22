import { dayjs } from "element-plus"

export const formatTime = (time: string) => {
    return dayjs(time).format('YYYY-MM-DD')
}