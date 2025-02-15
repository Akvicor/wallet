import http from "./axios";

export const typePeriodTypeFind = (params) => {
  return http.request({
    url: '/type/period/type',
    method: 'get',
    params
  })
}

export const typeTransactionTypeFind = (params) => {
  return http.request({
    url: '/type/transaction/type',
    method: 'get',
    params
  })
}

export const typeWalletTypeFind = (params) => {
  return http.request({
    url: '/type/wallet/type',
    method: 'get',
    params
  })
}

export const typeWalletPartitionAverageTypeFind = (params) => {
  return http.request({
    url: '/type/wallet_partition/average',
    method: 'get',
    params
  })
}