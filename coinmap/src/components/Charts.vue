<template>
  <div>
    <div id="flex">
        <div class="flex-item">
          <el-button size="mini" class="but" @click="click1">价格</el-button>
          <el-button size="mini" class="but" @click="click2">市值</el-button>
          <el-button size="mini" class="but" @click="click3">交易量</el-button>
          <el-button size="mini" class="but" @click="click4">比特币兑换</el-button>
        </div>
        <div class="flex-item">
          <el-button size="mini" class="but" @click="click9">USD</el-button>
          <el-button size="mini" class="but" @click="click10">CNY</el-button>
        </div>
        <div class="flex-item">
          <el-button size="mini" class="but" @click="click5">1D</el-button>
          <el-button size="mini" class="but" @click="click6">7D</el-button>
          <el-button size="mini" class="but" @click="click7">1M</el-button>
          <el-button size="mini" class="but" @click="click8">1Y</el-button>
        </div>
    </div>
    <div class="EchartPractice">
      <div id="main"></div>
    </div>
    
  </div>
</template>

<script>
export default {
  name: "EchartPractice",
  data() {
    return {
      coin: this.$route.query.coin,
      mode: "Price",
      date: "1Y",
      money: "$",
    };
  },
  methods: {
    click1() {
      this.mode = "Price";
      var myChart = document.getElementById("main");
      myChart.removeAttribute("_echarts_instance_");
      this.drawChartchoose();
    },
    click2() {
      this.mode = "MarketCap";
      var myChart = document.getElementById("main");
      myChart.removeAttribute("_echarts_instance_");
      this.drawChartchoose();
    },
    click3() {
      this.mode = "TradingView";
      var myChart = document.getElementById("main");
      myChart.removeAttribute("_echarts_instance_");
      this.drawChartchoose();
    },
    click4() {
      this.mode = "BtcExchange";
      var myChart = document.getElementById("main");
      myChart.removeAttribute("_echarts_instance_");
      this.drawChartchoose();
    },
    click5() {
      this.date = "1D";
      var myChart = document.getElementById("main");
      myChart.removeAttribute("_echarts_instance_");
      this.drawChartchoose();
    },
    click6() {
      this.date = "7D";
      var myChart = document.getElementById("main");
      myChart.removeAttribute("_echarts_instance_");
      this.drawChartchoose();
    },
    click7() {
      this.date = "1M";
      var myChart = document.getElementById("main");
      myChart.removeAttribute("_echarts_instance_");
      this.drawChartchoose();
    },
    click8() {
      this.date = "1Y";
      var myChart = document.getElementById("main");
      myChart.removeAttribute("_echarts_instance_");
      this.drawChartchoose();
    },
    click9() {
      this.money = "$";
      var myChart = document.getElementById("main");
      myChart.removeAttribute("_echarts_instance_");
      this.drawChartchoose();
    },
    click10() {
      this.money = "￥";
      var myChart = document.getElementById("main");
      myChart.removeAttribute("_echarts_instance_");
      this.drawChartchoose();
    },
    drawChartchoose() {
      if (this.mode == "Price") {
        this.drawChart();
      } else if (this.mode == "MarketCap") {
        this.drawChart2();
      } else if (this.mode == "TradingView") {
        this.drawChart3();
      } else if (this.mode == "BtcExchange") {
        this.drawChart4();
      }
    },
    async drawChart() {
      let myEchart = this.$echarts.init(document.getElementById("main"));
      let data = [];
      await this.$http
        .get(
          "data-api/v3/cryptocurrency/detail/chart?coinName=" +
            this.coin +
            "&range=" +
            this.date +
            "&convertId=2781",
          {
            headers: {
              token: window.sessionStorage["token"],
            },
          }
        )
        .then((res) => {
          res.data.map((i) => {
            if (this.money === "$") {
              data.push([parseInt(i.Time + "000"), i.Price]);
            }
            if (this.money === "￥") {
              data.push([parseInt(i.Time + "000"), i.ZhPrice]);
            }
          });
        });

      let option = {
        tooltip: {
          trigger: "axis",
          position: function (pt) {
            return [pt[0], "10%"];
          },
        },
        title: {
          left: "center",
          text: this.coin + " Chart",
        },
        toolbox: {
          feature: {
            dataZoom: {
              yAxisIndex: "none",
            },
            restore: {},
            saveAsImage: {},
          },
        },
        xAxis: {
          type: "time",
          boundaryGap: false,
        },
        yAxis: {
          scale: true,
          type: "value",
          boundaryGap: [0, "100%"],
        },
        dataZoom: [
          {
            type: "inside",
            start: 0,
            end: 20,
          },
          {
            start: 0,
            end: 20,
          },
        ],
        series: [
          {
            name: this.money + "Price:",
            type: "line",
            smooth: true,
            symbol: "none",
            areaStyle: {},
            data: data,
          },
        ],
        visualMap: {
          show: true,
          pieces: [
            {
              gt: 0,
              lte: data[0][1],
              color: "red",
            },
            {
              gt: data[0][1],
              color: "green",
            },
          ],
          seriesIndex: 0,
        },
      };
      myEchart.setOption(option);
    },
    async drawChart2() {
      let myEchart = this.$echarts.init(document.getElementById("main"));
      // let base = +new Date(1988, 9, 3);
      // let oneDay = 24 * 3600 * 1000;
      // let data = [[base, Math.random() * 300]];
      // for (let i = 1; i < 20000; i++) {
      //   let now = new Date((base += oneDay));
      //   data.push([+now, Math.round((Math.random() - 0.5) * 20 + data[i - 1][1])]);
      // }
      let data = [];
      await this.$http
        .get(
          "data-api/v3/cryptocurrency/detail/chart?coinName=" +
            this.coin +
            "&range=" +
            this.date +
            "&convertId=2781",
          {
            headers: {
              token: window.sessionStorage["token"],
            },
          }
        )
        .then((res) => {
          res.data.map((i) => {
            if (this.money === "$") {
              data.push([parseInt(i.Time + "000"), i.MarketCap]);
            }
            if (this.money === "￥") {
              data.push([parseInt(i.Time + "000"), i.ZhMarketCap]);
            }
          });
        });

      let option = {
        tooltip: {
          trigger: "axis",
          position: function (pt) {
            return [pt[0], "10%"];
          },
        },
        title: {
          left: "center",
          text: this.coin + " Chart",
        },
        toolbox: {
          feature: {
            dataZoom: {
              yAxisIndex: "none",
            },
            restore: {},
            saveAsImage: {},
          },
        },
        xAxis: {
          type: "time",
          boundaryGap: false,
        },
        yAxis: {
          scale: true,
          type: "value",
          boundaryGap: [0, "100%"],
        },
        dataZoom: [
          {
            type: "inside",
            start: 0,
            end: 20,
          },
          {
            start: 0,
            end: 20,
          },
        ],
        series: [
          {
            name: this.money + "MarketCap",
            type: "line",
            smooth: true,
            symbol: "none",
            areaStyle: {},
            data: data,
          },
        ],
      };
      myEchart.setOption(option);
    },
    async drawChart3() {
      let myEchart = this.$echarts.init(document.getElementById("main"));
      const upColor = "#ec0000";
      const upBorderColor = "#8A0000";
      const downColor = "#00da3c";
      const downBorderColor = "#008F28";
      // const dataCount = 2e5;
      // const data = generateOHLC(dataCount);
      let data = [];
      await this.$http
        .get(
          "data-api/v3/cryptocurrency/historical?coinName=" +
            this.coin +
            "&timeStart=0&timeEnd=1630936389",
          {
            headers: {
              token: window.sessionStorage["token"],
            },
          }
        )
        .then((res) => {
          res.data.map((i) => {
            if (this.money === "$") {
              data.push([
                i.TimeHigh,
                i.OpenPrice,
                i.HighPrice,
                i.LowPrice,
                i.ClosePrice,
                i.MarketCap,
              ]);
            }
            if (this.money === "￥") {
              data.push([
                i.TimeHigh,
                i.ZhOpenPrice,
                i.ZhHighPrice,
                i.ZhLowPrice,
                i.ZhClosePrice,
                i.ZhMarketCap,
              ]);
            }
          });
        });

      let option = {
        dataset: {
          source: data,
        },
        title: {
          left: "center",
          text: this.coin + " Chart",
        },
        tooltip: {
          trigger: "axis",
          axisPointer: {
            type: "line",
          },
        },
        toolbox: {
          feature: {
            dataZoom: {
              yAxisIndex: false,
            },
          },
        },
        grid: [
          {
            left: "10%",
            right: "10%",
            bottom: 200,
          },
          {
            left: "10%",
            right: "10%",
            height: 80,
            bottom: 80,
          },
        ],
        xAxis: [
          {
            type: "category",
            scale: true,
            boundaryGap: false,
            // inverse: true,
            axisLine: { onZero: false },
            splitLine: { show: false },
            min: "dataMin",
            max: "dataMax",
          },
          {
            type: "category",
            gridIndex: 1,
            scale: true,
            boundaryGap: false,
            axisLine: { onZero: false },
            axisTick: { show: false },
            splitLine: { show: false },
            axisLabel: { show: false },
            min: "dataMin",
            max: "dataMax",
          },
        ],
        yAxis: [
          {
            splitArea: {
              show: true,
            },
          },
          {
            scale: true,
            gridIndex: 1,
            splitNumber: 2,
            axisLabel: { show: false },
            axisLine: { show: false },
            axisTick: { show: false },
            splitLine: { show: false },
          },
        ],
        dataZoom: [
          {
            type: "inside",
            xAxisIndex: [0, 1],
            start: 10,
            end: 100,
          },
          {
            show: true,
            xAxisIndex: [0, 1],
            type: "slider",
            bottom: 10,
            start: 10,
            end: 100,
          },
        ],
        visualMap: {
          show: false,
          seriesIndex: 1,
          dimension: 6,
          pieces: [
            {
              value: 1,
              color: upColor,
            },
            {
              value: -1,
              color: downColor,
            },
          ],
        },
        series: [
          {
            type: "candlestick",
            itemStyle: {
              color: upColor,
              color0: downColor,
              borderColor: upBorderColor,
              borderColor0: downBorderColor,
            },
            encode: {
              x: 0,
              y: [1, 4, 3, 2],
            },
          },
          {
            name: "Volumn",
            type: "bar",
            xAxisIndex: 1,
            yAxisIndex: 1,
            itemStyle: {
              color: "#7fbe9e",
            },
            large: true,
            encode: {
              x: 0,
              y: 5,
            },
          },
        ],
      };
      myEchart.setOption(option);
    },
    async drawChart4() {
      let myEchart = this.$echarts.init(document.getElementById("main"));
      // let base = +new Date(1988, 9, 3);
      // let oneDay = 24 * 3600 * 1000;
      // let data = [[base, Math.random() * 300]];
      // for (let i = 1; i < 20000; i++) {
      //   let now = new Date((base += oneDay));
      //   data.push([+now, Math.round((Math.random() - 0.5) * 20 + data[i - 1][1])]);
      // }
      let data = [];
      await this.$http
        .get(
          "data-api/v3/cryptocurrency/detail/chart?coinName=" +
            this.coin +
            "&range=" +
            this.date +
            "&convertId=2781",
          {
            headers: {
              token: window.sessionStorage["token"],
            },
          }
        )
        .then((res) => {
          res.data.map((i) => {
            if (this.money === "$") {
              data.push([parseInt(i.Time + "000"), i.BitcoinRate]);
            }
            if (this.money === "￥") {
              data.push([parseInt(i.Time + "000"), i.BitcoinRate]);
            }
          });
        });

      let option = {
        tooltip: {
          trigger: "axis",
          position: function (pt) {
            return [pt[0], "10%"];
          },
        },
        title: {
          left: "center",
          text: this.coin + " Chart",
        },
        toolbox: {
          feature: {
            dataZoom: {
              yAxisIndex: "none",
            },
            restore: {},
            saveAsImage: {},
          },
        },
        xAxis: {
          type: "time",
          boundaryGap: false,
        },
        yAxis: {
          scale: true,
          type: "value",
          boundaryGap: [0, "100%"],
        },
        dataZoom: [
          {
            type: "inside",
            start: 0,
            end: 20,
          },
          {
            start: 0,
            end: 20,
          },
        ],
        series: [
          {
            name: "ExchangeRate",
            type: "line",
            smooth: true,
            symbol: "none",
            areaStyle: {},
            data: data,
            areaStyle: {
              normal: { color: "#ffffff" }, //改变区域颜色
            },
          },
        ],
      };
      myEchart.setOption(option);
    },
  },
  mounted() {
    this.drawChartchoose();
  },
};
</script>

<style scoped>
#main {
  width: inherit;
  height: 400px;
  margin: auto;
  margin-top: 10px;
}
.but {
  margin: 0;
}
#flex {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  align-content: center;
  justify-content: space-around;
}
.flex-item {
  margin-bottom: 1rem;
  margin-top: 1rem;
}
</style>
