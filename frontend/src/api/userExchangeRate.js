import http from "./axios";


export const userExchangeRateCreate = (data) => {
  return http.request({
    url: '/user_exchange_rate/create',
    method: 'post',
    data
  })
}

export const userExchangeRateFind = (params) => {
  return http.request({
    url: '/user_exchange_rate/find',
    method: 'get',
    params
  })
}

export const userExchangeRateUpdate = (data) => {
  return http.request({
    url: '/user_exchange_rate/update',
    method: 'post',
    data
  })
}

export const userExchangeRateDelete = (data) => {
  return http.request({
    url: '/user_exchange_rate/delete',
    method: 'post',
    data
  })
}
