<template>
  <div class="dashboard-page">
    <section class="dashboard-hero">
      <div>
        <p class="eyebrow">Card analytics</p>
        <h1>Dashboard</h1>
        <p class="hero-copy">Track authorization, reversal, recharge and withdrawal activity across your virtual card operations.</p>
      </div>
      <div class="operation-bar">
        <el-date-picker
          v-model="search.dateRange"
          type="daterange"
          :placeholder="$t('lang.start_date')"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          clearable
          class="date-range"
        ></el-date-picker>
        <el-button class="search-btn" type="primary" icon="search" @click="getChartData">{{$t('lang.search')}}</el-button>
      </div>
    </section>

    <section class="metric-grid">
      <article class="metric-card metric-card-primary">
        <span>{{ $t('lang.total_amount') }}</span>
        <strong>${{ chartData.totalAuthorizationAmount || '0.00' }}</strong>
        <p>{{ $t('lang.authorization_statistics') }}</p>
      </article>
      <article class="metric-card">
        <span>{{ $t('lang.total_quantity') }}</span>
        <strong>{{ chartData.totalAuthorizationCount || 0 }}</strong>
        <p>{{ $t('lang.authorization_statistics') }}</p>
      </article>
      <article class="metric-card">
        <span>{{ $t('lang.total_recharge_amount') }}</span>
        <strong>${{ chartData.totalRechargeAmount || '0.00' }}</strong>
        <p>{{ $t('lang.recharge_and_withdrawal_statistics') }}</p>
      </article>
      <article class="metric-card">
        <span>{{ $t('lang.total_withdrawal_amount') }}</span>
        <strong>${{ chartData.totalWithdrawalAmount || '0.00' }}</strong>
        <p>{{ $t('lang.recharge_and_withdrawal_statistics') }}</p>
      </article>
    </section>

    <section class="chart-grid">
      <article class="chart-card">
        <div class="chart-card-header">
          <span>01</span>
          <h2>{{ $t('lang.authorization_statistics') }}</h2>
        </div>
        <div ref="authorizationDom" class="chart"></div>
      </article>
      <article class="chart-card">
        <div class="chart-card-header">
          <span>02</span>
          <h2>{{ $t('lang.reversal_statistics') }}</h2>
        </div>
        <div ref="reversalDom" class="chart"></div>
      </article>
      <article class="chart-card chart-card-wide">
        <div class="chart-card-header">
          <span>03</span>
          <h2>{{ $t('lang.recharge_and_withdrawal_statistics') }}</h2>
        </div>
        <div ref="inoutDom" class="chart"></div>
      </article>
    </section>
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
const chartTextColor = '#dcecff'
const chartMutedColor = 'rgba(220, 236, 255, 0.58)'
const chartGridLine = 'rgba(139, 214, 255, 0.12)'
const amountColor = '#44d5ff'
const quantityColor = '#7dffcc'
const warningColor = '#ff8aa6'
const tooltipStyle = {
  backgroundColor: 'rgba(5, 16, 29, 0.92)',
  borderColor: 'rgba(139, 214, 255, 0.22)',
  textStyle: {
    color: chartTextColor
  }
}
const axisLineStyle = {
  lineStyle: {
    color: 'rgba(139, 214, 255, 0.2)'
  }
}
const splitLineStyle = {
  lineStyle: {
    color: chartGridLine
  }
}
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
      color: chartTextColor,
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
          fill: chartMutedColor,
          align: 'left'
        }
      }
    ]
  },
  // 提示框，trigger: 'axis' 会在鼠标悬停到x轴时触发
  tooltip: {
    ...tooltipStyle,
    trigger: 'axis',
    axisPointer: {
      type: 'shadow' // 默认为 'line'，'shadow' 表示阴影指示器
    }
  },
  // 图例，用于筛选系列
  legend: {
    data: ['Amount', 'Quantity'],
    bottom: 15, // 放在标题下面
    left: 'center',
    textStyle: {
      color: chartMutedColor
    }
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
    data: [],
    axisLabel: {
      color: chartMutedColor
    },
    axisLine: axisLineStyle
  },
  // Y轴：金额
  yAxis: [{
    type: 'value',
    name: `${t('lang.amount')}(usd)`,
    nameTextStyle: {
      color: chartMutedColor
    },
    axisLabel: {
      color: chartMutedColor,
      formatter: '{value} usd'
    },
    splitLine: splitLineStyle
  },{
    type: 'value',
    name: `${t('lang.quanlity')}`,
    nameTextStyle: {
      color: chartMutedColor
    },
    axisLabel: {
      color: chartMutedColor,
      formatter: '{value}'
    },
    splitLine: splitLineStyle
  }],
  // 系列数据
  series: [
    {
      name: 'Amount',
      type: 'line',
      yAxisIndex: 0, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: amountColor
      },
      smooth: true,
        lineStyle: {
          color: amountColor,
          width: 3
        }
    },
    {
      name: 'Quantity',
      type: 'line',
      yAxisIndex: 1, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: quantityColor
      },
      smooth: true,
        lineStyle: {
          color: quantityColor,
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
      color: chartTextColor,
      fontSize: 20,
      fontWeight: 'bold'
    }
  },
  // 提示框，trigger: 'axis' 会在鼠标悬停到x轴时触发
  tooltip: {
    ...tooltipStyle,
    trigger: 'axis',
    axisPointer: {
      type: 'shadow' // 默认为 'line'，'shadow' 表示阴影指示器
    }
  },
  // 图例，用于筛选系列
  legend: {
    data: ['Recharge', 'Withdrawal'],
    bottom: 15, // 放在标题下面
    left: 'center',
    textStyle: {
      color: chartMutedColor
    }
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
    data: [],
    axisLabel: {
      color: chartMutedColor
    },
    axisLine: axisLineStyle
  },
  // Y轴：金额
  yAxis: [{
    type: 'value',
    name: `${t('lang.recharge_amount')}(usd)`,
    nameTextStyle: {
      color: chartMutedColor
    },
    axisLabel: {
      color: chartMutedColor,
      formatter: '{value} usd'
    },
    splitLine: splitLineStyle
  },{
    type: 'value',
    name: `${t('lang.withdrawal_amount')}(usd)`,
    nameTextStyle: {
      color: chartMutedColor
    },
    axisLabel: {
      color: chartMutedColor,
      formatter: '{value} usd'
    },
    splitLine: splitLineStyle
  }],
  // 系列数据
  series: [
    {
      name: 'Recharge',
      type: 'line',
      yAxisIndex: 0, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: quantityColor
      },
      smooth: true,
        lineStyle: {
          color: quantityColor,
          width: 3
        }
    },
    {
      name: 'Withdrawal',
      type: 'line',
      yAxisIndex: 1, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: warningColor
      },
      smooth: true,
        lineStyle: {
          color: warningColor,
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
      color: chartTextColor,
      fontSize: 20,
      fontWeight: 'bold'
    }
  },
  // 提示框，trigger: 'axis' 会在鼠标悬停到x轴时触发
  tooltip: {
    ...tooltipStyle,
    trigger: 'axis',
    axisPointer: {
      type: 'shadow' // 默认为 'line'，'shadow' 表示阴影指示器
    }
  },
  // 图例，用于筛选系列
  legend: {
    data: ['Amount', 'Quantity'],
    bottom: 15, // 放在标题下面
    left: 'center',
    textStyle: {
      color: chartMutedColor
    }
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
    data: [],
    axisLabel: {
      color: chartMutedColor
    },
    axisLine: axisLineStyle
  },
  // Y轴：金额
  yAxis: [{
    type: 'value',
    name: `${t('lang.amount')}(usd)`,
    nameTextStyle: {
      color: chartMutedColor
    },
    axisLabel: {
      color: chartMutedColor,
      formatter: '{value} usd'
    },
    splitLine: splitLineStyle
  },{
    type: 'value',
    name: `${t('lang.quanlity')}`,
    nameTextStyle: {
      color: chartMutedColor
    },
    axisLabel: {
      color: chartMutedColor,
      formatter: '{value}'
    },
    splitLine: splitLineStyle
  }],
  // 系列数据
  series: [
    {
      name: 'Amount',
      type: 'line',
      yAxisIndex: 0, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: amountColor
      },
      smooth: true,
        lineStyle: {
          color: amountColor,
          width: 3
        }
    },
    {
      name: 'Quantity',
      type: 'line',
      yAxisIndex: 1, // 对应左侧 Y 轴
      data: [],
      itemStyle: {
        color: warningColor
      },
      smooth: true,
        lineStyle: {
          color: warningColor,
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
                fill: chartMutedColor,
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
                        fill: chartMutedColor,
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
                        fill: chartMutedColor,
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
.dashboard-page {
  --bg: #020812;
  --panel: rgba(5, 16, 29, 0.86);
  --panel-soft: rgba(8, 24, 43, 0.62);
  --line: rgba(139, 214, 255, 0.16);
  --line-strong: rgba(139, 214, 255, 0.3);
  --text: #f4fbff;
  --muted: rgba(225, 242, 255, 0.72);
  --dim: rgba(225, 242, 255, 0.5);
  --cyan: #44d5ff;
  --green: #7dffcc;
  --danger: #ff8aa6;
  position: relative;
  width: 100%;
  min-height: calc(100vh - 110px);
  padding: 24px;
  overflow: hidden;
  color: var(--text);
  background:
    radial-gradient(circle at 10% 4%, rgba(68, 213, 255, 0.12), transparent 30%),
    radial-gradient(circle at 88% 0%, rgba(47, 125, 255, 0.13), transparent 28%),
    linear-gradient(180deg, #020812 0%, #030b16 58%, #04101d 100%);
  border-radius: 24px;
}

.dashboard-page::before {
  content: "";
  position: absolute;
  width: 420px;
  height: 420px;
  right: -180px;
  top: 120px;
  border-radius: 50%;
  background: rgba(68, 213, 255, 0.12);
  filter: blur(80px);
  pointer-events: none;
}

.dashboard-hero,
.metric-grid,
.chart-grid {
  position: relative;
  z-index: 1;
}

.dashboard-hero {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: end;
  gap: 22px;
  margin-bottom: 22px;
}

.eyebrow {
  display: inline-flex;
  margin: 0 0 12px;
  color: var(--green);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.dashboard-hero h1 {
  margin: 0;
  font-size: clamp(2rem, 4vw, 3.6rem);
  line-height: 1;
  letter-spacing: -0.055em;
}

.hero-copy {
  max-width: 620px;
  margin: 14px 0 0;
  color: var(--muted);
  line-height: 1.75;
}

.operation-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border: 1px solid var(--line);
  border-radius: 22px;
  background: rgba(4, 15, 29, 0.76);
  backdrop-filter: blur(14px);
}

.date-range {
  width: min(320px, 60vw);
}

.date-range :deep(.el-input__wrapper) {
  min-height: 42px;
  border-radius: 14px;
  border: 1px solid rgba(139, 214, 255, 0.14);
  background: rgba(255, 255, 255, 0.055);
  box-shadow: none;
}

.date-range :deep(.el-input__inner),
.date-range :deep(.el-range-input),
.date-range :deep(.el-range-separator) {
  color: var(--text);
}

.search-btn {
  min-height: 42px;
  border: 0;
  border-radius: 14px;
  color: #04111e;
  font-weight: 800;
  background: linear-gradient(135deg, var(--green), var(--cyan));
  box-shadow: 0 14px 28px rgba(68, 213, 255, 0.18);
}

.search-btn:hover,
.search-btn:focus {
  color: #04111e;
  background: linear-gradient(135deg, var(--green), var(--cyan));
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 18px;
}

.metric-card {
  min-height: 140px;
  padding: 22px;
  border: 1px solid var(--line);
  border-radius: 24px;
  background: var(--panel-soft);
  backdrop-filter: blur(14px);
}

.metric-card-primary {
  background:
    linear-gradient(135deg, rgba(125, 255, 204, 0.14), rgba(68, 213, 255, 0.08)),
    var(--panel);
  border-color: rgba(125, 255, 204, 0.28);
}

.metric-card span {
  color: var(--dim);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.metric-card strong {
  display: block;
  margin-top: 18px;
  color: #ffffff;
  font-size: clamp(1.6rem, 3vw, 2.5rem);
  line-height: 1;
  letter-spacing: -0.045em;
}

.metric-card p {
  margin: 12px 0 0;
  color: var(--muted);
  line-height: 1.55;
}

.chart-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.chart-card {
  min-width: 0;
  padding: 22px;
  border: 1px solid var(--line);
  border-radius: 28px;
  background:
    linear-gradient(145deg, rgba(68, 213, 255, 0.07), transparent 34%),
    var(--panel);
  box-shadow: 0 24px 70px rgba(0, 0, 0, 0.22);
  backdrop-filter: blur(14px);
}

.chart-card-wide {
  grid-column: 1 / -1;
}

.chart-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.chart-card-header span {
  display: grid;
  width: 42px;
  height: 42px;
  place-items: center;
  border-radius: 14px;
  color: #04111e;
  background: linear-gradient(135deg, var(--green), var(--cyan));
  font-weight: 900;
}

.chart-card-header h2 {
  margin: 0;
  color: var(--text);
  font-size: 1.05rem;
}

.chart {
  width: 100%;
  height: 430px;
}

.chart-card-wide .chart {
  height: 460px;
}

@media (max-width: 1180px) {
  .dashboard-hero {
    grid-template-columns: 1fr;
  }

  .operation-bar {
    width: fit-content;
  }

  .metric-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .chart-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 680px) {
  .dashboard-page {
    padding: 16px;
    border-radius: 18px;
  }

  .operation-bar {
    width: 100%;
    align-items: stretch;
    flex-direction: column;
  }

  .date-range {
    width: 100%;
  }

  .metric-grid {
    grid-template-columns: 1fr;
  }

  .chart-card {
    padding: 16px;
    border-radius: 22px;
  }

  .chart {
    height: 360px;
  }
}
</style>
