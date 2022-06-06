<template>
  <div>
    <el-table
      :data="tableData"
      style="text-algin:center;font-size: 15px;margin-top: 5px;"
      :header-cell-style="{ color: '#000000' }"
      :cell-style="{ color: '#000000' }"
      :row-style="{ height: 70 + 'px' }"
      @row-click="RowClick"
    >
      //需要显示文字提示的表头列加上 :render-header="tipHelp"，tipHelp是一个方法
      <el-table-column prop="Index" width="50" align="center">
      </el-table-column>
      <el-table-column prop="Name" label="Name" align="center" min-width="150px"/>
      <el-table-column
        prop="Price"
        label="Price"
        align="center"
        style="margin-right:5%"
        min-width="150px"
      >
      </el-table-column>
      <el-table-column prop="h24" label="24h %" align="center">
        <template slot-scope="scope">
          <span
            v-if="parseFloat(scope.row.h24) >= parseFloat('0%')"
            style="color: green"
            >{{ scope.row.h24 }}</span
          >
          <span v-else style="color: red">{{ scope.row.h24 }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="d7" label="7d %" align="center">
        <template slot-scope="scope">
          <span
            v-if="parseFloat(scope.row.d7) >= parseFloat('0%')"
            style="color: green"
            >{{ scope.row.d7 }}</span
          >
          <span v-else style="color: red">{{ scope.row.d7 }}</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="MarketCap"
        label="Market Cap"
        align="center"
        :render-header="tipHelp1"
         min-width="170px"
      >
      </el-table-column>
      <el-table-column
        prop="Volume"
        label="Volume(24h)"
        align="center"
        :render-header="tipHelp2"
         min-width="170px"
      >
      </el-table-column>
      <el-table-column prop="Collection" label="Collection" align="center">
        <template slot-scope="scope">
          <el-button
            @click="handleClick(scope.row)"
            @click.stop
            type="danger"
            icon="el-icon-star-off"
            circle
          ></el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import coins from "../coins";
export default {
  data() {
    return {
      cdata: [],
      username: this.$route.query.username,
      x: 0,
      tableData: [],
      ismultiple: false,
      page: 1,
      pagesize: 20,
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
      var _this = this;
      await this.$http
        .get("data-api/v3/cryptocurrency/like", {
          headers: {
            token: window.sessionStorage["token"]
          }
        })
        .then(res => {
          if (res.status === 200) {
            for (var i = 0; i < res.data.length; i++) {
              _this.$http
                .get(
                  "data-api/v3/cryptocurrency/historical?coinName=" +
                    res.data[i] +
                    "&timeStart=1630936389&timeEnd=1630936389",
                  {
                    headers: {
                      token: window.sessionStorage["token"]
                    }
                  }
                )
                .then(res1 => {
                  _this.index++;
                  if (res1.status === 200) {
                    _this.tableData.push({
                      Index: _this.index,
                      Name: res1.data[0].Name,
                      Id: res1.data[0].Id,
                      Price: "$" + res1.data[0].HighPrice.toLocaleString(),
                      Volume:
                        "$" + parseInt(res1.data[0].Volume).toLocaleString(),
                      MarketCap:
                        "$" + parseInt(res1.data[0].MarketCap).toLocaleString(),
                      d7: 0 + "%",
                      h24: 0 + "%"
                    });
                  }
                })
                .catch(function(error) {
                  _this.$message.error("请求失败");
                });
            }
          }
        });
    },
    // 点击当前行传入要跳转模块的id
    RowClick(row) {
      this.$router.push({
        path: "/charts",
        query: {
          coin: row.Name
        }
      });
    },
    handleClick(row) {
      var _this = this;
      this.$http
        .delete(
          "data-api/v3/cryptocurrency/like?cid=" +
            coins.find(x => x.name == row.Name).id,
          {
            headers: {
              token: window.sessionStorage["token"]
            }
          }
        )
        .then(function(data) {
          _this.$message.success("取消收藏成功");
          _this.tableData = [];
          _this.index = 0;
          _this.getData();
        });
    }
  }
};
</script>

<style scoped>
</style>
