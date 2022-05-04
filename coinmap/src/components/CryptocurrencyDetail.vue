<template>
  <div>
    <h1>详细数据</h1>
    <el-table
      :data="tableData"
      style="text-algin:center; font-size: 15px; margin-top: 5px;"
      :header-cell-style="{ color: '#000000' }"
      :cell-style="{ color: '#000000' }"
      :row-style="{ height: 60 + 'px' }"
      @row-click="RowClick"
    >
      <el-table-column prop="Index" width="70" align="center" />
      <el-table-column prop="Name" label="Name" align="center"/>
      <el-table-column
        prop="Price"
        label="Price"
        align="center"
        style="margin-right:5%"
      />

      <el-table-column
        prop="MarketCap"
        label="Market Cap"
        align="center"
        :render-header="tipHelp1"
      />
      <el-table-column
        prop="Volume"
        label="Volume(24h)"
        align="center"
        :render-header="tipHelp2"
      />
    </el-table>
    <div class="block">
      <!-- 分页 -->
      <el-pagination
        @current-change="handleCurrentChange"
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="pager, jumper"
        style="margin-top: 2rem;text-align: center;margin-bottom: 2rem;"
        small
      >
      </el-pagination>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      coin: this.$route.query.coin,
      tableData: [],
      pageSize: 50,
      page: 1,
      total: 0,
      index: 0
    };
  },
  created() {
    this.getData();
  },
  methods: {
    tipHelp1() {
      return (
        <el-tooltip class="tooltip" effect="light" placement="bottom">
          //这里是提示语的具体内容
          <div slot="content">
            The total market value of a cryptocurrency's circulating supply.
          </div>
          <div slot="content">
            It is analogous to the free-float capitalization in the stock
            market.
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
      );
    },
    tipHelp2() {
      return (
        <el-tooltip class="tooltip" effect="light" placement="bottom">
          //这里是提示语的具体内容
          <div slot="content">
            A measure of how much of a cryptocurrency was traded in the last 24
            hours.
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
      );
    },
    async getData() {
      await this.$http
        .get(
          "data-api/v3/cryptocurrency/historical_page_num?coinName=" +
            this.coin,
          {
            headers: {
              token: window.sessionStorage["token"]
            }
          }
        )
        .then(res => {
          this.total = res.data;
        })
        .catch(function(error) {
          this.$message.error("获取数据失败");
        });
      this.index = (this.page - 1) * this.pageSize;
      await this.$http
        .get(
          "data-api/v3/cryptocurrency/historical_page?coinName=" +
            this.coin +
            "&pages=" +
            this.page +
            "&limits=" +
            this.pageSize,
          {
            headers: {
              token: window.sessionStorage["token"]
            }
          }
        )
        .then(res => {
          for (var i = 0; i < res.data.length; i++) {
            this.index++;
            this.tableData.push({
              Index: this.index,
              Name: res.data[i].Name,
              Id: res.data[i].Id,
              Price: "$" + res.data[i].HighPrice.toLocaleString(),
              Volume: "$" + parseInt(res.data[i].Volume).toLocaleString(),
              MarketCap: "$" + parseInt(res.data[i].MarketCap).toLocaleString(),
              d7: 0 + "%",
              h24: 0 + "%"
            });
          }
        })
        .catch(function(error) {
          this.$message.error("获取数据失败");
        });
    },
    // 分页跳转
    handleCurrentChange(val) {
      this.page = val;
      this.tableData = [];
      this.getData();
    },
    // 点击当前行传入要跳转模块的id
    RowClick(row) {
      this.$router.push({
        path: "/charts",
        query: {
          coin: row.Name
        }
      });
    }
  }
};
</script>

<style scoped>
</style>
