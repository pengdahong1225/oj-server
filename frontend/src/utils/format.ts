import { dayjs } from "element-plus"

export const formatTime = (time: string) => {
    let timestamp = Number(time);
    if (!isNaN(timestamp)) {
        timestamp *= 1000;
    }
    return dayjs(timestamp).format('YYYY-MM-DD')
} 