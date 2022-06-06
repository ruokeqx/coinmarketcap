<template>
  <div>

      <el-menu
        :default-active="'1'"
        mode="horizontal"
        @select="handleSelect"
      >
      <el-menu-item index="1">我卖的</el-menu-item>
      <el-menu-item index="2">我买的</el-menu-item>
  </el-menu>
    <el-table
      :data="menuSelected==1?sellTableData:buyTableData"
      style="text-algin:center; font-size: 15px; margin-top: 2rem"
      :header-cell-style="{ color: '#000000' }"
      :cell-style="{ color: '#000000' }"
      :row-style="{ height: 60 + 'px' }"
      @row-click="jump_to_transaction"
    >
      <el-table-column :label="menuSelected==1?'创建日期':'购买日期'" align="center" min-width="150px">
        <template slot-scope="scope">
          {{ dateToString(scope.row.TsCreaTime) }}
        </template>
      </el-table-column>
      <el-table-column label="交易币种" align="center">
        <template slot-scope="scope">
          {{ coins.find(x => x.id == scope.row.TsCid).name }}
        </template>
      </el-table-column>
      <el-table-column prop="TsNum" label="交易数量" align="center" />
      <el-table-column prop="Discount" label="折扣" align="center" />
      <el-table-column label="价格" align="center" min-width="200px">
        <template slot-scope="scope">
          {{scope.row.Cost.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" v-if="menuSelected==1" min-width="150px">
        <template slot-scope="scope">
          <el-button v-if="scope.row.TsStatus==0" @click="takeDown(scope.row.TsId)" type="small" style="margin-right: 5px" @click.stop
            >下架</el-button
          >
          <el-button v-if="scope.row.TsStatus==0" @click="discount(scope.row.TsId)" type="small" @click.stop
            >打折</el-button
          >
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
      coins,
      sellTableData: [],
      buyTableData: [],
      menuSelected:1,
    };
  },
  created() {
    this.getData();
  },
  methods: {
    getData() {
      this.$http
        .get(
          "data-api/v3/cryptocurrency/transaction/mysell",
          {
            headers: {
              token: window.sessionStorage["token"]
            }
          }
        )
        .then(res => {
          if (res.status == 200) {
            this.sellTableData = res.data;
          } else {
            this.$message.error("获取数据失败");
          }
        })
        .catch(() => {
          this.$message.error("获取数据失败");
        });
      this.$http
        .get(
          "data-api/v3/cryptocurrency/transaction/mybuy",
          {
            headers: {
              token: window.sessionStorage["token"]
            }
          }
        )
        .then(res => {
          if (res.status == 200) {
            this.buyTableData = res.data;
          } else {
            this.$message.error("获取数据失败");
          }
        })
        .catch(() => {
          this.$message.error("获取数据失败");
        });
    },
    dateToString(timeStamp) {
      let date = new Date(timeStamp * 1000);
      const toTwoStr=(x)=>x<10?"0"+x:""+x;
      return (
        date.getFullYear() +
        "-" +
        toTwoStr(date.getMonth() + 1) +
        "-" +
        toTwoStr(date.getDate()) +
        " " +
        toTwoStr(date.getHours()) +
        ":" +
        toTwoStr(date.getMinutes()) // +
        // ":" +
        // toTwoStr(date.getSeconds())
      );
    },
    takeDown(TsId) {
      let oldThis = this;
      this.$confirm("确认下架吗", "", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          oldThis.$http
            .post(
              "data-api/v3/cryptocurrency/transaction/close",
              { TsId },
              {
                headers: {
                  token: window.sessionStorage["token"]
                }
              }
            )
            .then(res => {
              if (res.data.code == 200) {
                oldThis.$message.success("下架成功");
                oldThis.getData();
              } else {
                oldThis.$alert(res.data.msg, '提示', {confirmButtonText: 'OK'});
              }
            })
            .catch((err) => {
              oldThis.$alert("错误: "+err.response.data.msg, '提示', {confirmButtonText: 'OK'});
            });
        })
        .catch(() => {
          oldThis.$message({
            type: "info",
            message: "取消下架"
          });
        });
    },
    discount(TsId){
      let oldThis = this;
      this.$prompt("请输入折扣(0-1之间)", "", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        inputPattern: /0(\.[0-9]*)?/,
        inputErrorMessage: "折扣格式错误"
      })
        .then(({ value }) => {
          oldThis.$http
            .post(
              "data-api/v3/cryptocurrency/transaction/discount",
              { TsId, discount: parseFloat(value) },
              {
                headers: {
                  token: window.sessionStorage["token"]
                }
              }
            )
            .then(res => {
              if (res.data.code == 200) {
                oldThis.$message.success("打折成功");
                oldThis.getData();
              } else {
                oldThis.$alert(res.data.msg, '提示', {confirmButtonText: 'OK'});
              }
            })
            .catch((err) => {
              oldThis.$alert("错误: "+err.response.data.msg, '提示', {confirmButtonText: 'OK'});
            });
        })
        .catch(() => {
        });
    },
    handleSelect (key) {
      this.menuSelected=parseInt(key)
    },
    jump_to_transaction(row){
      this.$router.push({name:'transaction',params: {TsId: row.TsId}});
    }
  }
};
</script>

<style></style>
