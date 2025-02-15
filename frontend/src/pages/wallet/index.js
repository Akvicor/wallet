import React, {useEffect, useState} from 'react'
import {Button, Form, Input, Table, Popconfirm, Modal, message, Tag, Select, Flex} from 'antd'
import {
  walletCreate,
  walletUpdate,
  walletDisable,
  walletEnable,
  walletFindNormal,
  walletUpdateSequence
} from "../../api/wallet";
import {ColorButtonProvider} from "../../theme/button";
import './wallet.css'
import dayjs from "dayjs";
import {currencyFind} from "../../api/currency";
import {typeWalletPartitionAverageTypeFind, typeWalletTypeFind} from "../../api/type";
import {userCardFind} from "../../api/userCard";
import {
  walletPartitionCreate,
  walletPartitionDisable,
  walletPartitionEnable,
  walletPartitionUpdate, walletPartitionUpdateSequence
} from "../../api/walletPartition";
import {
  PartitionAverageTypeMonth, PartitionAverageTypeNormal,
  PartitionAverageTypeQuarter,
  PartitionAverageTypeWeek,
  PartitionAverageTypeYear
} from "./type";
import {useSelector} from "react-redux";

dayjs.extend(require("dayjs/plugin/quarterOfYear"));
dayjs.extend(require("dayjs/plugin/isoWeek"));

const WalletType = [1, 2]

