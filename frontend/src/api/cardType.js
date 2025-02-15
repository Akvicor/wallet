import http from "./axios";


export const cardTypeCreate = (data) => {
  return http.request({
    url: '/card_type/create',
    method: 'post',
    data
  })
}

export const cardTypeFind = (params) => {
  return http.request({
    url: '/card_type/find',
    method: 'get',
    params
  })
}

export const cardTypeUpdate = (data) => {
  return http.request({
    url: '/card_type/update',
    method: 'post',
    data
  })
}

export const cardTypeDelete = (data) => {
  return http.request({
    url: '/card_type/delete',
    method: 'post',
    data
  })
}
