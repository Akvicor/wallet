import React, {useEffect, useState} from 'react'
import {Button, Form, Input, Table, Popconfirm, Modal, message, Flex} from 'antd'
import {bankCreate, bankFind, bankUpdate, bankDelete} from "../../api/bank";
import './bank.css'
import {useSelector} from "react-redux";

const Bank = () => {
  const mode = useSelector(state => state.mode)
  const [searchKeyword, setSearchKeyword] = useState({
    search: ''
  })
  const [tableData, setTableData] = useState([])
  const [inputBankDataAction, setInputBankDataAction] = useState('close');
  const [form] = Form.useForm()

  const updateTable = (search) => {
    bankFind(search).then(({data}) => {
      if (data.code === 0) {
        setTableData(data.data)
      }
    })
  }

  const handleDeleteBank = (id) => {
    bankDelete({id: id}).then(({data}) => {
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
  const handleSearchBank = ({keyword}) => {
    setSearchKeyword({
      search: keyword
    })
  }
  const handleInputBankDataShow = (action, data) => {
    if (data) {
      const cloneData = JSON.parse(JSON.stringify(data))
      form.setFieldsValue(cloneData)
    }
    setInputBankDataAction(action)
  }

  const handleInputBankDataOk = () => {
    form.validateFields().then((input) => {
      if (inputBankDataAction === 'create') {
        bankCreate(input).then(({data}) => {
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
          setInputBankDataAction('close')
          form.resetFields()
        })
      } else if (inputBankDataAction === 'update') {
        bankUpdate(input).then(({data}) => {
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
          setInputBankDataAction('close')
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
      dataIndex: 'abbr'
    }, {
      title: '电话',
      dataIndex: 'phone'
    }, {
      title: '操作',
      render: (rowData) => {
        return (
          <div>
            <Button type="primary" style={{marginRight: '5px'}}
                    onClick={() => handleInputBankDataShow('update', rowData)}>编辑</Button>
            <Popconfirm
              title={'删除银行'}
              description={"你确定删除" + rowData.name + "?"}
              onConfirm={() => handleDeleteBank(rowData.id)}
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
            <Button type="primary" onClick={() => handleInputBankDataShow('create')}>创建</Button>
            <Form
              layout="inline"
              onFinish={handleSearchBank}
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
              <Button type="primary" onClick={() => handleInputBankDataShow('create')}>创建</Button>
            </div>
            <Form
              layout="inline"
              onFinish={handleSearchBank}
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
          pagination={{
            pageSizeOptions: [10, 15, 20, 50, 100],
            responsive: true,
            showQuickJumper: true,
            showSizeChanger: true
          }}
          scroll={{
            x: 'max-content',
          }}
          rowKey={'id'}
      />
      <Modal
        title={inputBankDataAction === 'create' ? '创建' : inputBankDataAction === 'update' ? '编辑' : 'Unknown'}
        open={inputBankDataAction !== 'close'}
        onOk={handleInputBankDataOk}
        onCancel={() => {
          setInputBankDataAction('close');
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
            inputBankDataAction === 'update' &&
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
            name="abbr"
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
            label="电话"
            name="phone"
          >
            <Input placeholder={'请输入电话'}/>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default Bank
