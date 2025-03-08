import React, {useEffect, useState} from 'react'
import {Button, Form, Input, Table, Popconfirm, Modal, message, Select, TreeSelect, DatePicker, Tag, Flex} from 'antd'
import './transaction.css'
import {currencyFind} from "../../api/currency";
import {userTransactionCategoryFind} from "../../api/userTransactionCategory";
import {typeTransactionTypeFind} from "../../api/type";
import {
  transactionChecked,
  transactionCreate,
  transactionDelete,
  transactionFind,
  transactionUpdate
} from "../../api/transaction";
import dayjs from "dayjs";
import {walletFind} from "../../api/wallet";
import {
  TransactionTypeAutoExchange,
  TransactionTypeAutoTransfer,
  TransactionTypeExpense,
  TransactionTypeIncome
} from "./type";
import {useSelector} from "react-redux";

const Transaction = () => {
  const mode = useSelector(state => state.mode)
  const [searchKeyword, setSearchKeyword] = useState({
    search: '',
    from_partition_id: [],
    to_partition_id: [],
    from_currency_id: [],
    to_currency_id: [],
    category_id: [],
    transaction_type: []
  })
  const [selectedTransactionType, setSelectedTransactionType] = useState(0)
  const [selectedCategoryData, setSelectedCategoryData] = useState([])
  const [walletData, setWalletData] = useState([])
  const [currencyData, setCurrencyData] = useState([])
  const [categoryData, setCategoryData] = useState([])
  const [transactionTypeData, setTransactionTypeData] = useState([])
  const [tableData, setTableData] = useState([])
  const [tablePaginationData, setTablePaginationData] = useState({current: 1, pageSize: 10})
  const [tablePaginationRespData, setTablePaginationRespData] = useState({total: 0})
  const [tableFilterPartitionData, setTableFilterPartitionData] = useState([])
  const [tableFilterCurrencyData, setTableFilterCurrencyData] = useState([])
  const [tableFilterCategoryData, setTableFilterCategoryData] = useState([])
  const [tableFilterTransactionTypeData, setTableFilterTransactionTypeData] = useState([])
  const [tableLoadingData, setTableLoadingData] = useState(false)
  const [inputTransactionDataAction, setInputTransactionDataAction] = useState('close');
  const [form] = Form.useForm()

  const updateWalletData = () => {
    walletFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        let filters = []

        data.data.forEach((wallet) => {
          let childTitle = ''
          if (wallet.wallet_type === 1 || wallet.wallet_type === 2) {
            childTitle = '[钱包] '
          } else if (wallet.wallet_type === 3) {
            childTitle = '[债务] '
          } else if (wallet.wallet_type === 4) {
            childTitle = '[愿望单] '
          }
          let child = {
            title: childTitle + wallet.name,
            value: -wallet.id,
            selectable: false,
            children: [],
          }
          let filterChild = {
            text: childTitle + wallet.name,
            value: -wallet.id,
            children: [],
          }
          wallet.partition.forEach((part) => {
            let limit = ''
            if (part.limit !== '0') {
              limit = '/' + part.limit
            }
            if (wallet.wallet_type === 2) {
              child.children.push({
                value: part.id,
                title: '[' + part.name + '] ' + part.currency.code
              })
              filterChild.children.push({
                value: part.id,
                text: part.name
              })
            } else {
              child.children.push({
                value: part.id,
                title: '[' + part.name + '] ' + part.balance + limit + ' ' + part.currency.code
              })
              filterChild.children.push({
                value: part.id,
                text: part.name
              })
            }
          })
          options.push(child)
          filters.push(filterChild)
        })

        setWalletData(options)
        setTableFilterPartitionData(filters)
      }
    })
  }

  useEffect(() => {
    updateWalletData()
  }, [])
  useEffect(() => {
    currencyFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        let filters = []
        data.data.forEach((item) => {
          options.push({
            value: item.id,
            label: item.name
          })
          filters.push({
            value: item.id,
            text: item.name
          })
        })
        setCurrencyData(options)
        setTableFilterCurrencyData(filters)
      }
    })
  }, [])
  useEffect(() => {
    userTransactionCategoryFind().then(({data}) => {
      if (data.code === 0) {
        setSelectedCategoryData([])
        setCategoryData(data.data)

        let filters = []
        data.data.forEach((item) => {
          filters.push({
            value: item.id,
            text: item.name
          })
        })
        setTableFilterCategoryData(filters)
      }
    })
  }, [])
  useEffect(() => {
    typeTransactionTypeFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        let filters = []
        data.data.forEach((item) => {
          options.push({
            value: item.type,
            label: item.name,
            colour: item?.colour
          })
          filters.push({
            value: item.type,
            text: item.name
          })
        })
        setTransactionTypeData(options)
        setTableFilterTransactionTypeData(filters)
      }
    })
  }, [])

  const selectTransactionType = (type) => {
    let options = []
    categoryData.forEach((item) => {
      if (item.type === type) {
        options.push({
          value: item.id,
          label: item.name,
        })
      }
    })
    setSelectedTransactionType(type)
    setSelectedCategoryData(options)
  }

  const updateTable = () => {
    setTableLoadingData(true)
    transactionFind({
      search: searchKeyword.search,
      from_partition_id: searchKeyword.from_partition_id,
      to_partition_id: searchKeyword.to_partition_id,
      from_currency_id: searchKeyword.from_currency_id,
      to_currency_id: searchKeyword.to_currency_id,
      category_id: searchKeyword.category_id,
      transaction_type: searchKeyword.transaction_type,
      index: tablePaginationData.current,
      limit: tablePaginationData.pageSize,
    }).then(({data}) => {
      if (data.code === 0) {
        setTableData(data.data)
        setTableLoadingData(false)
        setTablePaginationRespData({
          total: data.page.total,
        })
      }
    })
  }

  useEffect(updateTable, [tablePaginationData, searchKeyword]);

  const handleCheckTransaction = (id) => {
    transactionChecked({id: id}).then(({data}) => {
      if (data.code !== 0) {
        message.open({
          type: 'warning',
          content: '确认失败: ' + data.msg
        })
      } else {
        message.open({
          type: 'success',
          content: '确认成功'
        })
      }
      updateTable()
      updateWalletData()
    })
  }
  const handleDeleteTransaction = (id) => {
    transactionDelete({id: id}).then(({data}) => {
      if (data.code !== 0) {
        message.open({
          type: 'warning',
          content: '删除失败: ' + data.msg
        })
      } else {
        message.open({
          type: 'success',
          content: '删除成功'
        })
      }
      updateTable()
      updateWalletData()
    })
  }
  const handleSearchTransaction = ({keyword}) => {
    const cloneData = JSON.parse(JSON.stringify(searchKeyword))
    cloneData.search = keyword
    setSearchKeyword(cloneData)
  }
  const handleInputTransactionDataShow = (action, data) => {
    if (data) {
      const cloneData = JSON.parse(JSON.stringify(data))
      cloneData.created_date = dayjs(cloneData.created * 1000)
      cloneData.created_time = dayjs(cloneData.created * 1000)
      selectTransactionType(cloneData.type)
      form.setFieldsValue(cloneData)
    } else {
      data = {
        created_date: dayjs(),
        created_time: dayjs(),
      }
      form.setFieldsValue(data)
    }
    setInputTransactionDataAction(action)
  }

  const handleInputTransactionDataOk = () => {
    form.validateFields().then((input) => {
      input.created = dayjs(input.created_date.format('YYYY-MM-DD') + ' ' + input.created_time.format('HH:mm:ss'))
      if (inputTransactionDataAction === 'create') {
        delete input.created_date
        delete input.created_time
        input.created = dayjs(input.created).unix()
        transactionCreate(input).then(({data}) => {
          if (data.code !== 0) {
            message.open({
              type: 'warning',
              content: '创建失败'
            })
            return
          }
          message.open({
            type: 'success',
            content: '创建成功'
          })
          updateTable()
          updateWalletData()
          setInputTransactionDataAction('close')
          form.resetFields()
        })
      } else if (inputTransactionDataAction === 'update') {
        input.created = dayjs(input.created).unix()
        transactionUpdate(input).then(({data}) => {
          if (data.code !== 0) {
            message.open({
              type: 'warning',
              content: '更新失败'
            })
            return
          }
          message.open({
            type: 'success',
            content: '更新成功'
          })
          updateTable()
          updateWalletData()
          setInputTransactionDataAction('close')
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

  const columns = [
    {
      title: '类型',
      key: 'transaction_type',
      dataIndex: 'category',
      render: (category) => {
        const result = transactionTypeData.find(item => item.value === category.type)
        return (
          <div key={'div' + result?.label}>
            <Tag bordered={false} color={result?.colour}>{result?.label}</Tag>
          </div>
        )
      },
      filters: tableFilterTransactionTypeData,
      onFilter: (value, record) => true
    }, {
      title: '分类',
      key: 'transaction_category',
      dataIndex: 'category',
      render: (category) => {
        return (
          <div key={'div' + category.id}>
            <Tag bordered={false} color={category?.colour}>{category?.name}</Tag>
          </div>
        )
      },
      filters: tableFilterCategoryData,
      onFilter: (value, record) => true
    }, {
      title: '源',
      dataIndex: 'from_partition',
      render: (from_partition) => {
        if (from_partition) {
          return from_partition.wallet.name + '(' + from_partition.name + ')'
        }
        return '-'
      },
      filters: tableFilterPartitionData,
      onFilter: (value, record) => true
    }, {
      title: '目标',
      dataIndex: 'to_partition',
      render: (to_partition) => {
        if (to_partition) {
          return to_partition.wallet.name + '(' + to_partition.name + ')'
        }
        return '-'
      },
      filters: tableFilterPartitionData,
      onFilter: (value, record) => true
    }, {
      title: '源金额',
      key: 'from_value',
      align: 'right',
      render: (rowData) => {
        if (rowData.from_value === '0') {
          return '-'
        }
        if (rowData.from_currency === null) {
          return '-'
        }
        return rowData.from_value + ' ' + rowData.from_currency.code
      },
      filters: tableFilterCurrencyData,
      onFilter: (value, record) => true
    }, {
      title: '目标金额',
      key: 'to_value',
      align: 'right',
      render: (rowData) => {
        if (rowData.to_value === '0') {
          return '-'
        }
        if (rowData.to_currency === null) {
          return '-'
        }
        return rowData.to_value + ' ' + rowData.to_currency.code
      },
      filters: tableFilterCurrencyData,
      onFilter: (value, record) => true
    }, {
      title: '手续费',
      dataIndex: 'fee',
      align: 'right',
      render: (fee) => {
        if (fee === '0') {
          return '-'
        }
        return fee
      }
    }, {
      title: '描述',
      dataIndex: 'description'
    }, {
      title: '交易时间',
      dataIndex: 'created',
      render: (created) => {
        if (created === 0) {
          return '-'
        }
        return dayjs(created * 1000).format('YYYY-MM-DD HH:mm:ss')
      }
    }, {
      title: '确认时间',
      key: 'checked',
      render: (rowData) => {
        if (rowData.checked === 0) {
          return (
            <Button danger type="primary" onClick={() => handleCheckTransaction(rowData.id)}>确认</Button>
          )
        }
        return dayjs(rowData.checked * 1000).format('YYYY-MM-DD HH:mm:ss')
      }
    }, {
      title: '操作',
      render: (rowData) => {
        return (
          <div>
            <Button type="primary" style={{marginRight: '5px'}}
                    onClick={() => handleInputTransactionDataShow('update', rowData)}>编辑</Button>
            <Popconfirm
              title={'删除交易'}
              description={"你确定删除交易?"}
              onConfirm={() => handleDeleteTransaction(rowData.id)}
              onCancel={() => {
              }}
              okText="确认"
              cancelText="取消"
            >
              <Button danger type="primary">删除</Button>
            </Popconfirm>
          </div>
        );
      }
    }
  ]

  return (
    <div>
      {
        mode.isWide ? (
          <Flex style={{width: '100%', marginBottom: '15px'}} justify='space-between' align='center'>
            <Button type="primary" onClick={() => handleInputTransactionDataShow('create')}>创建</Button>
            <Form
              layout="inline"
              onFinish={handleSearchTransaction}
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
              <Button type="primary" onClick={() => handleInputTransactionDataShow('create')}>创建</Button>
            </div>
            <Form
              layout="inline"
              onFinish={handleSearchTransaction}
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
        dataSource={tableData}
        scroll={{
          x: 'max-content',
        }}
        pagination={{
          total: tablePaginationRespData.total,
          current: tablePaginationData.current,
          pageSize: tablePaginationData.pageSize,
          loading: tableLoadingData,
          showSizeChanger: true,
          showQuickJumper: true,
          pageSizeOptions: [5, 10, 15, 20, 50, 100],
          responsive: true,
          onChange: (pageNumber, pageSize) => {
            let data = JSON.parse(JSON.stringify(tablePaginationData))
            data.current = pageNumber
            data.pageSize = pageSize
            setTablePaginationData(data)
          },
          showTotal: (total) => {
            return 'Total ' + total + ' items'
          }
        }}
        onChange={
          (pagination, filters, sorter, extra) => {
            const cloneData = JSON.parse(JSON.stringify(searchKeyword))
            cloneData.from_partition_id = filters.from_partition
            cloneData.to_partition_id = filters.to_partition
            cloneData.from_currency_id = filters.from_value
            cloneData.to_currency_id = filters.to_value
            cloneData.transaction_type = filters.transaction_type
            cloneData.category_id = filters.transaction_category
            setSearchKeyword(cloneData)
            // console.log('params', filters);
          }
        }
        rowKey={'id'}
      />
      <Modal
        title={inputTransactionDataAction === 'create' ? '创建' : inputTransactionDataAction === 'update' ? '编辑' : 'Unknown'}
        open={inputTransactionDataAction !== 'close'}
        onOk={handleInputTransactionDataOk}
        onCancel={() => {
          setInputTransactionDataAction('close');
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
            inputTransactionDataAction === 'update' &&
            <Form.Item
              name="id"
              hidden
            >
              <Input/>
            </Form.Item>
          }
          {
            inputTransactionDataAction === 'create' &&
            <Form.Item
              label="交易类型"
              name="transaction_type"
              rules={[
                {
                  required: true,
                  message: '请输入交易类型'
                }
              ]}
            >
              <Select
                placeholder="交易类型"
                onChange={(t) => {
                  selectTransactionType(t)
                }}
                allowClear
                options={transactionTypeData}
              />
            </Form.Item>
          }
          <Form.Item
            label="交易分类"
            name="category_id"
            rules={[
              {
                required: true,
                message: '请输入交易分类'
              }
            ]}
          >
            <Select
              placeholder="交易分类"
              onChange={() => {
              }}
              allowClear
              options={selectedCategoryData}
            />
          </Form.Item>
          {
            inputTransactionDataAction === 'create' &&
            selectedTransactionType !== TransactionTypeIncome &&
            <Form.Item
              label="源划分"
              name="from_partition_id"
              rules={[
                {
                  required: true,
                  message: '请输入源划分'
                }
              ]}
            >
              <TreeSelect
                style={{
                  width: '100%',
                }}
                dropdownStyle={{
                  maxHeight: 400,
                  overflow: 'auto',
                }}
                placeholder="请选择源划分"
                allowClear
                treeExpandAction='click'
                treeData={walletData}
              />
            </Form.Item>
          }
          {
            inputTransactionDataAction === 'create' &&
            selectedTransactionType !== TransactionTypeExpense &&
            <Form.Item
              label="目标划分"
              name="to_partition_id"
              rules={[
                {
                  required: true,
                  message: '请输入目标划分'
                }
              ]}
            >
              <TreeSelect
                style={{
                  width: '100%',
                }}
                dropdownStyle={{
                  maxHeight: 400,
                  overflow: 'auto',
                }}
                placeholder="请选择目标划分"
                allowClear
                treeExpandAction='click'
                treeData={walletData}
              />
            </Form.Item>
          }
          {
            inputTransactionDataAction === 'create' &&
            (selectedTransactionType === TransactionTypeIncome || selectedTransactionType === TransactionTypeExpense || selectedTransactionType === TransactionTypeAutoTransfer) &&
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
                placeholder="货币"
                onChange={() => {
                }}
                allowClear
                options={currencyData}
              />
            </Form.Item>
          }
          <Form.Item
            label="描述"
            name="description"
          >
            <Input placeholder={'请输入描述'}/>
          </Form.Item>
          {
            inputTransactionDataAction === 'create' &&
            (selectedTransactionType !== TransactionTypeExpense) &&
            (selectedTransactionType !== TransactionTypeIncome) &&
            (selectedTransactionType !== TransactionTypeAutoExchange) &&
            <Form.Item
              label="源价值"
              name="from_value"
              rules={[
                {
                  required: true,
                  message: '请输入源价值'
                }
              ]}
            >
              <Input placeholder={'请输入源价值'}/>
            </Form.Item>
          }
          {
            inputTransactionDataAction === 'create' &&
            (selectedTransactionType !== TransactionTypeAutoTransfer) &&
            <Form.Item
              label="目标价值"
              name="to_value"
              rules={[
                {
                  required: true,
                  message: '请输入目标价值'
                }
              ]}
            >
              <Input placeholder={'请输入目标价值'}/>
            </Form.Item>
          }
          <Form.Item
            label="交易日期"
            name="created_date"
            rules={[
              {
                required: true,
                message: '请输入交易日期'
              }
            ]}
          >
            <DatePicker picker='date' inputReadOnly style={{width: '100%'}} placeholder={'请输入交易日期'}/>
          </Form.Item>
          <Form.Item
            label="交易时间"
            name="created_time"
            rules={[
              {
                required: true,
                message: '请输入交易时间'
              }
            ]}
          >
            <DatePicker picker='time' inputReadOnly style={{width: '100%'}} placeholder={'请输入交易时间'}/>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default Transaction
