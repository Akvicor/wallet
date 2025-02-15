import React, {useEffect, useState} from 'react'
import {Calendar, Cascader, DatePicker, Divider, Flex, message, Modal, Table, Tag} from 'antd'
import './home.css'
import {
  transactionChart, transactionChartPie,
  transactionFindRange,
  transactionViewMonth,
  transactionViewYear
} from "../../api/transaction";
import dayjs from "dayjs";
import {typeTransactionTypeFind} from "../../api/type";
import MarkdownIt from 'markdown-it';
import MdEditor, {PluginComponent} from 'react-markdown-editor-lite';
import {Line} from '@ant-design/charts';
import 'react-markdown-editor-lite/lib/index.css';
import {userBindHomeTipsFind, userBindHomeTipsSave} from "../../api/userBind";
import {Pie} from '@ant-design/plots';
import {walletFind} from "../../api/wallet";

const mdParser = new MarkdownIt(/* Markdown-it options */);

class EditorView extends PluginComponent {
  static pluginName = 'editor-view'
  static align = 'right'
  static defaultConfig = {view: {menu: true, md: false, html: true}}

  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this)
    this.state = {
      view: this.getConfig('view')
    }
  }

  handleClick() {
    this.state.view.md = false
    this.props.editor.setView(this.state.view)
    userBindHomeTipsSave({content: this.props.editor.getMdValue()}).then(({data}) => {
      if (data.code !== 0) {
        message.open({
          type: 'warning',
          content: '保存失败'
        })
        return
      }
      message.open({
        type: 'success',
        content: '保存成功'
      })
    })
  }

  render() {
    return (
      <span
        className="button button-type-update"
        title="Update"
        onClick={this.handleClick}
      >
        更新
      </span>
    );
  }
}

MdEditor.use(EditorView)

