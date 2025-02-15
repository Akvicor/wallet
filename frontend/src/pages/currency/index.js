import React, {useEffect, useState} from 'react'
import {Button, Form, Input, Table, Popconfirm, Modal, message, Flex} from 'antd'
import {currencyCreate, currencyFind, currencyUpdate, currencyDelete} from "../../api/currency";
import './currency.css'
import {useSelector} from "react-redux";

const Currency = () => {
  const mode = useSelector(state => state.mode)
  const [searchKeyword, setSearchKeyword] = useState({
    search: ''
  })
  const [tableData, setTableData] = useState([])
  const [inputCurrencyDataAction, setInputCurrencyDataAction] = useState('close');
  const [form] = Form.useForm()

  const updateTable = (search) => {
    currencyFind(search).then(({data}) => {
      if (data.code === 0) {
        setTableData(data.data)
      }
    })
  }

  const handleDeleteCurrency = (id) => {
    currencyDelete({id: id}).then(({data}) => {
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
  const handleSearchCurrency = ({keyword}) => {
    setSearchKeyword({
      search: keyword
    })
  }
  const handleInputCurrencyDataShow = (action, data) => {
    if (data) {
      const cloneData = JSON.parse(JSON.stringify(data))
      form.setFieldsValue(cloneData)
    }
    setInputCurrencyDataAction(action)
  }

  const handleInputCurrencyDataOk = () => {
    form.validateFields().then((input) => {
      if (inputCurrencyDataAction === 'create') {
        currencyCreate(input).then(({data}) => {
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
          setInputCurrencyDataAction('close')
          form.resetFields()
        })
      } else if (inputCurrencyDataAction === 'update') {
        currencyUpdate(input).then(({data}) => {
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
          setInputCurrencyDataAction('close')
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
      title: '名称',
      dataIndex: 'name',
    }, {
      title: '英文名称',
      dataIndex: 'english_name'
    }, {
      title: '缩写',
      dataIndex: 'code'
    }, {
      title: '符号',
      dataIndex: 'symbol'
    }, {
      title: '操作',
      render: (rowData) => {
        return (
          <div>
            <Button type="primary" style={{marginRight: '5px'}}
                    onClick={() => handleInputCurrencyDataShow('update', rowData)}>编辑</Button>
            <Popconfirm
              title={'删除货币类型'}
              description={"你确定删除" + rowData.name + "?"}
              onConfirm={() => handleDeleteCurrency(rowData.id)}
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
            <Button type="primary" onClick={() => handleInputCurrencyDataShow('create')}>创建</Button>
            <Form
              layout="inline"
              onFinish={handleSearchCurrency}
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
            <div style={{marginBottom: 6}}>
              <Button type="primary" onClick={() => handleInputCurrencyDataShow('create')}>创建</Button>
            </div>
            <Form
              layout="inline"
              onFinish={handleSearchCurrency}
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
        title={inputCurrencyDataAction === 'create' ? '创建' : inputCurrencyDataAction === 'update' ? '编辑' : 'Unknown'}
        open={inputCurrencyDataAction !== 'close'}
        onOk={handleInputCurrencyDataOk}
        onCancel={() => {
          setInputCurrencyDataAction('close');
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
            inputCurrencyDataAction === 'update' &&
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
            label="英文名称"
            name="english_name"
            rules={[
              {
                required: true,
                message: '请输入名称'
              }
            ]}
          >
            <Input placeholder={'请输入英文名称'}/>
          </Form.Item>
          <Form.Item
            label="缩写"
            name="code"
            rules={[
              {
                required: true,
                message: '请输入缩写'
              }
            ]}
          >
            <Input placeholder={'请输入缩写'}/>
          </Form.Item>
          <Form.Item
            label="符号"
            name="symbol"
            rules={[
              {
                required: true,
                message: '请输入符号'
              }
            ]}
          >
            <Input placeholder={'请输入符号'}/>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default Currency
