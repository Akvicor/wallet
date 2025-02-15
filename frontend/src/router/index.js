import {createBrowserRouter, Navigate} from 'react-router-dom'
import Main from '../pages/main'
import Home from '../pages/home'
import User from "../pages/user";
import Bank from "../pages/bank";
import Login from '../pages/login'
import {AdminAuth, RouterAuth, ViewerAuth} from './routerAuth'
import Currency from "../pages/currency";
import CardType from "../pages/cardType";
import UserExchangeRate from "../pages/userExchangeRate";
import UserCard from "../pages/userCard";
import Wallet from "../pages/wallet";
import UserTransactionCategory from "../pages/userTransactionCategory";
import Transaction from "../pages/transaction";
import Wishlist from "../pages/wishlist";
import Debt from "../pages/debt";
import PeriodPay from "../pages/periodPay";
import UserCenter from "../pages/userCenter";


const routes = [
  {
    path: '/',
    Component: Main,
    children: [
      {
        path: '/',
        element: (<RouterAuth><Navigate to="/home" replace/></RouterAuth>)
      },
      {
        path: 'home',
        Component: Home
      },
      {
        path: 'transaction',
        element: (<ViewerAuth><Transaction/></ViewerAuth>)
      },
      {
        path: 'wallet',
        element: (<ViewerAuth><Wallet/></ViewerAuth>)
      },
      {
        path: 'user-card',
        element: (<ViewerAuth><UserCard/></ViewerAuth>)
      },
      {
        path: 'wishlist',
        element: (<ViewerAuth><Wishlist/></ViewerAuth>)
      },
      {
        path: 'period-pay',
        element: (<ViewerAuth><PeriodPay/></ViewerAuth>)
      },
      {
        path: 'debt',
        element: (<ViewerAuth><Debt/></ViewerAuth>)
      },
      {
        path: 'admin',
        children: [
          {
            path: 'user-center',
            element: (<ViewerAuth><UserCenter/></ViewerAuth>)
          },
          {
            path: 'user',
            element: (<AdminAuth><User/></AdminAuth>)
          },
          {
            path: 'bank',
            element: (<AdminAuth><Bank/></AdminAuth>)
          },
          {
            path: 'currency',
            element: (<AdminAuth><Currency/></AdminAuth>)
          },
          {
            path: 'card-type',
            element: (<AdminAuth><CardType/></AdminAuth>)
          },
          {
            path: 'user-exchange-rate',
            element: (<ViewerAuth><UserExchangeRate/></ViewerAuth>)
          },
          {
            path: 'user-transaction-category',
            element: (<ViewerAuth><UserTransactionCategory/></ViewerAuth>)
          }
        ]
      }
    ]
  }, {
    path: '/login',
    Component: Login
  }
]

export default createBrowserRouter(routes)
