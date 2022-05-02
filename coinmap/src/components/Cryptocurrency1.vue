<template>
<div>
  
  <el-row class = 'row'  style="width:80%;text-algin:center; margin-left: 10%; margin-top: 10px">
    <el-col :span="12"><div class="grid-content bg-purple" style="font-weight: bolder;
      font-size: 18px;">Today's Cryptocurrency Prices by Market Cap</div></el-col>
  </el-row>
  <el-table :data="tableData"  style="width:80%;text-algin:center; margin-left: 10%; font-size: 15px" :header-cell-style="{color: '#000000'}" :cell-style="{color: '#000000'}" :row-style="{height:70+'px'}" >
     //需要显示文字提示的表头列加上 :render-header="tipHelp"，tipHelp是一个方法
     <el-table-column
      prop="Index"
      width="50"
      align="center"
      >
    </el-table-column>
    <!-- <el-table-column
      prop="Name"
      label="Name"
      align="center"
    >
    </el-table-column> -->
    <el-table-column prop="Name" label="Name" align="center">
      <template slot-scope="scope">
        <div
          size="mini"
          id="main"
          @click="RowClick(scope.row.Name)">{{scope.row.Name}}</div>
      </template>
    </el-table-column>
    <el-table-column
      prop="Price"
      label="Price"
      align="center"
      style="margin-right:5%"
    >
    </el-table-column>
    <!-- <el-table-column
      prop="h24"
      label="24h %"
      align="center"
    >
    </el-table-column> -->
    
    <el-table-column
      prop="MarketCap"
      label="Market Cap"
      align="center"
      :render-header="tipHelp1"
    >
    </el-table-column>
    <el-table-column
      prop="Volume"
      label="Volume(24h)"
      align="center"
      :render-header="tipHelp2"
    >
    </el-table-column>
    <!-- <el-table-column
      prop="Name"
      label="Last 7 Days"
      align="center"
    >
    <template slot-scope="scope">
      <div>
        {{ drawChart4(scope.row) }}
        <div :id="scope.row.Index" class="tiger-trend-charts"></div>
      </div>
    </template>
    </el-table-column> -->
  </el-table>
    <div class="block">
    <!-- 分页 -->
    <el-pagination  
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :current-page="page"
      :page-sizes="[20, 50, 80 ,100]"
      :page-size="pagesize"
      layout="total, sizes, prev, pager, next, jumper"
      style="margin-left: 25%; margin-top:10px;"
      background>
    </el-pagination>
    </div>
   <el-backtop :top="50">
    <div
      style="{
        height: 100%;
        width: 100%;
        background-color: #f2f5f6;
        box-shadow: 0 0 6px rgba(0,0,0, .12);
        text-align: center;
        line-height: 40px;
        color: #1989fa;
      }"
    >
      ⬆
    </div>
  </el-backtop>
</div>
</template>

