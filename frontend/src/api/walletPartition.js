import http from "./axios";


export const walletPartitionCreate = (data) => {
  return http.request({
    url: '/wallet_partition/create',
    method: 'post',
    data
  })
}

export const walletPartitionUpdate = (data) => {
  return http.request({
    url: '/wallet_partition/update',
    method: 'post',
    data
  })
}

export const walletPartitionUpdateSequence = (data) => {
  return http.request({
    url: '/wallet_partition/update/sequence',
    method: 'post',
    data
  })
}

export const walletPartitionDisable = (data) => {
  return http.request({
    url: '/wallet_partition/disable',
    method: 'post',
    data
  })
}

export const walletPartitionEnable = (data) => {
  return http.request({
    url: '/wallet_partition/enable',
    method: 'post',
    data
  })
}