const Home = () => {
  const [monthData, setMonthData] = useState({})
  const [monthReqData, setMonthReqData] = useState({target: dayjs().unix(), direction: 0})
  const [yearData, setYearData] = useState({})
  const [yearReqData, setYearReqData] = useState({target: dayjs().unix(), direction: 0})
  const [showOpen, setShowOpen] = useState({show: false, title: 'Transactions', start: 0, end: 0});
  const [showData, setShowData] = useState([])
  const [transactionTypeData, setTransactionTypeData] = useState([])
  const [editorValue, setEditorValue] = useState('## Hi')
  const [chartFromPartition, setChartFromPartitionData] = useState([])
  const [chartReqData, setChartReqData] = useState({
    from_partition_id: [],
    unit: 'day',
    start: dayjs().startOf('month').unix(),
    end: dayjs().endOf('day').unix()
  })
  const [chartData, setChartData] = useState([])
  const [chartPieData, setChartPieData] = useState([])
  const [chartPieInfoData, setChartPieInfoData] = useState({total: 0})

  useEffect(() => {
    userBindHomeTipsFind().then(({data}) => {
      if (data.code !== 0) {
        message.open({
          type: 'warning',
          content: '获取Tips失败'
        })
        return
      }
      setEditorValue(data.data.content)
    })
  }, []);

  useEffect(() => {
    transactionViewMonth(monthReqData).then(({data}) => {
      if (data.code === 0) {
        setMonthData(data.data.days)
      }
    })
  }, [monthReqData]);
  useEffect(() => {
    transactionViewYear(yearReqData).then(({data}) => {
      if (data.code === 0) {
        setYearData(data.data.months)
      }
    })
  }, [yearReqData]);
  useEffect(() => {
    typeTransactionTypeFind().then(({data}) => {
      if (data.code === 0) {
        let options = []
        data.data.forEach((item) => {
          options.push({
            value: item.type,
            label: item.name,
            colour: item.colour
          })
        })
        setTransactionTypeData(options)
      }
    })
  }, [])
  useEffect(() => {
    transactionChart(chartReqData).then(({data}) => {
      if (data.code === 0) {
        setChartData(data.data)
      }
    })
    transactionChartPie(chartReqData).then(({data}) => {
      if (data.code === 0) {
        let info = JSON.parse(JSON.stringify(chartPieInfoData))
        info.total = 0
        data.data.forEach((item) => {
          info.total += item.value
        })
        setChartPieInfoData(info)
        setChartPieData(data.data)
      }
    })
    // eslint-disable-next-line
  }, [chartReqData]);

  const updateChartFromPartitionData = () => {
    walletFind().then(({data}) => {
      if (data.code === 0) {
        let options = []

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
            label: childTitle + wallet.name,
            value: -wallet.id,
            disableCheckbox: true,
            children: [],
          }
          wallet.partition.forEach((part) => {
            child.children.push({
              value: part.id,
              label: part.name
            })
          })
          options.push(child)
        })

        setChartFromPartitionData(options)
      }
    })
  }
  useEffect(() => {
    updateChartFromPartitionData()
  }, [])

  const lineConfig = {
    data: chartData,
    xField: 'date',
    yField: 'value',
    colorField: 'currency',
    // shapeField: 'smooth',
    // point: {
    //   shapeField: 'square',
    //   sizeField: 4,
    // },
    axis: {
      y: {
        labelFormatter: (v) => `${v}`.replace(/\d{1,3}(?=(\d{3})+$)/g, (s) => `${s},`),
      },
    },
    scale: {
      color: {
        range: ['#30BF78', '#1979C9', '#F4664A', '#FAAD14']
      }
    },
    style: {
      minWidth: 800,
      lineWidth: 2,
      lineDash: (data) => {
        if (data[0].currency === 'AVG') return [4, 4];
      },
      opacity: (data) => {
        if (data[0].currency !== 'AVG') return 0.5;
      },
    },
    slider: {
      x: {labelFormatter: (d) => d},
      y: {labelFormatter: '~s'},
    },
  };
  // color: (data) => {
  //   console.log(data)
  //   return data[0].currency === 'AVG' ? '#93D072' : '#2D71E7';
  // },
  const pieConfig = {
    data: chartPieData,
    angleField: 'value',
    colorField: 'type',
    color: (d) => {
      const item = chartPieData.find((item) => item.type === d);
      if (item) return item.colour;
      return 'red';
    },
    innerRadius: 0.6,
    label: {
      text: (d) => `${d.type}\n${d.value}`,
      style: {
        fontWeight: 'bold',
      },
    },
    legend: false,
    // legend: {
    //   color: {
    //     title: false,
    //     position: 'right',
    //     rowPadding: 5,
    //   },
    // },
    interaction: {
      elementHighlight: true,
    },
    tooltip: {
      title: 'type',
    },
    state: {
      inactive: {opacity: 0.5},
    },
    annotations: [
      {
        type: 'text',
        style: {
          text: chartPieInfoData.total.toFixed(2),
          x: '50%',
          y: '50%',
          textAlign: 'center',
          fontSize: 40,
          fontStyle: 'bold',
        },
      },
    ],
  }

  const columns = [
    {
      title: '类型',
      dataIndex: 'category',
      render: (category) => {
        const result = transactionTypeData.find(item => item.value === category.type)
        return (
          <div key={'div' + result.label}>
            <Tag bordered={false} color={result.colour}>{result.label}</Tag>
          </div>
        )
      }
    }, {
      title: '分类',
      dataIndex: 'category',
      render: (category) => {
        return (
          <div key={'div' + category.id}>
            <Tag bordered={false} color={category.colour}>{category.name}</Tag>
          </div>
        )
      }
    }, {
      title: '源',
      dataIndex: 'from_partition',
      render: (from_partition) => {
        if (from_partition) {
          return from_partition.wallet.name + '(' + from_partition.name + ')'
        }
        return '-'
      }
    }, {
      title: '目标',
      dataIndex: 'to_partition',
      render: (to_partition) => {
        if (to_partition) {
          return to_partition.wallet.name + '(' + to_partition.name + ')'
        }
        return '-'
      }
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
      }
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
      }
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
          return '未确认'
        }
        return dayjs(rowData.checked * 1000).format('YYYY-MM-DD HH:mm:ss')
      }
    }
  ]

  const showModal = (date, info) => {
    let title = 'Transaction'
    let start = 0
    let end = 0
    if (info.source === 'date') {
      title = date.format('YYYY-MM-DD')
      start = date.startOf('day').unix()
      end = date.endOf('day').unix()
    } else if (info.source === 'month') {
      title = date.format('YYYY-MM')
      start = date.startOf('month').unix()
      end = date.endOf('month').unix()
    } else {
      return
    }
    transactionFindRange({
      start: start,
      end: end,
    }).then(({data}) => {
      if (data.code === 0) {
        if (data.data !== null && data.data.length !== 0) {
          setShowData(data.data)
          setShowOpen({show: true, title: title, start: start, end: end});
        }
      }
    })
  };
  const handleClose = () => {
    setShowOpen({show: false, title: 'Transactions', start: 0, end: 0});
  };

  const handleEditorChange = ({html, text}) => {
    setEditorValue(text)
  }

  const getDateListData = (value) => {
    let key = value.format('YYYY-MM-DD')
    if (!(monthData.hasOwnProperty(key))) {
      return {
        income_merge: [],
        expense_merge: []
      }
    }
    return monthData[key]
  }
  const getMonthListData = (value) => {
    let key = value.format('YYYY-MM')
    if (!(yearData.hasOwnProperty(key))) {
      return {
        income_merge: [],
        expense_merge: []
      }
    }
    return yearData[key].merge
  }

  const cellChildRender = (listData) => {
    if (listData === null) {
      return (
        <ul className="events"></ul>
      )
    }
    if (listData.expense_merge.length === 0 && listData.income_merge.length === 0) {
      return (
        <ul className="events"></ul>
      )
    }
    if (listData.expense_merge.length === 0) {
      return (
        <ul className="events">
          {
            listData.income_merge.map((item) => (
              <li key={item.currency.id}>
                <Tag color='green' key={item.currency.id}>
                  {'+' + item.value + item.currency.symbol}
                </Tag>
              </li>
            ))
          }
        </ul>
      )
    }
    if (listData.income_merge.length === 0) {
      return (
        <ul className="events">
          {
            listData.expense_merge.map((item) => (
              <li key={item.currency.id}>
                <Tag color='red' key={item.currency.id}>
                  {'-' + item.value + item.currency.symbol}
                </Tag>
              </li>
            ))
          }
        </ul>
      )
    }
    return (
      <ul className="events">
        {
          listData.income_merge.map((item) => (
            <li key={item.currency.id}>
              <Tag color='green' key={item.currency.id}>
                {'+' + item.value + item.currency.symbol}
              </Tag>
            </li>
          ))
        }
        {
          listData.expense_merge.map((item) => (
            <li key={item.currency.id}>
              <Tag color='red' key={item.currency.id}>
                {'-' + item.value + item.currency.symbol}
              </Tag>
            </li>
          ))
        }
      </ul>
    )
  };
  const cellRender = (current, info) => {
    if (info.type === 'date') {
      const listData = getDateListData(current);
      return cellChildRender(listData);
    } else if (info.type === 'month') {
      const listData = getMonthListData(current);
      return cellChildRender(listData);
    }
    return info.originNode;
  };
  const onPanelChange = (current, mode) => {
    if (mode === 'month') {
      let data = JSON.parse(JSON.stringify(monthReqData))
      data.target = current.unix()
      setMonthReqData(data)
    } else if (mode === 'year') {
      let data = JSON.parse(JSON.stringify(yearReqData))
      data.target = current.unix()
      setYearReqData(data)
    }
  };
  return (
    <div>
      <div style={{display: 'flex', width: '100%', overflowX: 'auto'}}>
        <Calendar onPanelChange={onPanelChange} cellRender={cellRender} onSelect={showModal}
                  style={{minWidth: 900}}/>
        <Modal title={showOpen.title} open={showOpen.show} onOk={handleClose} onCancel={handleClose}
               style={{minWidth: 900}}>
          <Table
            columns={columns}
            scroll={{
              x: 'max-content',
            }}
            dataSource={showData}
            pagination={{
              pageSize: 5,
            }}
            rowKey={'id'}
          />
        </Modal>
      </div>
      <Divider/>
      <div style={{minWidth: 350, width: '100%', overflowX: 'auto'}}>
        <Flex style={{marginBottom: '5px', width: '100%'}} justify='space-between' align='flex-start'>
          <Cascader
            defaultValue={['day']}
            options={[
              {value: 'day', label: 'day'},
              {value: 'month', label: 'month'},
              {value: 'year', label: 'year'}
            ]}
            onChange={(value) => {
              let val = value[0]
              let data = JSON.parse(JSON.stringify(chartReqData))
              data.unit = val
              data.start = dayjs().startOf(val).unix()
              data.end = dayjs().endOf(val).unix()
              setChartReqData(data)
            }}
            style={{minWidth: 80}}
          />
          <Cascader
            style={{
              marginLeft: '5px',
              minWidth: 160,
              width: '100%',
            }}
            options={chartFromPartition}
            onChange={(value) => {
              let data = JSON.parse(JSON.stringify(chartReqData))
              data.from_partition_id = []
              value.forEach((item) => {
                if (item.length > 1) {
                  data.from_partition_id.push(item[item.length - 1])
                }
              })
              setChartReqData(data)
            }}
            multiple
            maxTagCount="responsive"
          />
        </Flex>
        <Flex style={{width: '100%'}} justify='space-between' align='flex-start'>
          <DatePicker
            style={{width: '100%'}}
            picker={chartReqData.unit === 'day' ? 'date' : chartReqData.unit}
            defaultValue={dayjs(chartReqData.start * 1000)}
            inputReadOnly
            onChange={(start) => {
              let data = JSON.parse(JSON.stringify(chartReqData))
              data.start = start.unix()
              setChartReqData(data)
            }}
            placeholder={'请输入交易日期'}/>
          <DatePicker
            style={{marginLeft: '5px', width: '100%'}}
            picker={chartReqData.unit === 'day' ? 'date' : chartReqData.unit}
            defaultValue={dayjs(chartReqData.end * 1000)}
            inputReadOnly
            onChange={(start) => {
              let data = JSON.parse(JSON.stringify(chartReqData))
              data.end = start.unix()
              setChartReqData(data)
            }}
            placeholder={'请输入交易日期'}/>
        </Flex>
        <Divider/>
        <Pie {...pieConfig} />
        <Divider/>
        <Line {...lineConfig} />
      </div>
      <Divider/>
      <MdEditor value={editorValue}
                view={EditorView.defaultConfig.view}
                canView={{menu: true, md: true, html: true, both: true, fullScreen: false, hideMenu: true}}
                renderHTML={text => mdParser.render(text)} onChange={handleEditorChange}/>
    </div>
  )
}

export default Home
