import http from "./axios";


export const walletCreate = (data) => {
  return http.request({
    url: '/wallet/create',
    method: 'post',
    data
  })
}

export const walletFind = (params) => {
  return http.request({
    url: '/wallet/find',
    method: 'get',
    params
  })
}

export const walletFindNormal = (params) => {
  return http.request({
    url: '/wallet/find/normal',
    method: 'get',
    params
  })
}

export const walletFindDebt = (params) => {
  return http.request({
    url: '/wallet/find/debt',
    method: 'get',
    params
  })
}

export const walletFindWishlist = (params) => {
  return http.request({
    url: '/wallet/find/wishlist',
    method: 'get',
    params
  })
}

export const walletUpdate = (data) => {
  return http.request({
    url: '/wallet/update',
    method: 'post',
    data
  })
}

export const walletUpdateSequence = (data) => {
  return http.request({
    url: '/wallet/update/sequence',
    method: 'post',
    data
  })
}

export const walletDisable = (data) => {
  return http.request({
    url: '/wallet/disable',
    method: 'post',
    data
  })
}

export const walletEnable = (data) => {
  return http.request({
    url: '/wallet/enable',
    method: 'post',
    data
  })
}
