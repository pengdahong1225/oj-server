/**
 * @description token验证模块
 * 1.将base64字符串解码成json字符串
 * 2.将json字符串解码成json对象
 */

import store from '@/store'

// 判断token过期时间
// 一旦发现token过期了，就删除信息
export function checkExpire(token) {
	const expire = parseExpire(token)
	const now = new Date().getTime() / 1000
	if (expire - now < 0) {
		console.log('token过期:', expire)
		store.dispatch('user/logout')
		return false
	}
	return true
}


function parseExpire(token) {
	const payload = token.split('.')[1]
	const decodedPayload = atob(payload)
	const obj = JSON.parse(decodedPayload)

	return obj.exp
}
