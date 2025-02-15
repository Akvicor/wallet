import React, {useEffect, useState} from 'react'
import {
  Button,
  Form,
  Input,
  Table,
  Popconfirm,
  Modal,
  message,
  Tag,
  Select,
  InputNumber,
  DatePicker, Flex,
} from 'antd'
import {
  userCardCreate,
  userCardFind,
  userCardUpdate,
  userCardDisable,
  userCardEnable,
  userCardValidRequest, userCardValidInput, userCardValidCancel, userCardUpdateSequence
} from "../../api/userCard";
import {ColorButtonProvider} from "../../theme/button";
import './userCard.css'
import dayjs from "dayjs";
import {currencyFind} from "../../api/currency";
import {bankFind} from "../../api/bank";
import {cardTypeFind} from "../../api/cardType";
import {
  userCardCurrencyCreate,
  userCardCurrencyDelete,
  userCardCurrencyUpdateBalance
} from "../../api/userCardCurrency";
import {useSelector} from "react-redux";

const UserCard = () => {
  const mode = useSelector(state => state.mode)
  const [searchKeyword, setSearchKeyword] = useState({
    search: ''
  })
  const [updateOrder, setUpdateOrder] = useState(false)
  const [showDisabled, setShowDisabled] = useState(false)
  const [verifyCodeInputOpen, setVerifyCodeInputOpen] = useState('');
  const [bankData, setBankData] = useState([])
  const [currencyData, setCurrencyData] = useState([])
  const [cardTypeData, setCardTypeData] = useState([])
  const [tableData, setTableData] = useState([])
  const [tableFilterData, setTableFilterData] = useState([])
  const [inputUserCardDataAction, setInputUserCardDataAction] = useState('close');
  const [form] = Form.useForm()
  const [inputUserCardCurrencyDataAction, setInputUserCardCurrencyDataAction] = useState('close');
  const [formCurrency] = Form.useForm()

  useEffect(() => {
    bankFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        data.data.forEach((item) => {
          options.push({
            value: item.id,
            label: item.name
          })
        })
        setBankData(options)
      }
    })
  }, [])
  useEffect(() => {
    currencyFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        data.data.forEach((item) => {
          options.push({
            value: item.id,
            label: item.name
          })
        })
        setCurrencyData(options)
      }
    })
  }, [])
  useEffect(() => {
    cardTypeFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        data.data.forEach((item) => {
          options.push({
            value: item.id,
            label: item.name
          })
        })
        setCardTypeData(options)
      }
    })
  }, [])

  const updateTable = ({search}) => {
    userCardFind({search: search, all: true}).then(({data}) => {
      if (data.code === 0) {
        setTableData(data.data)
        setTableFilterData(JSON.parse(JSON.stringify(data.data)).filter(item => item.disabled === 0))
      }
    })
  }

  const getTableData = () => {
    if (showDisabled) {
      return tableData
    } else {
      return tableFilterData
    }
  }

  const updateSequence = (id, target) => {
    userCardUpdateSequence({id: id, target: target}).then(({data}) => {
      if (data.code === 0) {
        updateTable(searchKeyword)
      }
    })
  }

  const showVerifyCodeInputModal = (key) => {
    if (key === 'cvv') {
      setVerifyCodeInputOpen(key)
    } else if (key === 'password') {
      setVerifyCodeInputOpen(key)
    }
  };
  const handleVerifyCodeInputOk = () => {
    setVerifyCodeInputOpen('');
  };
  const handleVerifyCodeInputCancel = () => {
    setVerifyCodeInputOpen('');
  };

  const handleDisableUserCard = (disable, id) => {
    if (disable) {
      userCardDisable({id: id}).then(({data}) => {
        if (data.code !== 0) {
          message.open({
            type: 'warning',
            content: '停用失败: ' + data.msg
          })
          return
        }
        message.open({
          type: 'success',
          content: '停用成功'
        })
        updateTable(searchKeyword)
      })
    } else {
      userCardEnable({id: id}).then(({data}) => {
        if (data.code !== 0) {
          message.open({
            type: 'warning',
            content: '启用失败' + data.msg
          })
          return
        }
        message.open({
          type: 'success',
          content: '启用成功'
        })
        updateTable(searchKeyword)
      })
    }
  }
  const handleDeleteCurrency = (id) => {
    userCardCurrencyDelete({id: id}).then(({data}) => {
      if (data.code !== 0) {
        message.open({
          type: 'warning',
          content: '删除失败: ' + data.msg
        })
        return
      }
      message.open({
        type: 'success',
        content: '删除成功'
      })
      updateTable(searchKeyword)
    })
  }
  const handleVerifyCodeRequest = (key, method) => {
    userCardValidRequest({key: key, method: method}).then(({data}) => {
      if (data.code !== 0) {
        message.open({
          type: 'warning',
          content: '验证码发送失败: ' + data.msg
        })
        return
      }
      message.open({
        type: 'success',
        content: '验证码发送成功'
      })
    })
  }
  const handleVerifyCodeInput = (input) => {
    userCardValidInput(input).then(({data}) => {
      if (data.code !== 0) {
        message.open({
          type: 'warning',
          content: '验证失败: ' + data.msg
        })
        return
      }
      message.open({
        type: 'success',
        content: '验证成功'
      })
      updateTable(searchKeyword)
      setVerifyCodeInputOpen('')
    })
  }
  const handleVerifyCodeCancel = (key) => {
    userCardValidCancel({key: key}).then(({data}) => {
      if (data.code !== 0) {
        message.open({
          type: 'warning',
          content: '取消失败: ' + data.msg
        })
        return
      }
      message.open({
        type: 'success',
        content: '取消成功'
      })
      updateTable(searchKeyword)
      setVerifyCodeInputOpen('')
    })
  }
  const handleSearchUserCard = ({keyword}) => {
    setSearchKeyword({
      search: keyword
    })
  }
  const handleInputUserCardDataShow = (action, data) => {
    if (data) {
      const cloneData = JSON.parse(JSON.stringify(data))
      cloneData.type_id = cloneData.types.map(item => {
        return item.id
      })
      cloneData.currency_id = cloneData.currency.map(item => {
        return item.id
      })
      cloneData.exp_date = dayjs(cloneData.exp_date * 1000)
      cloneData.cvv = '---'
      cloneData.password = '------'
      form.setFieldsValue(cloneData)
    } else {
      const cloneData = {}
      cloneData.statement_closing_day = 0
      cloneData.payment_due_day = 0
      cloneData.cvv = '000'
      cloneData.limit = 0
      cloneData.fee = '0'
      form.setFieldsValue(cloneData)
    }
    setInputUserCardDataAction(action)
  }
  const handleInputCurrencyDataShow = (action, data) => {
    if (data) {
      const cloneData = JSON.parse(JSON.stringify(data))
      formCurrency.setFieldsValue(cloneData)
    }
    setInputUserCardCurrencyDataAction(action)
  }

  const handleInputUserCardDataOk = () => {
    form.validateFields().then((input) => {
      input.exp_date = dayjs(input.exp_date).unix()
      if (inputUserCardDataAction === 'create') {
        userCardCreate(input).then(({data}) => {
          if (data.code !== 0) {
            message.open({
              type: 'warning',
              content: '创建失败: ' + data.msg
            })
            return
          }
          message.open({
            type: 'success',
            content: '创建成功'
          })
          updateTable(searchKeyword)
        })
      } else if (inputUserCardDataAction === 'update') {
        userCardUpdate(input).then(({data}) => {
          if (data.code !== 0) {
            message.open({
              type: 'warning',
              content: '更新失败: ' + data.msg
            })
            return
          }
          message.open({
            type: 'success',
            content: '更新成功'
          })
          updateTable(searchKeyword)
        })
      }
      setInputUserCardDataAction('close')
      form.resetFields()
    }).catch(() => {
      message.open({
        type: 'warning',
        content: '请检查输入'
      })
    })
  }

  const handleInputCurrencyDataOk = () => {
    formCurrency.validateFields().then((input) => {
      if (inputUserCardCurrencyDataAction === 'create') {
        userCardCurrencyCreate(input).then(({data}) => {
          if (data.code !== 0) {
            message.open({
              type: 'warning',
              content: '创建失败: ' + data.msg
            })
            return
          }
          message.open({
            type: 'success',
            content: '创建成功'
          })
          updateTable(searchKeyword)
          setInputUserCardCurrencyDataAction('close')
          formCurrency.resetFields()
        })
      } else if (inputUserCardCurrencyDataAction === 'update') {
        userCardCurrencyUpdateBalance(input).then(({data}) => {
          if (data.code !== 0) {
            message.open({
              type: 'warning',
              content: '更新失败: ' + data.msg
            })
            return
          }
          message.open({
            type: 'success',
            content: '更新成功'
          })
          updateTable(searchKeyword)
          setInputUserCardCurrencyDataAction('close')
          formCurrency.resetFields()
        })
      }
    }).catch(() => {
      message.open({
        type: 'warning',
        content: '请检查输入'
      })
    })
  }

  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
    }, {
      title: '卡号',
      dataIndex: 'number'
    }, {
      title: '银行',
      dataIndex: 'bank',
      render: (bank) => {
        return bank.name
      }
    }, {
      title: '类型',
      dataIndex: 'types',
      render: (types) => {
        return (
          <>
            {
              types.map(item => {
                return (
                  <div key={'div' + item.name}>
                    <Tag color='green' key={item.name}>
                      {item.name}
                    </Tag>
                  </div>
                )
              })
            }
          </>
        );
      }
    }, {
      title: '货币余额',
      dataIndex: 'id',
      render: (id) => {
        return (
          <div>
            <ColorButtonProvider color="green">
              <Button type="default" style={{marginRight: '5px'}}
                      onClick={() => handleInputCurrencyDataShow('create', {user_card_id: id})}>创建</Button>
            </ColorButtonProvider>
          </div>
        );
      }
    }, {
      title: '账单日',
      dataIndex: 'statement_closing_day',
      render: (statement_closing_day) => {
        if (statement_closing_day === 0) {
          return '-'
        }
        return statement_closing_day
      }
    }, {
      title: '还款日',
      dataIndex: 'payment_due_day',
      render: (payment_due_day) => {
        if (payment_due_day === 0) {
          return '-'
        }
        return payment_due_day
      }
    }, {
      title: '过期时间',
      dataIndex: 'exp_date',
      render: (exp_date) => {
        return dayjs(exp_date * 1000).format('MM/YY')
      }
    }, {
      title: 'CVV',
      dataIndex: 'cvv',
      render: (cvv) => {
        return <Button type="text" onClick={(e) => {
          if (e.target.textContent === '000') {
            e.target.textContent = cvv
          } else {
            e.target.textContent = '000'
          }
        }}>000</Button>
      }
    }, {
      title: '密码',
      dataIndex: 'password',
      render: (password) => {
        return <Button type="text" onClick={(e) => {
          if (e.target.textContent === '000000') {
            e.target.textContent = password
          } else {
            e.target.textContent = '000000'
          }
        }}>000000</Button>
      }
    }, {
      title: '管理费',
      dataIndex: 'fee',
      align: 'right',
      render: (fee) => {
        if (fee === '0') {
          return '-'
        }
        return fee
      }
    }, {
      title: '限制',
      dataIndex: 'limit',
      align: 'right',
      render: (limit) => {
        if (limit === '0') {
          return '-'
        }
        return limit
      }
    }, {
      title: '货币',
      key: 'currency',
      render: (rowData) => {
        let currency = rowData.currency
        return (
          <>
            {
              currency.map(item => {
                let colour = 'geekblue'
                if (item.currency_id === rowData.master_currency_id) {
                  colour = 'red'
                }
                let balance = ' (' + item.balance + item.currency.symbol + ')'
                if (rowData.hide_balance) {
                  balance = ''
                }
                return (
                  <div key={'div' + item.currency.name}>
                    <Tag color={colour} key={item.currency.name}>
                      {item.currency.name + balance}
                    </Tag>
                  </div>
                )
              })
            }
          </>
        );
      }
    }, {
      title: '描述',
      dataIndex: 'description'
    }, {
      title: '状态',
      dataIndex: 'disabled',
      render: (disabled) => {
        let color
        let status
        if (disabled === 0) {
          color = 'green';
          status = 'active'
        } else {
          color = 'red';
          status = 'disabled'
        }
        return (
          <>
            <Tag color={color} key={disabled}>
              {status}
            </Tag>
            {
              disabled !== 0 &&
              <Tag color={color} key={'time' + disabled}>
                {dayjs(disabled * 1000).format('YYYY-MM-DD HH:mm:ss')}
              </Tag>
            }
          </>
        );
      }
    }, {
      title: '操作',
      render: (rowData) => {
        let disableStatus = rowData.disabled === 0
        let disableMsg = disableStatus ? '停用' : '启用'
        return (
          <Flex style={{width: '100%'}} justify='flex-start' align='flex-start'>
            <Button type="primary" style={{marginRight: '5px'}}
                    onClick={() => handleInputUserCardDataShow('update', rowData)}>编辑</Button>
            <ColorButtonProvider danger={disableStatus} color="green">
              <Popconfirm
                title={disableMsg + '银行卡'}
                description={"你确定" + disableMsg + "银行卡?"}
                onConfirm={() => handleDisableUserCard(disableStatus, rowData.id)}
                onCancel={() => {
                }}
                okText="确认"
                cancelText="取消"
              >
                <Button danger={disableStatus} type="primary">{disableMsg}</Button>
              </Popconfirm>
            </ColorButtonProvider>
          </Flex>
        );
      }
    }
  ]
  if (updateOrder) {
    columns.push({
      title: '排序',
      render: (rowData) => {
        return (
          <Flex style={{width: '100%'}} justify='flex-start' align='flex-start'>
            <Button type="default" style={{marginRight: '5px'}}>{rowData.sequence}</Button>
            <Button type="default" style={{marginRight: '5px'}}
                    onClick={() => updateSequence(rowData.id, rowData.sequence - 1)}>上</Button>
            <Button type="default" onClick={() => updateSequence(rowData.id, rowData.sequence + 1)}>下</Button>
          </Flex>
        );
      }
    })
  }
  const columnsExpand = (row) => {
    const columns = [
      {
        title: '货币',
        dataIndex: 'currency',
        render: (currency) => {
          return currency.name
        }
      }, {
        title: '余额',
        dataIndex: 'balance',
        render: (balance) => {
          if (row.hide_balance) {
            return (
              <span onClick={(e) => {
                if (e.target.textContent === '???') {
                  e.target.textContent = balance
                } else {
                  e.target.textContent = '???'
                }
              }}>???</span>
            )
          }
          return balance
        }
      }, {
        title: '操作',
        render: (rowData) => {
          return (
            <Flex style={{width: '100%'}} justify='flex-start' align='flex-start'>
              <Button type="default" style={{marginRight: '5px', display: 'none'}}
                      onClick={() => handleInputCurrencyDataShow('update', rowData)} hidden>编辑</Button>
              <Popconfirm
                title={'删除货币余额'}
                description={"你确定删除货币余额?"}
                onConfirm={() => handleDeleteCurrency(rowData.id)}
                onCancel={() => {
                }}
                okText="确认"
                cancelText="取消"
              >
                <Button danger type="default">删除</Button>
              </Popconfirm>
            </Flex>
          );
        }
      }
    ]
    row.currency.forEach((item, index) => {
      item.key = index
    })
    return <Table columns={columns} dataSource={row.currency} pagination={false} size="small"/>
  }

  useEffect(() => {
    updateTable(searchKeyword)
  }, [searchKeyword]);

  return (
    <div>
      {
        mode.isWide ? (
          <Flex style={{width: '100%', marginBottom: '15px'}} justify='space-between' align='center'>
            <div>
              <Button type="primary" onClick={() => handleInputUserCardDataShow('create')}>创建</Button>
              <Button type="default" onClick={() => showVerifyCodeInputModal('cvv')}
                      style={{marginLeft: 5}}>CVV</Button>
              <Button type="default" onClick={() => showVerifyCodeInputModal('password')}
                      style={{marginLeft: 5}}>密码</Button>
              <Button type="default" onClick={() => setUpdateOrder(!updateOrder)} style={{marginLeft: 5}}>排序</Button>
              <Button type="default" onClick={() => setShowDisabled(!showDisabled)}
                      style={{marginLeft: 5}}>隐藏</Button>
            </div>
            <Form
              layout="inline"
              onFinish={handleSearchUserCard}
            >
              <Form.Item name="keyword">
                <Input placeholder='请输入关键词'/>
              </Form.Item>
              <Form.Item>
                <Button htmlType='submit' type='primary'>搜索</Button>
              </Form.Item>
            </Form>
          </Flex>
        ) : (
          <div style={{marginBottom: 15}}>
            <Flex style={{width: '100%', marginBottom: 6}} justify='flex-start' align='flex-start'>
              <Button type="primary" onClick={() => handleInputUserCardDataShow('create')}>创建</Button>
              <Button type="default" onClick={() => showVerifyCodeInputModal('cvv')}
                      style={{marginLeft: 5}}>CVV</Button>
              <Button type="default" onClick={() => showVerifyCodeInputModal('password')}
                      style={{marginLeft: 5}}>密码</Button>
              <Button type="default" onClick={() => setUpdateOrder(!updateOrder)}
                      style={{marginLeft: 5}}>排序</Button>
              <Button type="default" onClick={() => setShowDisabled(!showDisabled)}
                      style={{marginLeft: 5}}>隐藏</Button>
            </Flex>
            <Form
              layout="inline"
              onFinish={handleSearchUserCard}
            >
              <Form.Item name="keyword">
                <Input placeholder='请输入关键词'/>
              </Form.Item>
              <Form.Item>
                <Button htmlType='submit' type='primary'>搜索</Button>
              </Form.Item>
            </Form>
          </div>
        )
      }
      <Table
        columns={columns}
        scroll={{
          x: 'max-content',
        }}
        expandable={{
          expandedRowRender: columnsExpand,
          expandIcon: ({expanded, onExpand, record}) => {
            return expanded ? (
              <Button type="default" onClick={e => onExpand(record, e)}>收起</Button>
            ) : (
              <Button type="default" onClick={e => onExpand(record, e)}>展开</Button>
            )
          }
        }}
        dataSource={getTableData()}
        pagination={{
          pageSizeOptions: [10, 15, 20, 50, 100],
          responsive: true,
          showQuickJumper: true,
          showSizeChanger: true
        }}
        rowKey={'id'}
      />
      <Modal
        title={inputUserCardDataAction === 'create' ? '创建' : inputUserCardDataAction === 'update' ? '编辑' : 'Unknown'}
        open={inputUserCardDataAction !== 'close'}
        onOk={handleInputUserCardDataOk}
        onCancel={() => {
          setInputUserCardDataAction('close');
          form.resetFields()
        }}
        okText="确定"
        cancelText="取消"
      >
        <Form
          form={form}
          labelCol={{
            span: 6
          }}
          wrapperCol={{
            span: 18
          }}
        >
          {
            inputUserCardDataAction === 'update' &&
            <Form.Item
              name="id"
              hidden
            >
              <Input/>
            </Form.Item>
          }
          <Form.Item
            label="名称"
            name="name"
            rules={[
              {
                required: true,
                message: '请输入名称'
              }
            ]}
          >
            <Input placeholder={'请输入名称'}/>
          </Form.Item>
          <Form.Item
            label="卡号"
            name="number"
            rules={[
              {
                required: true,
                message: '请输入卡号'
              }
            ]}
          >
            <Input placeholder={'请输入卡号'}/>
          </Form.Item>
          <Form.Item
            label="银行"
            name="bank_id"
            rules={[
              {
                required: true,
                message: '请输入银行'
              }
            ]}
          >
            <Select
              placeholder="银行"
              onChange={() => {
              }}
              allowClear
              options={bankData}
            />
          </Form.Item>
          <Form.Item
            label="类型"
            name="type_id"
            rules={[
              {
                required: true,
                message: '请输入类型'
              }
            ]}
          >
            <Select
              mode="multiple"
              placeholder="类型"
              onChange={() => {
              }}
              showSearch={false}
              allowClear
              options={cardTypeData}
            />
          </Form.Item>
          <Form.Item
            label="账单日"
            name="statement_closing_day"
            rules={[
              {
                required: true,
                message: '请输入账单日'
              }
            ]}
          >
            <InputNumber style={{width: '100%'}} placeholder={'请输入账单日'}/>
          </Form.Item>
          <Form.Item
            label="还款日"
            name="payment_due_day"
            rules={[
              {
                required: true,
                message: '请输入还款日'
              }
            ]}
          >
            <InputNumber style={{width: '100%'}} placeholder={'请输入还款日'}/>
          </Form.Item>
          <Form.Item
            label="过期时间"
            name="exp_date"
            rules={[
              {
                required: true,
                message: '请输入过期时间'
              }
            ]}
          >
            <DatePicker format="MM/YYYY" style={{width: '100%'}} picker="month"/>
          </Form.Item>
          <Form.Item
            label="CVV"
            name="cvv"
            rules={[
              {
                required: true,
                message: '请输入CVV'
              }
            ]}
          >
            <Input placeholder={'请输入CVV'}/>
          </Form.Item>
          <Form.Item
            label="密码"
            name="password"
            rules={[
              {
                required: true,
                message: '请输入密码'
              }
            ]}
          >
            <Input placeholder={'请输入密码'}/>
          </Form.Item>
          <Form.Item
            label="限制"
            name="limit"
            rules={[
              {
                required: true,
                message: '请输入限制'
              }
            ]}
          >
            <Input placeholder={'请输入限制'}/>
          </Form.Item>
          <Form.Item
            label="管理费"
            name="fee"
            rules={[
              {
                required: true,
                message: '请输入管理费'
              }
            ]}
          >
            <Input placeholder={'请输入管理费'}/>
          </Form.Item>
          <Form.Item
            label="主货币"
            name="master_currency_id"
            rules={[
              {
                required: true,
                message: '请输入主货币'
              }
            ]}
          >
            <Select
              placeholder="主货币"
              onChange={() => {
              }}
              allowClear
              options={currencyData}
            />
          </Form.Item>
          {
            inputUserCardDataAction === 'create' &&
            <Form.Item
              label="货币"
              name="currency_id"
              rules={[
                {
                  required: true,
                  message: '请输入货币'
                }
              ]}
            >
              <Select
                mode="multiple"
                placeholder="货币"
                onChange={() => {
                }}
                showSearch={false}
                allowClear
                options={currencyData}
              />
            </Form.Item>
          }
          <Form.Item
            label="隐藏余额"
            name="hide_balance"
            rules={[
              {
                required: true,
                message: '请选择是否隐藏余额'
              }
            ]}
          >
            <Select
              placeholder="隐藏余额"
              onChange={() => {
              }}
              allowClear
              options={[{value: true, label: '是'}, {value: false, label: '否'}]}
            />
          </Form.Item>
          <Form.Item
            label="描述"
            name="description"
          >
            <Input placeholder={'请输入描述'}/>
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title={inputUserCardCurrencyDataAction === 'create' ? '创建' : inputUserCardCurrencyDataAction === 'update' ? '编辑' : 'Unknown'}
        open={inputUserCardCurrencyDataAction !== 'close'}
        onOk={handleInputCurrencyDataOk}
        onCancel={() => {
          setInputUserCardCurrencyDataAction('close');
          formCurrency.resetFields()
        }}
        okText="确定"
        cancelText="取消"
      >
        <Form
          form={formCurrency}
          labelCol={{
            span: 6
          }}
          wrapperCol={{
            span: 18
          }}
        >
          {
            inputUserCardCurrencyDataAction === 'update' &&
            <Form.Item
              name="id"
              hidden
            >
              <Input/>
            </Form.Item>
          }
          {
            inputUserCardCurrencyDataAction === 'update' &&
            <Form.Item
              label="余额"
              name="balance"
              rules={[
                {
                  required: true,
                  message: '请输入余额'
                }
              ]}
            >
              <Input placeholder={'请输入余额'}/>
            </Form.Item>
          }
          {
            inputUserCardCurrencyDataAction === 'create' &&
            <Form.Item
              name="user_card_id"
              hidden
            >
              <Input/>
            </Form.Item>
          }
          {
            inputUserCardCurrencyDataAction === 'create' &&
            <Form.Item
              label="货币类型"
              name="currency_id"
              rules={[
                {
                  required: true,
                  message: '请输入货币类型'
                }
              ]}
            >
              <Select
                placeholder="货币类型"
                onChange={() => {
                }}
                allowClear
                options={currencyData}
              />
            </Form.Item>
          }
        </Form>
      </Modal>
      <Modal title="请输入绑定手机收到的验证码" open={verifyCodeInputOpen !== ''} onOk={handleVerifyCodeInputOk}
             onCancel={handleVerifyCodeInputCancel}>
        <Form
          layout="inline"
          fields={[{
            name: 'key',
            value: verifyCodeInputOpen
          }]}
          onFinish={handleVerifyCodeInput}
        >
          <Flex style={{width: '100%', marginBottom: 6}} justify='flex-start' align='flex-start'>
            <Form.Item name="key" hidden>
              <Input/>
            </Form.Item>
            <Form.Item name="verify_code">
              <Input placeholder='请输入验证码'/>
            </Form.Item>
            <Form.Item>
              <Button type="primary" ghost onClick={() => {
                handleVerifyCodeRequest(verifyCodeInputOpen, 'phone')
              }}>短信验证</Button>
            </Form.Item>
            <Form.Item>
              <Button type="primary" ghost onClick={() => {
                handleVerifyCodeRequest(verifyCodeInputOpen, 'mail')
              }}>邮箱验证</Button>
            </Form.Item>
          </Flex>
          <Flex style={{width: '100%'}} justify='flex-start' align='flex-start'>
            <Form.Item>
              <Button htmlType='submit' type='primary'>验证</Button>
            </Form.Item>
            <Form.Item>
              <Button type='default' onClick={() => {
                handleVerifyCodeCancel(verifyCodeInputOpen)
              }}>取消</Button>
            </Form.Item>
          </Flex>
        </Form>
      </Modal>
    </div>
  )
}

export default UserCard