<script>
import { staticArray } from '../../public/coins'
import { GetParam1 } from '../../public/page1-id'
import { GetParam2 } from '../../public/page2-id'
import { GetParam3 } from '../../public/page3-id'
import { GetParam4 } from '../../public/page4-id'
import { GetParam5 } from '../../public/page5-id'
import { GetParam6 } from '../../public/page6-id'
import wsUrl from '../config.js'
export default {
  data () {
    return {
      coin: this.$route.query.coin,
      x: 0,
      coins: staticArray(),
      tableData: [],
      ismultiple: false,
      page: 1,
      pagesize:20,
      index: 0,
      path: wsUrl,
      socket: '',
      params: [],
      param1: GetParam1(),
      param2: GetParam2(),
      param3: GetParam3(),
      param4: GetParam4(),
      param5: GetParam5(),
      param6: GetParam6(),
      socdata: [],
      cstyle: { color: '#000000' },
    }
  },
  created () {
    this.getData()
  },
  mounted () {
    // 初始化
    this.init()
  },
  methods: {
    tipHelp1 () {
      return (
        <el-tooltip class="tooltip" effect="light" placement="bottom">
        //这里是提示语的具体内容
        <div slot="content">
        The total market value of a cryptocurrency's circulating supply.
        </div>
        <div slot="content">
        It is analogous to the free-float capitalization in the stock market.
        </div>
        <div slot="content">
        Market Cap = Current Price x Circulating Supply.
        </div>
        //这里是表头的名称，并加上一个icon
        <div>
        Market Cap
        <i
          class="el-icon-question"
          style="margin-left:3px;font-size:15px;"
        />
        </div>
        </el-tooltip>
      )
    },
    tipHelp2 () {
      return (
        <el-tooltip class="tooltip" effect="light" placement="bottom">
        //这里是提示语的具体内容
        <div slot="content">
        A measure of how much of a cryptocurrency was traded in the last 24 hours.
        </div>
        //这里是表头的名称，并加上一个icon
        <div>
        Volume(24h)
        <i
          class="el-icon-question"
          style="margin-left:3px;font-size:15px;"
        />
        </div>
        </el-tooltip>
      )
    },
    async getData () {
        this.index=(this.page-1)*this.pagesize
        await this.$http.get('data-api/v3/cryptocurrency/historical_page?coinName=' + this.coin + '&page=' + this.page + '&limits=' + this.pagesize , {
          headers: {
            'token':  window.sessionStorage['token']
          }
        }).then(res => {
          console.log(res.data)
          for(var i=0;i<res.data.length;i++){
              this.index++
              this.tableData.push({ Index: this.index, Name: res.data[i].Name, Id: res.data[i].Id, Price: '$' + res.data[i].HighPrice.toLocaleString(), Volume: '$' + parseInt(res.data[i].Volume).toLocaleString(), MarketCap: '$' + parseInt(res.data[i].MarketCap).toLocaleString(), d7: 0 + '%', h24: 0 + '%' })
          }
        })
        .catch(function(error){
          console.log('请求失败')
          }
        )
      // [1,1027,1839,52,5994]
      // this.params.push(1)
      // this.params.push(1027)
      // this.params.push(1839)
      // this.params.push(52)
      // this.params.push(5994)
      console.log(this.tableData)
      console.log(this.params)
    },
    init: function () {
      if (typeof (WebSocket) === 'undefined') {
        alert('您的浏览器不支持socket')
      } else {
        // 实例化socket
        this.socket = new WebSocket(this.path)
        // 监听socket连接
        this.socket.onopen = this.open
        // 监听socket错误信息
        this.socket.onerror = this.error
        // 监听socket消息
        this.socket.onmessage = this.getMessage
      }
    },
    open: function () {
      console.log('socket连接成功')
      if (this.page === 1) {
        this.send('[' + this.param1 + ']')
      }
      if (this.page === 2) {
        this.send('[' + this.param2 + ']')
      }
      if (this.page === 3) {
        this.send('[' + this.param3 + ']')
      }
      if (this.page === 4) {
        this.send('[' + this.param4 + ']')
      }
      if (this.page === 5) {
        this.send('[' + this.param5 + ']')
      }
      if (this.page === 6) {
        this.send('[' + this.param6 + ']')
      }
    },
    error: function () {
      console.log('连接错误')
    },
    getMessage: function (msg) {
      // console.log(JSON.parse(msg.data).d.cr)
      this.socdat = JSON.parse(msg.data).d.cr
      console.log(this.socdat)
      this.tableData.forEach(element => {
        if (element.Id === this.socdat.id) {
          element.d7 = this.socdat.p7d.toLocaleString() + '%'
          element.h24 = this.socdat.p24h.toLocaleString() + '%'
          element.Price = '$' + this.socdat.p.toLocaleString()
          element.MarketCap = '$' + this.socdat.mc.toLocaleString()
          element.Volume = '$' + this.socdat.v.toLocaleString()
        }
      })
    },
    send: function (params) {
      console.log('发送' + params)
      this.socket.send(params)
    },
    close: function () {
      console.log('socket已经关闭')
    },
    // 每页数据的数量
    handleSizeChange (val) {
      console.log(`每页 ${val} 条`)
      this.pagesize=val
      this.tableData = []
      this.getData()
      this.init()
    },
    // 分页跳转
    handleCurrentChange (val) {
      console.log(`当前页: ${val}`)
      this.page = val
      this.tableData = []
      this.getData()
      this.init()
    },
     // 点击当前行传入要跳转模块的id
    RowClick (name) {
      console.log(name)
      this.$router.push({
        path: '/charts',
        query: {
          coin: name
        }
      })
    },
    handleClick(row){
      var param = new FormData();
      param.append('coinName',row.Name)
      var config ={
        headers: {
            'token':  window.sessionStorage['token']
          },
      }
      console.log(row.Name)
      this.$message.success('收藏成功');
      this.$http.post('data-api/v3/cryptocurrency/like',param,config).then(function(data){
        console.log(data)
      })
    },
    scClick(){
      console.log("sc")
      this.$router.push({
        path: '/collection',
        query: {
          username: this.username
        }
      })
    },
    hbClick(){
      console.log("hb")
    },
    gmClick(){
      console.log("gm")
    }
    // async drawChart4(Row) {
    //     let myEchart = this.$echarts.init(document.getElementById(Row.Index));
    //     // let base = +new Date(1988, 9, 3);
    //     // let oneDay = 24 * 3600 * 1000;
    //     // let data = [[base, Math.random() * 300]];
    //     // for (let i = 1; i < 20000; i++) {
    //     //   let now = new Date((base += oneDay));
    //     //   data.push([+now, Math.round((Math.random() - 0.5) * 20 + data[i - 1][1])]);
    //     // }
    //     // console.log(data)
    //     let data = [];
    //     await this.$http.get('data-api/v3/cryptocurrency/detail/chart?coinName=' + Row.Name + '&range=1M&convertId=2781').then(res => {
    //       console.log(res.data)
    //       res.data.map(i => {
    //         if(this.money === '2781'){
    //           data.push([parseInt(i.Time+'000'), i.Price]);
    //         }
    //         if(this.money === '2787'){
    //           data.push([parseInt(i.Time+'000'), i.ZhPrice]);
    //         }
    //       });
    //       console.log(data)
    //     })
        
    //     let option = {
    //       tooltip: {
    //         trigger: 'axis',
    //         position: function (pt) {
    //           return [pt[0], '10%'];
    //         }
    //       },
    //       title: {
    //         left: 'center',
    //         text: this.coin+' Chart'
    //       },
    //       toolbox: {
    //         feature: {
    //           dataZoom: {
    //             yAxisIndex: 'none'
    //           },
    //           restore: {},
    //           saveAsImage: {}
    //         }
    //       },
    //       xAxis: {
    //         type: 'time',
    //         boundaryGap: false
    //       },
    //       yAxis: {
    //         scale:true,
    //         type: 'value',
    //         boundaryGap: [0, '100%']
    //       },
    //       dataZoom: [
    //         {
    //           type: 'inside',
    //           start: 0,
    //           end: 20
    //         },
    //         {
    //           start: 0,
    //           end: 20
    //         }
    //       ],
    //       series: [
    //         {
    //           name: 'ExchangeRate',
    //           type: 'line',
    //           smooth: true,
    //           symbol: 'none',
    //           areaStyle: {},
    //           data: data,
    //           areaStyle : {
    //             normal : {color : '#ffffff',}//改变区域颜色
    //           }
    //         }
    //       ],
    //     };
    //     myEchart.setOption(option);
    //   }
  },
  destroyed () {
    // 销毁监听
    this.socket.onclose = this.close
  }
}
// 其他table的属性和方法可根据实际情况对应地去使用，用法大多是大同小异的，可以去官网查看文档喔
</script>

<style scoped>
.el-table{
  margin-top: 5px;
}
.user{
  margin-left: 83%;
}
</style>