const Wallet = () => {
  const mode = useSelector(state => state.mode)
  const [searchKeyword, setSearchKeyword] = useState({
    search: ''
  })
  const [updateOrder, setUpdateOrder] = useState(false)
  const [showDisabled, setShowDisabled] = useState(false)
  const [userCardData, setUserCardData] = useState([])
  const [currencyData, setCurrencyData] = useState([])
  const [walletTypeData, setWalletTypeData] = useState([])
  const [walletPartitionAverageTypeData, setWalletPartitionAverageTypeData] = useState([])
  const [tableData, setTableData] = useState([])
  const [tableFilterData, setTableFilterData] = useState([])
  const [inputWalletDataAction, setInputWalletDataAction] = useState('close');
  const [form] = Form.useForm()
  const [inputPartitionDataAction, setInputPartitionDataAction] = useState('close');
  const [formPartition] = Form.useForm()


  useEffect(() => {
    userCardFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        data.data.forEach((item) => {
          let bank = item.bank
          let name = item.name
          if (bank) {
            name = '[' + bank.name + '] ' + name
          }
          options.push({
            value: item.id,
            label: name
          })
        })
        setUserCardData(options)
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
    typeWalletTypeFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        data.data.forEach((item) => {
          if (WalletType.indexOf(item.type) > -1) {
            options.push({
              value: item.type,
              label: item.name
            })
          }
        })
        setWalletTypeData(options)
      }
    })
  }, [])
  useEffect(() => {
    typeWalletPartitionAverageTypeFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        data.data.forEach((item) => {
          options.push({
            value: item.type,
            label: item.name
          })
        })
        setWalletPartitionAverageTypeData(options)
      }
    })
  }, [])

  const updateTable = ({search}) => {
    walletFindNormal({search: search, all: true, all_partition: true}).then(({data}) => {
      if (data.code === 0) {
        setTableData(data.data)
        setTableFilterData(JSON.parse(JSON.stringify(data.data)).filter(item => item.disabled === 0).map(item => {
          item.partition = item.partition.filter(part => part.disabled === 0)
          return item
        }))
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

  const updateWalletSequence = (id, target) => {
    walletUpdateSequence({id: id, target: target}).then(({data}) => {
      if (data.code === 0) {
        updateTable(searchKeyword)
      }
    })
  }

  const updateWalletPartitionSequence = (walletID, id, target) => {
    walletPartitionUpdateSequence({wallet_id: walletID, id: id, target: target}).then(({data}) => {
      if (data.code === 0) {
        updateTable(searchKeyword)
      }
    })
  }

  const handleDisableWallet = (disable, id) => {
    if (disable) {
      walletDisable({id: id}).then(({data}) => {
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
      walletEnable({id: id}).then(({data}) => {
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
  const handleDisablePartition = (disable, wallet_id, id) => {
    if (disable) {
      walletPartitionDisable({wallet_id: wallet_id, id: id}).then(({data}) => {
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
      walletPartitionEnable({wallet_id: wallet_id, id: id}).then(({data}) => {
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
  const handleSearchWallet = ({keyword}) => {
    setSearchKeyword({
      search: keyword
    })
  }
  const handleInputWalletDataShow = (action, data) => {
    if (data) {
      const cloneData = JSON.parse(JSON.stringify(data))
      form.setFieldsValue(cloneData)
    }
    setInputWalletDataAction(action)
  }
  const handleInputPartitionDataShow = (action, data) => {
    if (data) {
      const cloneData = JSON.parse(JSON.stringify(data))
      formPartition.setFieldsValue(cloneData)
    }
    setInputPartitionDataAction(action)
  }

  const handleInputWalletDataOk = () => {
    form.validateFields().then((input) => {
      if (inputWalletDataAction === 'create') {
        walletCreate(input).then(({data}) => {
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
          setInputWalletDataAction('close')
          form.resetFields()
        })
      } else if (inputWalletDataAction === 'update') {
        walletUpdate(input).then(({data}) => {
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
          setInputWalletDataAction('close')
          form.resetFields()
        })
      }
    }).catch(() => {
      message.open({
        type: 'warning',
        content: '请检查输入'
      })
    })
  }

  const handleInputPartitionDataOk = () => {
    formPartition.validateFields().then((input) => {
      if (inputPartitionDataAction === 'create') {
        walletPartitionCreate(input).then(({data}) => {
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
          setInputPartitionDataAction('close')
          formPartition.resetFields()
        })
      } else if (inputPartitionDataAction === 'update') {
        walletPartitionUpdate(input).then(({data}) => {
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
          setInputPartitionDataAction('close')
          formPartition.resetFields()
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
      title: '描述',
      dataIndex: 'description',
    }, {
      title: '划分',
      dataIndex: 'id',
      render: (id) => {
        return (
          <div>
            <ColorButtonProvider color="green">
              <Button type="default" style={{marginRight: '5px'}}
                      onClick={() => handleInputPartitionDataShow('create', {wallet_id: id})}>创建</Button>
            </ColorButtonProvider>
          </div>
        );
      }
    }, {
      title: '货币',
      key: 'partition',
      render: (rowData) => {
        return (
          <>
            {rowData.partition.map(item => {
              if (item.disabled !== 0) {
                return ''
              }
              let avg = ''
              if (item.average !== PartitionAverageTypeNormal) {
                let avgDaySum = 1
                if (item.average === PartitionAverageTypeWeek) {
                  avgDaySum = dayjs().endOf('isoWeek').diff(dayjs().endOf('day'), 'day') + 1
                } else if (item.average === PartitionAverageTypeMonth) {
                  avgDaySum = dayjs().endOf('month').diff(dayjs().endOf('day'), 'day') + 1
                } else if (item.average === PartitionAverageTypeQuarter) {
                  avgDaySum = dayjs().endOf('quarter').diff(dayjs().endOf('day'), 'day') + 1
                } else if (item.average === PartitionAverageTypeYear) {
                  avgDaySum = dayjs().endOf('year').diff(dayjs().endOf('day'), 'day') + 1
                }
                avg = '日均(' + (Math.round((parseFloat(item.balance) / avgDaySum) * 100) / 100) + item.currency.symbol + ')'
              }
              let balance = ' (' + item.balance + item.currency.symbol + ')'
              if (rowData.wallet_type === 2) {
                balance = ''
              }
              return (
                <div key={'partition' + item.id}>
                  <Tag color='blue' key={'partition' + item.id + item.name}>
                    {item.name + ' | ' + item.currency.name + balance + ' ' + avg}
                  </Tag>
                </div>
              )
            })}
          </>
        );
      }
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
                    onClick={() => handleInputWalletDataShow('update', rowData)}>编辑</Button>
            <ColorButtonProvider danger={disableStatus} color="green">
              <Popconfirm
                title={disableMsg + '钱包'}
                description={"你确定" + disableMsg + "钱包?"}
                onConfirm={() => handleDisableWallet(disableStatus, rowData.id)}
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
                    onClick={() => updateWalletSequence(rowData.id, rowData.sequence - 1)}>上</Button>
            <Button type="default" onClick={() => updateWalletSequence(rowData.id, rowData.sequence + 1)}>下</Button>
          </Flex>
        );
      }
    })
  }
  const columnsExpand = (row) => {
    const columns = [
      {
        title: '划分',
        dataIndex: 'name',
      }, {
        title: '货币',
        dataIndex: 'currency',
        render: (currency) => {
          return currency.name
        }
      }, {
        title: '余额',
        dataIndex: 'balance',
        align: 'right',
        render: (balance) => {
          if (row.wallet_type === 2) {
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
        title: '限额',
        dataIndex: 'limit',
        align: 'right',
        render: (limit) => {
          if (limit === '0') return "-"
          return limit
        }
      }, {
        title: '均分',
        dataIndex: 'average',
        render: (average) => {
          const result = walletPartitionAverageTypeData.find(avg => avg.value === average)
          if (result) return result.label
          return "未知"
        }
      }, {
        title: '绑定银行卡',
        dataIndex: 'card',
        render: (card) => {
          return card.name
        }
      }, {
        title: '描述',
        dataIndex: 'description',
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
            <div>
              <Tag color={color} key={disabled}>
                {status}
              </Tag>
              {
                disabled !== 0 &&
                <Tag color={color} key={'time' + disabled}>
                  {dayjs(disabled * 1000).format('YYYY-MM-DD HH:mm:ss')}
                </Tag>
              }
            </div>
          );
        }
      }, {
        title: '操作',
        render: (rowData) => {
          let disableStatus = rowData.disabled === 0
          let disableMsg = disableStatus ? '停用' : '启用'
          return (
            <Flex style={{width: '100%'}} justify='flex-start' align='flex-start'>
              <Button type="default" style={{marginRight: '5px'}}
                      onClick={() => handleInputPartitionDataShow('update', rowData)}>编辑</Button>
              <ColorButtonProvider danger={disableStatus} color="green">
                <Popconfirm
                  title={disableMsg + '划分'}
                  description={"你确定" + disableMsg + "划分?"}
                  onConfirm={() => handleDisablePartition(disableStatus, row.id, rowData.id)}
                  onCancel={() => {
                  }}
                  okText="确认"
                  cancelText="取消"
                >
                  <Button danger={disableStatus} type="default">{disableMsg}</Button>
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
                      onClick={() => updateWalletPartitionSequence(rowData.wallet_id, rowData.id, rowData.sequence - 1)}>上</Button>
              <Button type="default"
                      onClick={() => updateWalletPartitionSequence(rowData.wallet_id, rowData.id, rowData.sequence + 1)}>下</Button>
            </Flex>
          );
        }
      })
    }
    for (let i = 0; i < row.partition.length; i++) {
      row.partition[i].key = i
    }
    return <Table columns={columns} dataSource={row.partition} pagination={false} size="small"/>
  }

  useEffect(() => {
    updateTable(searchKeyword)
  }, [searchKeyword]);

  return (
    <div>
      {
        mode.isWide ? (
          <Flex style={{width: '100%', marginBottom: '15px'}} justify='space-between' align='center'>
            <div style={{marginBottom: 6}}>
              <Button type="primary" onClick={() => handleInputWalletDataShow('create')}
                      style={{marginRight: 6}}>创建</Button>
              <Button type="default" onClick={() => setUpdateOrder(!updateOrder)} style={{marginRight: 6}}>排序</Button>
              <Button type="default" onClick={() => setShowDisabled(!showDisabled)}>隐藏</Button>
            </div>
            <Form
              layout="inline"
              onFinish={handleSearchWallet}
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
          <div style={{marginBottom: '15px'}}>
            <div style={{marginBottom: 6}}>
              <Button type="primary" onClick={() => handleInputWalletDataShow('create')}
                      style={{marginRight: 6}}>创建</Button>
              <Button type="default" onClick={() => setUpdateOrder(!updateOrder)} style={{marginRight: 6}}>排序</Button>
              <Button type="default" onClick={() => setShowDisabled(!showDisabled)}>隐藏</Button>
            </div>
            <Form
              layout="inline"
              onFinish={handleSearchWallet}
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
        title={inputWalletDataAction === 'create' ? '创建' : inputWalletDataAction === 'update' ? '编辑' : 'Unknown'}
        open={inputWalletDataAction !== 'close'}
        onOk={handleInputWalletDataOk}
        onCancel={() => {
          setInputWalletDataAction('close');
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
            inputWalletDataAction === 'update' &&
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
            label="类型"
            name="wallet_type"
            rules={[
              {
                required: true,
                message: '请输入类型'
              }
            ]}
          >
            <Select
              placeholder="类型"
              onChange={() => {
              }}
              allowClear
              options={walletTypeData}
            />
          </Form.Item>
          {
            inputWalletDataAction !== 'update' &&
            <Form.Item
              label="银行卡"
              name="cards_id"
              rules={[
                {
                  required: true,
                  message: '请输入银行卡'
                }
              ]}
            >
              <Select
                mode="multiple"
                placeholder="银行卡"
                onChange={() => {
                }}
                showSearch={false}
                allowClear
                options={userCardData}
              />
            </Form.Item>
          }
          <Form.Item
            label="描述"
            name="description"
          >
            <Input placeholder={'请输入描述'}/>
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title={inputPartitionDataAction === 'create' ? '创建' : inputPartitionDataAction === 'update' ? '编辑' : 'Unknown'}
        open={inputPartitionDataAction !== 'close'}
        onOk={handleInputPartitionDataOk}
        onCancel={() => {
          setInputPartitionDataAction('close');
          formPartition.resetFields()
        }}
        okText="确定"
        cancelText="取消"
      >
        <Form
          form={formPartition}
          labelCol={{
            span: 6
          }}
          wrapperCol={{
            span: 18
          }}
        >
          {
            inputPartitionDataAction === 'update' &&
            <Form.Item
              name="id"
              hidden
            >
              <Input/>
            </Form.Item>
          }
          <Form.Item
            name="wallet_id"
            hidden
          >
            <Input/>
          </Form.Item>
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
          {
            inputPartitionDataAction === 'create' &&
            <Form.Item
              label="银行卡"
              name="card_id"
              rules={[
                {
                  required: true,
                  message: '请输入银行卡'
                }
              ]}
            >
              <Select
                placeholder="银行卡"
                onChange={() => {
                }}
                allowClear
                options={userCardData}
              />
            </Form.Item>
          }
          {
            inputPartitionDataAction === 'create' &&
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
          <Form.Item
            label="限额"
            name="limit"
            rules={[
              {
                required: true,
                message: '请输入限额'
              }
            ]}
          >
            <Input placeholder={'请输入限额'}/>
          </Form.Item>
          <Form.Item
            label="均分类型"
            name="average"
            rules={[
              {
                required: true,
                message: '请输入均分类型'
              }
            ]}
          >
            <Select
              placeholder="均分类型"
              onChange={() => {
              }}
              allowClear
              options={walletPartitionAverageTypeData}
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
    </div>
  )
}

export default Wallet
