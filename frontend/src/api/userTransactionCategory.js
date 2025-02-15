import http from "./axios";


export const userTransactionCategoryCreate = (data) => {
  return http.request({
    url: '/user_transaction_category/create',
    method: 'post',
    data
  })
}

export const userTransactionCategoryFind = (params) => {
  return http.request({
    url: '/user_transaction_category/find',
    method: 'get',
    params
  })
}

export const userTransactionCategoryUpdate = (data) => {
  return http.request({
    url: '/user_transaction_category/update',
    method: 'post',
    data
  })
}

export const userTransactionCategoryUpdateSequence = (data) => {
  return http.request({
    url: '/user_transaction_category/update/sequence',
    method: 'post',
    data
  })
}

export const userTransactionCategoryDelete = (data) => {
  return http.request({
    url: '/user_transaction_category/delete',
    method: 'post',
    data
  })
}
