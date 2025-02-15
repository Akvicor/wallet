import http from "./axios";


export const transactionCreate = (data) => {
  return http.request({
    url: '/user_transaction/create',
    method: 'post',
    data
  })
}

export const transactionFind = (data) => {
  return http.request({
    url: '/user_transaction/find',
    method: 'post',
    data
  })
}

export const transactionFindRange = (params) => {
  return http.request({
    url: '/user_transaction/find/range',
    method: 'get',
    params
  })
}

export const transactionUpdate = (data) => {
  return http.request({
    url: '/user_transaction/update',
    method: 'post',
    data
  })
}

export const transactionChecked = (data) => {
  return http.request({
    url: '/user_transaction/checked',
    method: 'post',
    data
  })
}

export const transactionDelete = (data) => {
  return http.request({
    url: '/user_transaction/delete',
    method: 'post',
    data
  })
}

export const transactionViewDay = (params) => {
  return http.request({
    url: '/user_transaction/view/day',
    method: 'get',
    params
  })
}

export const transactionViewMonth = (params) => {
  return http.request({
    url: '/user_transaction/view/month',
    method: 'get',
    params
  })
}

export const transactionViewYear = (params) => {
  return http.request({
    url: '/user_transaction/view/year',
    method: 'get',
    params
  })
}

export const transactionViewTotal = (params) => {
  return http.request({
    url: '/user_transaction/view/total',
    method: 'get',
    params
  })
}

export const transactionChart = (data) => {
  return http.request({
    url: '/user_transaction/chart',
    method: 'post',
    data
  })
}

export const transactionChartPie = (data) => {
  return http.request({
    url: '/user_transaction/chart/pie',
    method: 'post',
    data
  })
}
