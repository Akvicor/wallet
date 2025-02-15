const config = [
  {
    path: '/home',
    name: 'home',
    label: '首页',
    role: 'viewer',
    icon: 'CalendarOutlined',
    url: '/home/index'
  },
  {
    path: '/transaction',
    name: 'transaction',
    label: '交易',
    role: 'viewer',
    icon: 'FileSearchOutlined',
    url: '/transaction/index'
  },
  {
    path: '/wallet',
    name: 'wallet',
    label: '钱包',
    role: 'viewer',
    icon: 'WalletOutlined',
    url: '/wallet/index'
  },
  {
    path: '/user-card',
    name: 'user-card',
    label: '银行卡',
    role: 'viewer',
    icon: 'CreditCardOutlined',
    url: '/user-card/index'
  },
  {
    path: '/wishlist',
    name: 'wishlist',
    label: '愿望单',
    role: 'viewer',
    icon: 'GiftOutlined',
    url: '/wishlist/index'
  },
  {
    path: '/period-pay',
    name: 'period-pay',
    label: '订阅',
    role: 'viewer',
    icon: 'HourglassOutlined',
    url: '/period-pay/index'
  },
  {
    path: '/debt',
    name: 'debt',
    label: '债务',
    role: 'viewer',
    icon: 'HourglassOutlined',
    url: '/debt/index'
  },
  {
    path: '/admin',
    label: '管理',
    role: 'viewer',
    icon: 'SettingOutlined',
    children: [
      {
        path: '/admin/user-center',
        name: 'user-center',
        label: '用户中心',
        role: 'viewer',
        icon: 'UserOutlined'
      },
      {
        path: '/admin/user',
        name: 'user',
        label: '用户',
        role: 'admin',
        icon: 'UsergroupAddOutlined'
      },
      {
        path: '/admin/bank',
        name: 'bank',
        label: '银行',
        role: 'admin',
        icon: 'BankOutlined'
      },
      {
        path: '/admin/currency',
        name: 'currency',
        label: '货币类型',
        role: 'admin',
        icon: 'PayCircleOutlined'
      },
      {
        path: '/admin/card-type',
        name: 'card-type',
        label: '银行卡类型',
        role: 'admin',
        icon: 'CreditCardOutlined'
      },
      {
        path: '/admin/user-exchange-rate',
        name: 'user-exchange-rate',
        label: '汇率',
        role: 'viewer',
        icon: 'TransactionOutlined'
      },
      {
        path: '/admin/user-transaction-category',
        name: 'user-transaction-category',
        label: '交易分类',
        role: 'viewer',
        icon: 'FlagOutlined'
      }
    ]
  }
]

const AdminMenu = JSON.parse(JSON.stringify(config)).filter(() => {
  return true
})

const UserMenu = JSON.parse(JSON.stringify(config)).filter(val => val.role !== 'admin').map(item => {
  if (item.children) {
    item.children = item.children.filter(val => {
      return val.role !== 'admin'
    })
  }
  return item
})

const ViewerMenu = JSON.parse(JSON.stringify(config)).filter(val => val.role !== 'admin' && val.role !== 'user').map(item => {
  if (item.children) {
    item.children = item.children.filter(val => {
      return val.role !== 'admin' && val.role !== 'user'
    })
  }
  return item
})

export {AdminMenu, UserMenu, ViewerMenu}
