import http from "./axios";

export const userCardCurrencyCreate = (data) => {
  return http.request({
    url: '/user_card_currency/create',
    method: 'post',
    data
  })
}

export const userCardCurrencyUpdateBalance = (data) => {
  return http.request({
    url: '/user_card_currency/update/balance',
    method: 'post',
    data
  })
}

export const userCardCurrencyDelete = (data) => {
  return http.request({
    url: '/user_card_currency/delete',
    method: 'post',
    data
  })
}
