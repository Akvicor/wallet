import React, {useEffect, useState} from 'react'
import {Button, Form, Input, Table, Popconfirm, Modal, message, Select, Tag, Flex} from 'antd'
import './userTransactionCategory.css'
import {typeTransactionTypeFind} from "../../api/type";
import {
  userTransactionCategoryCreate,
  userTransactionCategoryDelete,
  userTransactionCategoryFind, userTransactionCategoryUpdate, userTransactionCategoryUpdateSequence
} from "../../api/userTransactionCategory";
import {useSelector} from "react-redux";
import {RandomColour16} from "./colour";

const UserTransactionCategory = () => {
  const mode = useSelector(state => state.mode)
  const [searchKeyword, setSearchKeyword] = useState({
    search: ''
  })
  const [updateOrder, setUpdateOrder] = useState(false)
  const [transactionTypeData, setTransactionTypeData] = useState([])
  const [tableData, setTableData] = useState([])
  const [inputUserTransactionCategoryDataAction, setInputUserTransactionCategoryDataAction] = useState('close');
  const [form] = Form.useForm()

  useEffect(() => {
    typeTransactionTypeFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        data.data.forEach((item) => {
          options.push({
            value: item.type,
            label: item.name
          })
        })
        setTransactionTypeData(options)
        updateTable(searchKeyword)
      }
    })
  }, [searchKeyword])

  const updateTable = (search) => {
    userTransactionCategoryFind(search).then(({data}) => {
      if (data.code === 0) {
        setTableData(data.data)
      }
    })
  }

  const updateSequence = (id, target) => {
    userTransactionCategoryUpdateSequence({id: id, target: target}).then(({data}) => {
      if (data.code === 0) {
        updateTable(searchKeyword)
      }
    })
  }

  const handleDeleteUserTransactionCategory = (id) => {
    userTransactionCategoryDelete({id: id}).then(({data}) => {
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
  const handleSearchUserTransactionCategory = ({keyword}) => {
    setSearchKeyword({
      search: keyword
    })
  }
  const handleInputUserTransactionCategoryDataShow = (action, data) => {
    if (data) {
      const cloneData = JSON.parse(JSON.stringify(data))
      form.setFieldsValue(cloneData)
    } else {
      const cloneData = {}
      cloneData.colour = RandomColour16()
      form.setFieldsValue(cloneData)
    }
    setInputUserTransactionCategoryDataAction(action)
  }

  const handleInputUserTransactionCategoryOk = () => {
    form.validateFields().then((input) => {
      if (inputUserTransactionCategoryDataAction === 'create') {
        userTransactionCategoryCreate(input).then(({data}) => {
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
          setInputUserTransactionCategoryDataAction('close')
          form.resetFields()
        })
      } else if (inputUserTransactionCategoryDataAction === 'update') {
        userTransactionCategoryUpdate(input).then(({data}) => {
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
          setInputUserTransactionCategoryDataAction('close')
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
      dataIndex: 'type',
      render: (type) => {
        const result = transactionTypeData.find(item => item.value === type)
        if (result) {
          return result.label
        } else {
          return '未知'
        }
      }
    }, {
      title: '名称',
      dataIndex: 'name'
    }, {
      title: '颜色',
      dataIndex: 'colour',
      render: (colour) => {
        return (
          <Tag color={colour}>
            {colour}
          </Tag>
        );
      }
    }, {
      title: '操作',
      render: (rowData) => {
        return (
          <div>
            <Button type="primary" style={{marginRight: '5px'}}
                    onClick={() => handleInputUserTransactionCategoryDataShow('update', rowData)}>编辑</Button>
            <Popconfirm
              title={'删除汇率'}
              description={"你确定删除 " + rowData.name + "?"}
              onConfirm={() => handleDeleteUserTransactionCategory(rowData.id)}
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
  if (updateOrder) {
    columns.push({
      title: '排序',
      render: (rowData) => {
        return (
          <div>
            <Button type="default" style={{marginRight: '5px'}}>{rowData.sequence}</Button>
            <Button type="default" style={{marginRight: '5px'}}
                    onClick={() => updateSequence(rowData.id, rowData.sequence - 1)}>上</Button>
            <Button type="default" onClick={() => updateSequence(rowData.id, rowData.sequence + 1)}>下</Button>
          </div>
        );
      }
    })
  }

  return (
    <div>
      {
        mode.isWide ? (
          <Flex style={{width: '100%', marginBottom: '15px'}} justify='space-between' align='center'>
            <div style={{marginBottom: 6}}>
              <Button type="primary" onClick={() => handleInputUserTransactionCategoryDataShow('create')}
                      style={{marginRight: 6}}>创建</Button>
              <Button type="default" onClick={() => setUpdateOrder(!updateOrder)}>排序</Button>
            </div>
            <Form
              layout="inline"
              onFinish={handleSearchUserTransactionCategory}
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
              <Button type="primary" onClick={() => handleInputUserTransactionCategoryDataShow('create')}
                      style={{marginRight: 6}}>创建</Button>
              <Button type="default" onClick={() => setUpdateOrder(!updateOrder)}>排序</Button>
            </div>
            <Form
              layout="inline"
              onFinish={handleSearchUserTransactionCategory}
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
        title={inputUserTransactionCategoryDataAction === 'create' ? '创建' : inputUserTransactionCategoryDataAction === 'update' ? '编辑' : 'Unknown'}
        open={inputUserTransactionCategoryDataAction !== 'close'}
        onOk={handleInputUserTransactionCategoryOk}
        onCancel={() => {
          setInputUserTransactionCategoryDataAction('close');
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
            inputUserTransactionCategoryDataAction === 'update' &&
            <Form.Item
              name="id"
              hidden
            >
              <Input/>
            </Form.Item>
          }
          <Form.Item
            label="交易类型"
            name="type"
            rules={[
              {
                required: true,
                message: '请输入交易类型'
              }
            ]}
          >
            <Select
              placeholder="交易类型"
              onChange={() => {
              }}
              allowClear
              options={transactionTypeData}
            />
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
          <Form.Item
            label="颜色"
            name="colour"
            rules={[
              {
                required: true,
                message: '请输入颜色'
              }
            ]}
          >
            <Input placeholder={'请输入颜色'}/>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default UserTransactionCategory
