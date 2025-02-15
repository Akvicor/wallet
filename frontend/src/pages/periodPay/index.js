import React, {useEffect, useState} from 'react'
import {
  Button,
  Form,
  Input,
  Table,
  Popconfirm,
  Modal,
  message,
  Select,
  Tag,
  DatePicker,
  InputNumber,
  Flex,
  Collapse
} from 'antd'
import './periodPay.css'
import {currencyFind} from "../../api/currency";
import {
  periodPayCreate,
  periodPayDelete,
  periodPayFind,
  periodPaySummary,
  periodPayUpdate,
  periodPayUpdateNext
} from "../../api/periodPay";
import dayjs from "dayjs";
import {typePeriodTypeFind} from "../../api/type";
import {ColorButtonProvider} from "../../theme/button";
import {PeriodTypeDayInterval, PeriodTypeMonthInterval, PeriodTypeYearInterval} from "./type";
import {useSelector} from "react-redux";

const { Panel } = Collapse;

const PeriodPay = () => {
  const mode = useSelector(state => state.mode)
  const [searchKeyword, setSearchKeyword] = useState({
    search: ''
  })
  const [selectedPeriodType, setSelectedPeriodType] = useState(0)
  const [currencyData, setCurrencyData] = useState([])
  const [periodTypeData, setPeriodTypeData] = useState([])
  const [tableSummaryTotalData, setTableSummaryTotalData] = useState([])
  const [tableSummarySearchData, setTableSummarySearchData] = useState([])
  const [tableData, setTableData] = useState([])
  const [inputPeriodPayDataAction, setInputPeriodPayDataAction] = useState('close');
  const [form] = Form.useForm()

  useEffect(() => {
    typePeriodTypeFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        data.data.forEach((item) => {
          options.push({
            value: item.type,
            label: item.name
          })
        })
        setPeriodTypeData(options)
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

  const updateTable = (search) => {
    periodPayFind(search).then(({data}) => {
      if (data.code === 0) {
        setTableData(data.data)
      }
    })
    periodPaySummary(search).then(({data}) => {
      if (data.code === 0) {
        setTableSummaryTotalData(data.data.total)
        setTableSummarySearchData(data.data.search)
      }
    })
  }

  const handleDeletePeriodPay = (id) => {
    periodPayDelete({id: id}).then(({data}) => {
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
  const handleUpdateNext = (id) => {
    periodPayUpdateNext({id: id}).then(({data}) => {
      if (data.code !== 0) {
        message.open({
          type: 'warning',
          content: '更新失败: ' + data.msg
        })
      } else {
        message.open({
          type: 'success',
          content: '更新成功'
        })
      }
      updateTable(searchKeyword)
    })
  }
  const handleSearchPeriodPay = ({keyword}) => {
    setSearchKeyword({
      search: keyword
    })
  }
  const handleInputPeriodPayDataShow = (action, data) => {
    if (data) {
      const cloneData = JSON.parse(JSON.stringify(data))
      setSelectedPeriodType(cloneData.period_type)
      cloneData.start_at = dayjs(cloneData.start_at * 1000)
      cloneData.next_of_period = dayjs(cloneData.next_of_period * 1000)
      cloneData.expiration_date = dayjs(cloneData.expiration_date * 1000)
      form.setFieldsValue(cloneData)
    } else {
      const cloneData = {}
      cloneData.start_at = dayjs()
      cloneData.next_of_period = dayjs()
      cloneData.expiration_date = dayjs('2100-01-01')
      cloneData.expiration_times = -1
      form.setFieldsValue(cloneData)
    }
    setInputPeriodPayDataAction(action)
  }

  const handleInputPeriodPayDataOk = () => {
    form.validateFields().then((input) => {
      input.start_at = dayjs(input.start_at).unix()
      input.next_of_period = dayjs(input.next_of_period).unix()
      input.expiration_date = dayjs(input.expiration_date).unix()
      if (inputPeriodPayDataAction === 'create') {
        periodPayCreate(input).then(({data}) => {
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
          setInputPeriodPayDataAction('close')
          form.resetFields()
        })
      } else if (inputPeriodPayDataAction === 'update') {
        periodPayUpdate(input).then(({data}) => {
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
          setInputPeriodPayDataAction('close')
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

  const columnsSummary = [
    {
      title: '天',
      key: 'day',
      align: 'right',
      render: (rowData) => {
        if (rowData.day === '0') {
          return '-'
        }
        if (rowData.currency === null) {
          return '-'
        }
        return rowData.day + ' ' + rowData.currency.code
      }
    }, {
      title: '月',
      key: 'month',
      align: 'right',
      render: (rowData) => {
        if (rowData.month === '0') {
          return '-'
        }
        if (rowData.currency === null) {
          return '-'
        }
        return rowData.month + ' ' + rowData.currency.code
      }
    }, {
      title: '年',
      key: 'year',
      align: 'right',
      render: (rowData) => {
        if (rowData.year === '0') {
          return '-'
        }
        if (rowData.currency === null) {
          return '-'
        }
        return rowData.year + ' ' + rowData.currency.code
      }
    }, {
      title: '天汇购',
      key: 'day_merge',
      align: 'right',
      render: (rowData) => {
        if (rowData.day_merge === '0') {
          return '-'
        }
        if (rowData.currency === null) {
          return '-'
        }
        return rowData.day_merge + ' ' + rowData.currency.code
      }
    }, {
      title: '月汇购',
      key: 'month_merge',
      align: 'right',
      render: (rowData) => {
        if (rowData.month_merge === '0') {
          return '-'
        }
        if (rowData.currency === null) {
          return '-'
        }
        return rowData.month_merge + ' ' + rowData.currency.code
      }
    }, {
      title: '年汇购',
      key: 'year_merge',
      align: 'right',
      render: (rowData) => {
        if (rowData.year_merge === '0') {
          return '-'
        }
        if (rowData.currency === null) {
          return '-'
        }
        return rowData.year_merge + ' ' + rowData.currency.code
      }
    }
  ]

  const columns = [
    {
      title: '名称',
      dataIndex: 'name'
    }, {
      title: '金额',
      key: 'value',
      align: 'right',
      render: (rowData) => {
        if (rowData.value === '0') {
          return '-'
        }
        if (rowData.currency === null) {
          return '-'
        }
        return rowData.value + ' ' + rowData.currency.code
      }
    }, {
      title: '类型',
      dataIndex: 'period_type',
      render: (period_type) => {
        const result = periodTypeData.find(period => period.value === period_type)
        if (result) return result.label
        return "-"
      }
    }, {
      title: '倒计时',
      dataIndex: 'next_of_period',
      render: (next_of_period) => {
        let countdown = 0
        let now = dayjs().startOf('day')
        if (next_of_period > now.unix()) {
          countdown = dayjs(next_of_period * 1000).diff(now, 'day')
        }
        return countdown + '天'
      }
    }, {
      title: '下一次',
      dataIndex: 'next_of_period',
      render: (next_of_period) => {
        let color
        let now = dayjs().startOf('day').unix()
        if (next_of_period === now) {
          color = 'red';
        } else if (next_of_period > now) {
          color = 'green';
        } else {
          color = 'black';
        }
        return (
          <Tag color={color} key={'next_period' + next_of_period}>
            {dayjs(next_of_period * 1000).format('YYYY-MM-DD')}
          </Tag>
        )
      }
    }, {
      title: '剩余提醒次数',
      dataIndex: 'expiration_times',
      render: (expiration_times) => {
        if (expiration_times < 0) {
          return '-'
        }
        return expiration_times
      }
    }, {
      title: '开始时间',
      dataIndex: 'start_at',
      render: (start_at) => {
        return dayjs(start_at * 1000).format('YYYY-MM-DD')
      }
    }, {
      title: '过期时间',
      dataIndex: 'expiration_date',
      render: (expiration_date) => {
        return dayjs(expiration_date * 1000).format('YYYY-MM-DD')
      }
    }, {
      title: '描述',
      dataIndex: 'description'
    }, {
      title: '操作',
      render: (rowData) => {
        return (
          <div>
            <ColorButtonProvider color="green">
              <Popconfirm
                title={'更新订阅'}
                description={"你确定更新订阅?"}
                onConfirm={() => handleUpdateNext(rowData.id)}
                onCancel={() => {
                }}
                okText="确认"
                cancelText="取消"
              >
                <Button type="primary" style={{marginRight: '5px'}}>更新</Button>
              </Popconfirm>
            </ColorButtonProvider>
            <Button type="primary" style={{marginRight: '5px'}}
                    onClick={() => handleInputPeriodPayDataShow('update', rowData)}>编辑</Button>
            <Popconfirm
              title={'删除订阅'}
              description={"你确定删除 " + rowData.name + "?"}
              onConfirm={() => handleDeletePeriodPay(rowData.id)}
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
            <Button type="primary" onClick={() => handleInputPeriodPayDataShow('create')}>创建</Button>
            <Form
              layout="inline"
              onFinish={handleSearchPeriodPay}
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
              <Button type="primary" onClick={() => handleInputPeriodPayDataShow('create')}>创建</Button>
            </div>
            <Form
              layout="inline"
              onFinish={handleSearchPeriodPay}
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
      <Collapse style={{marginBottom: '15px'}}>
        <Panel header="总计" key="1">
          <Table pagination={false} columns={columnsSummary} dataSource={tableSummaryTotalData} />
        </Panel>
        <Panel header="过滤" key="2">
          <Table pagination={false} columns={columnsSummary} dataSource={tableSummarySearchData} />
        </Panel>
      </Collapse>
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
        title={inputPeriodPayDataAction === 'create' ? '创建' : inputPeriodPayDataAction === 'update' ? '编辑' : 'Unknown'}
        open={inputPeriodPayDataAction !== 'close'}
        onOk={handleInputPeriodPayDataOk}
        onCancel={() => {
          setInputPeriodPayDataAction('close');
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
            inputPeriodPayDataAction === 'update' &&
            <Form.Item
              name="id"
              hidden
            >
              <Input/>
            </Form.Item>
          }
          <Form.Item
            label="类型"
            name="period_type"
            rules={[
              {
                required: true,
                message: '请输入类型'
              }
            ]}
          >
            <Select
              placeholder="类型"
              onChange={(t) => {
                setSelectedPeriodType(t)
              }}
              allowClear
              options={periodTypeData}
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
          <Form.Item
            label="金额"
            name="value"
            rules={[
              {
                required: true,
                message: '请输入金额'
              }
            ]}
          >
            <Input placeholder={'请输入金额'}/>
          </Form.Item>
          <Form.Item
            label="开始日期"
            name="start_at"
            rules={[
              {
                required: true,
                message: '请输入开始日期'
              }
            ]}
          >
            <DatePicker format="YYYY-MM-DD" style={{width: '100%'}}/>
          </Form.Item>
          <Form.Item
            label="下一次日期"
            name="next_of_period"
            rules={[
              {
                required: true,
                message: '请输入下一次日期'
              }
            ]}
          >
            <DatePicker format="YYYY-MM-DD" style={{width: '100%'}}/>
          </Form.Item>
          <Form.Item
            label="截止日期"
            name="expiration_date"
            rules={[
              {
                required: true,
                message: '请输入截止日期'
              }
            ]}
          >
            <DatePicker format="YYYY-MM-DD" style={{width: '100%'}}/>
          </Form.Item>
          <Form.Item
            label="截止次数"
            name="expiration_times"
            rules={[
              {
                required: true,
                message: '请输入截止次数'
              }
            ]}
          >
            <InputNumber placeholder={'-1代表无限制'} style={{width: '100%'}}/>
          </Form.Item>
          {
            (selectedPeriodType === PeriodTypeDayInterval ||
              selectedPeriodType === PeriodTypeMonthInterval ||
              selectedPeriodType === PeriodTypeYearInterval) &&
            <Form.Item
              label="间隔"
              name="interval_of_period"
              rules={[
                {
                  required: true,
                  message: '请输入间隔'
                }
              ]}
            >
              <InputNumber placeholder={'请输入间隔'} style={{width: '100%'}}/>
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
    </div>
  )
}

export default PeriodPay
