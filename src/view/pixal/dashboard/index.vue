<template>
  <div  style="width: 100%;">
    <div class="operation-bar">
        <el-date-picker
          v-model="search.dateRange"
          type="daterange"
          :placeholder="$t('lang.start_date')"
          style="margin-right: 10px;max-width:300px"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          clearable
          class="max-width:100px"
        ></el-date-picker>
        <el-button type="primary" icon="search" @click="getChartData">{{$t('lang.search')}}</el-button>
      </div>
    <div class="mt-5">
        <div ref="authorizationDom" class="chart" style="width: 100%; height: 500px;"></div>
        <div ref="reversalDom" class="chart" style="width: 100%; height: 500px;"></div>
        <div ref="inoutDom" class="chart" style="width: 100%; height: 500px;"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, reactive } from 'vue';
import { cardReportByDay} from '@/api/finance';
// 引入 ECharts 核心模块和必要组件
import * as echarts from 'echarts';
import { useI18n } from 'vue-i18n';
import { graphic } from 'echarts/lib/export';
const { t } = useI18n();
// import {
//   BarChart
// } from 'echarts/charts';
// import {
//   TitleComponent,
//   TooltipComponent,
//   GridComponent,
//   LegendComponent
// } from 'echarts/components';
// import {
//   CanvasRenderer
// } from 'echarts/renderers';

// 注册必要组件
// echarts.use([
//   TitleComponent,
//   TooltipComponent,
//   GridComponent,
//   LegendComponent,
//   BarChart,
//   CanvasRenderer
// ]);

// 1. 定义 DOM 元素的 ref
const authorizationDom = ref(null);
// 2. 存储 ECharts 实例
let authorizationChart = null;
// 1. 定义 DOM 元素的 ref
const reversalDom = ref(null);
// 2. 存储 ECharts 实例
let reversalChart = null;
const inoutDom = ref(null);
// 2. 存储 ECharts 实例
let inoutChart = null;
const chartData = ref([])
const totalAuthorizationAmount = ref(0)
const totalAuthorizationCount = ref(0)
// 3. 定义图表配置
const option = {
  // 标题
  title: {
    text: t('lang.authorization_statistics'),
    left: 'center',
    top:'0px',
    textStyle: {
      fontSize: 20,
      fontWeight: 'bold'
    }
  },
  graphic: {
    elements: [
      {
        type: 'text',
        left: '60',
        top: '40',
        style: {
          text: ``,
          fontSize: 12,
          align: 'left'
        }
      }
    ]
  },
  // 提示框，trigger: 'axis' 会在鼠标悬停到x轴时触发
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'shadow' // 默认为 'line'，'shadow' 表示阴影指示器
    }
  },
  // 图例，用于筛选系列
  legend: {
    data: ['Amount', 'Quantity'],
    bottom: 15, // 放在标题下面
    left: 'center'
  },
  // 网格配置，调整图表在容器中的位置
  grid: {
    left: '3%',
    right: '4%',
    top: '90px',
    // bottom: '3%',
    containLabel: true // 防止标签溢出
  },
  // X轴：时间（这里用月份作为类别）
  xAxis: {
    type: 'category',
    data: []
  },
  // Y轴：金额
  yAxis: [{
    type: 'value',
    name: `${t('lang.amount')}(usd)`,
    axisLabel: {
      formatter: '{value} usd'
    }
  },{
    type: 'value',
    name: `${t('lang.quanlity')}`,
    axisLabel: {
      formatter: '{value}'
    }
  }],
  // 系列数据
  series: [
    {
      name: 'Amount',
      type: 'line',
      yAxisIndex: 0, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: '#5470C6' // 收入的颜色
      },
      smooth: true,
        lineStyle: {
          width: 3
        }
    },
    {
      name: 'Quantity',
      type: 'line',
      yAxisIndex: 1, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: '#EE6666' // 支出的颜色
      },
      smooth: true,
        lineStyle: {
          width: 3
        }
    }
  ]
};
const inoutOption = {
  // 标题
  title: {
    text: t('lang.recharge_and_withdrawal_statistics'),
    left: 'center',
    top:'0px',
    textStyle: {
      fontSize: 20,
      fontWeight: 'bold'
    }
  },
  // 提示框，trigger: 'axis' 会在鼠标悬停到x轴时触发
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'shadow' // 默认为 'line'，'shadow' 表示阴影指示器
    }
  },
  // 图例，用于筛选系列
  legend: {
    data: ['Recharge', 'Withdrawal'],
    bottom: 15, // 放在标题下面
    left: 'center'
  },
  // 网格配置，调整图表在容器中的位置
  grid: {
    left: '3%',
    right: '4%',
    top: '90px',
    // bottom: '3%',
    containLabel: true // 防止标签溢出
  },
  // X轴：时间（这里用月份作为类别）
  xAxis: {
    type: 'category',
    data: []
  },
  // Y轴：金额
  yAxis: [{
    type: 'value',
    name: `${t('lang.recharge_amount')}(usd)`,
    axisLabel: {
      formatter: '{value} usd'
    }
  },{
    type: 'value',
    name: `${t('lang.withdrawal_amount')}(usd)`,
    axisLabel: {
      formatter: '{value} usd'
    }
  }],
  // 系列数据
  series: [
    {
      name: 'Recharge',
      type: 'line',
      yAxisIndex: 0, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: '#5470C6' // 收入的颜色
      },
      smooth: true,
        lineStyle: {
          width: 3
        }
    },
    {
      name: 'Withdrawal',
      type: 'line',
      yAxisIndex: 1, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: '#EE6666' // 支出的颜色
      },
      smooth: true,
        lineStyle: {
          width: 3
        }
    }
  ]
};
const reversalOption = {
  // 标题
  title: {
    text: t('lang.reversal_statistics'),
    left: 'center',
    top:'0px',
    textStyle: {
      fontSize: 20,
      fontWeight: 'bold'
    }
  },
  // 提示框，trigger: 'axis' 会在鼠标悬停到x轴时触发
  tooltip: {
    trigger: 'axis',
    axisPointer: {
      type: 'shadow' // 默认为 'line'，'shadow' 表示阴影指示器
    }
  },
  // 图例，用于筛选系列
  legend: {
    data: ['Amount', 'Quantity'],
    bottom: 15, // 放在标题下面
    left: 'center'
  },
  // 网格配置，调整图表在容器中的位置
  grid: {
    left: '3%',
    right: '4%',
    // bottom: '3%',
    top: '90px',
    containLabel: true // 防止标签溢出
  },
  // X轴：时间（这里用月份作为类别）
  xAxis: {
    type: 'category',
    data: []
  },
  // Y轴：金额
  yAxis: [{
    type: 'value',
    name: `${t('lang.amount')}(usd)`,
    axisLabel: {
      formatter: '{value} usd'
    }
  },{
    type: 'value',
    name: `${t('lang.quanlity')}`,
    axisLabel: {
      formatter: '{value}'
    }
  }],
  // 系列数据
  series: [
    {
      name: 'Amount',
      type: 'line',
      yAxisIndex: 0, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: '#5470C6' // 收入的颜色
      },
      smooth: true,
        lineStyle: {
          width: 3
        }
    },
    {
      name: 'Quantity',
      type: 'line',
      yAxisIndex: 1, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: '#EE6666' // 支出的颜色
      },
      smooth: true,
        lineStyle: {
          width: 3
        }
    }
  ]
};
const search = reactive({
  startTime: '',
  endTime: '',

});

