import React, {useEffect, useState} from 'react'
import {Button, Form, Input, Table, Popconfirm, Modal, message, Select, Flex} from 'antd'
import {
  userExchangeRateCreate,
  userExchangeRateFind,
  userExchangeRateUpdate,
  userExchangeRateDelete
} from "../../api/userExchangeRate";
import './userExchangeRate.css'
import {currencyFind} from "../../api/currency";
import {useSelector} from "react-redux";

const UserExchangeRate = () => {
  const mode = useSelector(state => state.mode)
  const [searchKeyword, setSearchKeyword] = useState({
    search: ''
  })
  const [currencyData, setCurrencyData] = useState([])
  const [tableData, setTableData] = useState([])
  const [inputUserExchangeRateDataAction, setInputUserExchangeRateDataAction] = useState('close');
  const [form] = Form.useForm()

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

  const updateTable = (search) => {
    userExchangeRateFind(search).then(({data}) => {
      if (data.code === 0) {
        setTableData(data.data)
      }
    })
  }

  const handleDeleteUserExchangeRate = (id) => {
    userExchangeRateDelete({id: id}).then(({data}) => {
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
      updateTable(searchKeyword)
    })
  }
  const handleSearchUserExchangeRate = ({keyword}) => {
    setSearchKeyword({
      search: keyword
    })
  }
  const handleInputUserExchangeRateDataShow = (action, data) => {
    if (data) {
      const cloneData = JSON.parse(JSON.stringify(data))
      form.setFieldsValue(cloneData)
    }
    setInputUserExchangeRateDataAction(action)
  }

  const handleInputUserExchangeRateDataOk = () => {
    form.validateFields().then((input) => {
      if (inputUserExchangeRateDataAction === 'create') {
        userExchangeRateCreate(input).then(({data}) => {
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
          updateTable(searchKeyword)
          setInputUserExchangeRateDataAction('close')
          form.resetFields()
        })
      } else if (inputUserExchangeRateDataAction === 'update') {
        userExchangeRateUpdate(input).then(({data}) => {
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
          updateTable(searchKeyword)
          setInputUserExchangeRateDataAction('close')
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
      title: '源货币',
      dataIndex: 'from_currency',
      render: (from_currency) => {
        return from_currency.name
      }
    }, {
      title: '目标货币',
      dataIndex: 'to_currency',
      render: (to_currency) => {
        return to_currency.name
      }
    }, {
      title: '汇率',
      dataIndex: 'rate'
    }, {
      title: '操作',
      render: (rowData) => {
        return (
          <div>
            <Button type="primary" style={{marginRight: '5px'}}
                    onClick={() => handleInputUserExchangeRateDataShow('update', rowData)}>编辑</Button>
            <Popconfirm
              title={'删除汇率'}
              description={"你确定删除 " + rowData.from_currency.name + "->" + rowData.to_currency.name + "?"}
              onConfirm={() => handleDeleteUserExchangeRate(rowData.id)}
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

  useEffect(() => {
    updateTable(searchKeyword)
  }, [searchKeyword]);

  return (
    <div>
      {
        mode.isWide ? (
          <Flex style={{width: '100%', marginBottom: '15px'}} justify='space-between' align='center'>
            <Button type="primary" onClick={() => handleInputUserExchangeRateDataShow('create')}>创建</Button>
            <Form
              layout="inline"
              onFinish={handleSearchUserExchangeRate}
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
              <Button type="primary" onClick={() => handleInputUserExchangeRateDataShow('create')}>创建</Button>
            </div>
            <Form
              layout="inline"
              onFinish={handleSearchUserExchangeRate}
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
          dataSource={tableData}
          pagination={{
            pageSizeOptions: [10, 15, 20, 50, 100],
            responsive: true,
            showQuickJumper: true,
            showSizeChanger: true
          }}
          rowKey={'id'}
      />
      <Modal
        title={inputUserExchangeRateDataAction === 'create' ? '创建' : inputUserExchangeRateDataAction === 'update' ? '编辑' : 'Unknown'}
        open={inputUserExchangeRateDataAction !== 'close'}
        onOk={handleInputUserExchangeRateDataOk}
        onCancel={() => {
          setInputUserExchangeRateDataAction('close');
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
            inputUserExchangeRateDataAction === 'update' &&
            <Form.Item
              name="id"
              hidden
            >
              <Input/>
            </Form.Item>
          }
          <Form.Item
            label="源货币"
            name="from_currency_id"
            rules={[
              {
                required: true,
                message: '请输入源货币'
              }
            ]}
          >
            <Select
              placeholder="源货币"
              onChange={() => {
              }}
              allowClear
              options={currencyData}
            />
          </Form.Item>
          <Form.Item
            label="目标货币"
            name="to_currency_id"
            rules={[
              {
                required: true,
                message: '请输入目标货币'
              }
            ]}
          >
            <Select
              placeholder="目标货币"
              onChange={() => {
              }}
              allowClear
              options={currencyData}
            />
          </Form.Item>
          <Form.Item
            label="汇率"
            name="rate"
            rules={[
              {
                required: true,
                message: '请输入汇率'
              }
            ]}
          >
            <Input placeholder={'请输入汇率'}/>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default UserExchangeRate
