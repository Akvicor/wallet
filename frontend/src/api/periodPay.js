import http from "./axios";


export const periodPayCreate = (data) => {
  return http.request({
    url: '/user_period_pay/create',
    method: 'post',
    data
  })
}

export const periodPayFind = (params) => {
  return http.request({
    url: '/user_period_pay/find',
    method: 'get',
    params
  })
}

export const periodPaySummary = (params) => {
  return http.request({
    url: '/user_period_pay/summary',
    method: 'get',
    params
  })
}

export const periodPayUpdate = (data) => {
  return http.request({
    url: '/user_period_pay/update',
    method: 'post',
    data
  })
}

export const periodPayDelete = (data) => {
  return http.request({
    url: '/user_period_pay/delete',
    method: 'post',
    data
  })
}

export const periodPayUpdateNext = (data) => {
  return http.request({
    url: '/user_period_pay/update/next',
    method: 'post',
    data
  })
}