const getChartData = async () => {
    if(search.dateRange && search.dateRange.length ===2){
        search.startTime = search.dateRange[0]
        search.endTime = search.dateRange[1]
    }
  cardReportByDay(search).then(res =>{
    if(res.code ===0){
      chartData.value = res.data
      chartData.value.authorizationAmountArr = chartData.value.map(item => item.authorizationAmount)
      chartData.value.authorizationCountArr = chartData.value.map(item => item.authorizationCount)
      chartData.value.reversalAmountArr = chartData.value.map(item => item.reversalAmount)
      chartData.value.reversalCountArr = chartData.value.map(item => item.reversalCount)
      chartData.value.rechargeAmountArr = chartData.value.map(item => item.cardRechargeAmount)
      chartData.value.withdrawalAmountArr = chartData.value.map(item => item.cardWithdrawAmount)
      chartData.value.totalAuthorizationAmount = chartData.value.authorizationAmountArr.reduce((a, b) => parseFloat(a) + parseFloat(b), 0).toFixed(2)
      chartData.value.totalAuthorizationCount = chartData.value.authorizationCountArr.reduce((a, b) => a + b, 0)
      chartData.value.totalReversalAmount = chartData.value.reversalAmountArr.reduce((a, b) => parseFloat(a) + parseFloat(b), 0).toFixed(2)
      chartData.value.totalReversalCount = chartData.value.reversalCountArr.reduce((a, b) => a + b, 0)
      chartData.value.totalRechargeAmount = chartData.value.rechargeAmountArr .reduce((a, b) => parseFloat(a) + parseFloat(b), 0).toFixed(2)
      chartData.value.totalWithdrawalAmount = chartData.value.withdrawalAmountArr .reduce((a, b) => parseFloat(a) + parseFloat(b), 0).toFixed(2)
      chartData.value.xAxis = chartData.value.map(item => item.reportDay)
      authorizationChart.setOption({
        xAxis: {
          data: chartData.value.xAxis
        },
        graphic: {
          elements: [
            {
              type: 'text',
              left: '60',
              top: '20',
              style: {
                text: `${t('lang.total_amount')}: $${chartData.value.totalAuthorizationAmount} \n${t('lang.total_quantity')}: ${chartData.value.totalAuthorizationCount}`,
                fontSize: 12,
                // fontWeight: 'bold',
                // fill: '#1890ff',
                align: 'left'
              }
            }
          ]
        },
        series: [
          {
            name: 'Amount',
            type: 'line',
            data: chartData.value.authorizationAmountArr
          },
          {
            name: 'Quantity',
            type: 'line',
            data: chartData.value.authorizationCountArr
          }
        ]
      });
        reversalChart.setOption({
            xAxis: {
            data: chartData.value.xAxis
            },
            graphic: {
                elements: [
                    {
                    type: 'text',
                    left: '60',
                    top: '20',
                    style: {
                        text: `${t('lang.total_amount')}: $${chartData.value.totalReversalAmount} \n${t('lang.total_quantity')}: ${chartData.value.totalReversalCount}`,
                        fontSize: 12,
                        align: 'left'
                    }
                    }
                ]
            },
            series: [
            {
                name: 'Amount',
                type: 'line',
                data: chartData.value.reversalAmountArr
            },
            {
                name: 'Quantity',
                type: 'line',
                data: chartData.value.reversalCountArr
            }
            ]
        });
        inoutChart.setOption({
            xAxis: {
            data: chartData.value.xAxis
            },
            graphic: {
                elements: [
                    {
                    type: 'text',
                    left: '60',
                    top: '20',
                    style: {
                        text: `${t('lang.total_recharge_amount')}: $${chartData.value.totalRechargeAmount} \n${t('lang.total_withdrawal_amount')}: $${chartData.value.totalWithdrawalAmount}`,
                        fontSize: 12,
                        align: 'left'
                    }
                    }
                ]
            },
            series: [
            {
                name: 'Recharge',
                type: 'line',
                data: chartData.value.rechargeAmountArr
            },
            {
                name: 'Withdrawal',
                type: 'line',
                data: chartData.value.withdrawalAmountArr
            }
            ]
        });

    }
  })
};
// 4. 定义图表自适应调整大小的函数
const resizeChart = () => {
  authorizationChart?.resize();
    reversalChart?.resize();
    inoutChart?.resize();
};
// 5. 在组件挂载后初始化图表
onMounted(() => {
    getChartData()
  if (authorizationDom.value) {
    // 初始化 ECharts 实例
    authorizationChart = echarts.init(authorizationDom.value);
    // 设置配置项
    authorizationChart.setOption(option);
    // 添加窗口大小变化的监听
    window.addEventListener('resize', resizeChart);
  }
  if (reversalDom.value) {
    // 初始化 ECharts 实例
    reversalChart = echarts.init(reversalDom.value);
    // 设置配置项
    reversalChart.setOption(reversalOption);
    // 添加窗口大小变化的监听
    window.addEventListener('resize', resizeChart);
  }
  if (inoutDom.value) {
    // 初始化 ECharts 实例
    inoutChart = echarts.init(inoutDom.value);
    // 设置配置项
    inoutChart.setOption(inoutOption);
    // 添加窗口大小变化的监听
    window.addEventListener('resize', resizeChart);
  }
});

// 6. 在组件卸载前销毁图表并移除监听
onBeforeUnmount(() => {
  // 移除监听
  window.removeEventListener('resize', resizeChart);
  // 销毁实例
  if (authorizationChart) {
    authorizationChart.dispose();
    authorizationChart = null;
  }
    if (reversalChart) {
        reversalChart.dispose();
        reversalChart = null;
    }
    if (inoutChart) {
        inoutChart.dispose();
        inoutChart = null;
    }
});

// 你也可以暴露一个方法来更新数据
// (这个例子中数据是静态的，但你可以通过 props 或 pinia/vuex 来动态更新)
// const updateChart = (newData) => {
//   myChart.setOption({
//     series: [
//       { data: newData.income },
//       { data: newData.expenditure }
//     ]
//   });
// }
// defineExpose({ updateChart });

</script>

<style scoped>
/* 样式可以保持为空，因为大小是在 div 上直接设置的 */
.chart {
  width: 100%;
  height: 500px;
  margin-bottom: 40px;
  background: #f9f9f9;
  border-radius: 20px;
  padding: 20px;
}
</style>