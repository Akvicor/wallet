import React, {useEffect, useState} from 'react'
import {Button, Form, Input, Table, Popconfirm, Modal, message, Flex} from 'antd'
import {cardTypeCreate, cardTypeFind, cardTypeUpdate, cardTypeDelete} from "../../api/cardType";
import './cardType.css'
import {useSelector} from "react-redux";

const CardType = () => {
  const mode = useSelector(state => state.mode)
  const [searchKeyword, setSearchKeyword] = useState({
    search: ''
  })
  const [tableData, setTableData] = useState([])
  const [inputCardTypeDataAction, setInputCardTypeDataAction] = useState('close');
  const [form] = Form.useForm()

  const updateTable = (search) => {
    cardTypeFind(search).then(({data}) => {
      if (data.code === 0) {
        setTableData(data.data)
      }
    })
  }

  const handleDeleteCardType = (id) => {
    cardTypeDelete({id: id}).then(({data}) => {
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
  const handleSearchCardType = ({keyword}) => {
    setSearchKeyword({
      search: keyword
    })
  }
  const handleInputCardTypeDataShow = (action, data) => {
    if (data) {
      const cloneData = JSON.parse(JSON.stringify(data))
      form.setFieldsValue(cloneData)
    }
    setInputCardTypeDataAction(action)
  }

  const handleInputCardTypeDataOk = () => {
    form.validateFields().then((input) => {
      if (inputCardTypeDataAction === 'create') {
        cardTypeCreate(input).then(({data}) => {
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
          setInputCardTypeDataAction('close')
          form.resetFields()
        })
      } else if (inputCardTypeDataAction === 'update') {
        cardTypeUpdate(input).then(({data}) => {
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
          setInputCardTypeDataAction('close')
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
      title: '操作',
      render: (rowData) => {
        return (
          <div>
            <Button type="primary" style={{marginRight: '5px'}}
                    onClick={() => handleInputCardTypeDataShow('update', rowData)}>编辑</Button>
            <Popconfirm
              title={'删除银行卡类型'}
              description={"你确定删除" + rowData.name + "?"}
              onConfirm={() => handleDeleteCardType(rowData.id)}
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
            <Button type="primary" onClick={() => handleInputCardTypeDataShow('create')}>创建</Button>
            <Form
              layout="inline"
              onFinish={handleSearchCardType}
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
              <Button type="primary" onClick={() => handleInputCardTypeDataShow('create')}>创建</Button>
            </div>
            <Form
              layout="inline"
              onFinish={handleSearchCardType}
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
        title={inputCardTypeDataAction === 'create' ? '创建' : inputCardTypeDataAction === 'update' ? '编辑' : 'Unknown'}
        open={inputCardTypeDataAction !== 'close'}
        onOk={handleInputCardTypeDataOk}
        onCancel={() => {
          setInputCardTypeDataAction('close');
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
            inputCardTypeDataAction === 'update' &&
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
        </Form>
      </Modal>
    </div>
  )
}

export default CardType
