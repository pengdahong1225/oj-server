function setInLocalStorage (key, value) {
  localStorage.setItem(key, JSON.stringify(value))
}
function getFromLocalStorage (key) {
  return JSON.parse(localStorage.getItem(key))
}

// keys
const InfoKey = 'oj_user_info'
const TokenKey = 'oj_user_token'

// 存储用户信息
export const setUserInfo = (obj) => {
  setInLocalStorage(InfoKey, obj)
}

// 存储token
export const setToken = (obj) => {
  setInLocalStorage(TokenKey, obj)
}

// 读取用户信息
export const getUserInfo = () => {
  return getFromLocalStorage(InfoKey)
}

// 读取token
export const getToken = () => {
  return getFromLocalStorage(TokenKey)
}

// 移除个人信息
export const removeInfo = () => {
  localStorage.removeItem(InfoKey)
  localStorage.removeItem(TokenKey)
}
