<template>
  <div>
    <h1>实时价格</h1>
    <el-table
      :data="tableData"
      style="text-algin:center; font-size: 15px;margin-top: 5px;"
      :header-cell-style="{ color: '#000000' }"
      :cell-style="{ color: '#000000' }"
      :row-style="{ height: 60 + 'px' }"
      @row-click="RowClick"
    >
      <el-table-column prop="Index" width="50" align="center"/>
      <el-table-column prop="Name" label="Name" align="center"/>
      <el-table-column
        prop="Price"
        label="Price"
        align="center"
        style="margin-right:5%"
      />
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
      >
      </el-table-column>
      <el-table-column
        prop="Volume"
        label="Volume(24h)"
        align="center"
        :render-header="tipHelp2"
      >
      </el-table-column>
      <el-table-column prop="Collection" label="Collection" align="center">
        <template slot-scope="scope">
          <el-button
            @click="lickClick(scope.row)"
            @click.stop
            type="primary"
            icon="el-icon-star-off"
            circle
          ></el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="block">
      <!-- 分页 -->
      <el-pagination
        @current-change="handleCurrentChange"
        :current-page="page"
        :page-size="20"
        layout="pager, jumper"
        :total="544"
        style="margin-top: 2rem;text-align: center;margin-bottom: 2rem;"
        small
      >
      </el-pagination>
    </div>
  </div>
</template>

<script>
import coins from '../coins'
export default {
  data() {
    return {
      tableData: [],
      page: 1,
      pagesize: 20,
      index: 0,
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
      this.params = [];
      this.index = (this.page - 1) * this.pagesize;
      for (var i = this.index; i < this.page * this.pagesize; i++) {
        await this.$http
          .get(
            "data-api/v3/cryptocurrency/historical?coinName=" +
              coins[i]["name"] +
              "&timeStart=1630936389&timeEnd=1630936389",
            {
              headers: {
                token: window.sessionStorage["token"]
              }
            }
          )
          .then(res => {
            this.index++;
            if (res.status === 200) {
              this.tableData.push({
                Index: this.index,
                Name: res.data[0].Name,
                Id: res.data[0].Id,
                Price: "$" + res.data[0].HighPrice.toLocaleString(),
                Volume: "$" + parseInt(res.data[0].Volume).toLocaleString(),
                MarketCap:
                  "$" + parseInt(res.data[0].MarketCap).toLocaleString(),
                d7: 0 + "%",
                h24: 0 + "%"
              });
              this.params.push(res.data[0].Id);
            }
          })
          .catch(() => {
            // this.$message.error("获取数据失败");
          });
      }
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
        path: "/cryptocurrency_detail",
        query: {
          coin: row.Name
        }
      });
    },
    lickClick(row) {
      var param = new FormData();
      param.append("cid", coins.find((x)=>x.name==row.Name).id);
      var config = {
        headers: {
          token: window.sessionStorage["token"]
        }
      };
      this.$http
        .post("data-api/v3/cryptocurrency/like", param, config)
        .then(data => {
          this.$message.success("收藏成功");
        })
        .catch(e => {
          this.$message.error("收藏失败: "+e.response.data.msg);
        });
    }
  }
};
</script>

<style scoped>
</style>
