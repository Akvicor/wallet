import http from "./axios";


export const currencyCreate = (data) => {
  return http.request({
    url: '/currency/create',
    method: 'post',
    data
  })
}

export const currencyFind = (params) => {
  return http.request({
    url: '/currency/find',
    method: 'get',
    params
  })
}

export const currencyUpdate = (data) => {
  return http.request({
    url: '/currency/update',
    method: 'post',
    data
  })
}

export const currencyDelete = (data) => {
  return http.request({
    url: '/currency/delete',
    method: 'post',
    data
  })
}
