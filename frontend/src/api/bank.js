import http from "./axios";


export const bankCreate = (data) => {
  return http.request({
    url: '/bank/create',
    method: 'post',
    data
  })
}

export const bankFind = (params) => {
  return http.request({
    url: '/bank/find',
    method: 'get',
    params
  })
}

export const bankUpdate = (data) => {
  return http.request({
    url: '/bank/update',
    method: 'post',
    data
  })
}

export const bankDelete = (data) => {
  return http.request({
    url: '/bank/delete',
    method: 'post',
    data
  })
}
