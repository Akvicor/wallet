import http from "./axios";


export const userCardCreate = (data) => {
  return http.request({
    url: '/user_card/create',
    method: 'post',
    data
  })
}

export const userCardFind = (params) => {
  return http.request({
    url: '/user_card/find',
    method: 'get',
    params
  })
}

export const userCardValidRequest = (data) => {
  return http.request({
    url: '/user_card/valid/request',
    method: 'post',
    data
  })
}

export const userCardValidInput = (data) => {
  return http.request({
    url: '/user_card/valid/input',
    method: 'post',
    data
  })
}

export const userCardValidCancel = (data) => {
  return http.request({
    url: '/user_card/valid/cancel',
    method: 'post',
    data
  })
}

export const userCardUpdate = (data) => {
  return http.request({
    url: '/user_card/update',
    method: 'post',
    data
  })
}

export const userCardUpdateSequence = (data) => {
  return http.request({
    url: '/user_card/update/sequence',
    method: 'post',
    data
  })
}

export const userCardDisable = (data) => {
  return http.request({
    url: '/user_card/disable',
    method: 'post',
    data
  })
}

export const userCardEnable = (data) => {
  return http.request({
    url: '/user_card/enable',
    method: 'post',
    data
  })
}
