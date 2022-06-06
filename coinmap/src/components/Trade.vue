<template>
  <div>
    <div style="display:flex;gap:20px">
      <el-select
        style="flex:1;"
        v-model="selectedCoin"
        filterable
        placeholder="请选择币种"
        @change="getData()"
      >
        <el-option
          v-for="item in coins"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        >
        </el-option>
      </el-select>
      <el-button
        icon="el-icon-close"
        @click="
          selectedCoin = null;
          getData();
        "
      ></el-button>
    </div>

    <el-table
      :data="tableData"
      style="text-algin:center; font-size: 15px; margin-top: 20px;"
      :header-cell-style="{ color: '#000000' }"
      :cell-style="{ color: '#000000' }"
      :row-style="{ height: 60 + 'px' }"
      @row-click="jump_to_transaction"
    >
      <el-table-column label="创建日期" align="center" min-width="150px">
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
      <el-table-column label="价格" align="center" min-width="150px">
        <template slot-scope="scope">
          {{scope.row.Cost.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center">
        <template slot-scope="scope">
          <el-button v-if="scope.row.SellerId!=uid" @click="buyClick(scope.row.TsId)" type="small" @click.stop
            >购买</el-button
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
    let uid=JSON.parse(atob(window.sessionStorage.getItem("token").split('.')[1]))['uid'];
    return {
      uid,
      coins,
      selectedCoin: null,
      tableData: []
    };
  },
  created() {
    this.getData();
  },
  methods: {
    getData() {
      let cid = this.selectedCoin == null ? -1 : this.selectedCoin;
      this.$http
        .post(
          "data-api/v3/cryptocurrency/transaction/search",
          { cid },
          {
            headers: {
              token: window.sessionStorage["token"]
            }
          }
        )
        .then(res => {
          if (res.status == 200) {
            this.tableData = res.data;
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
    buyClick(TsId) {
      let oldThis = this;
      this.$confirm("确认购买吗", "", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          oldThis.$http
            .post(
              "data-api/v3/cryptocurrency/transaction/buy",
              { TsId },
              {
                headers: {
                  token: window.sessionStorage["token"]
                }
              }
            )
            .then(res => {
              if (res.data.code == 200) {
                oldThis.$message.success("购买成功");
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
            message: "取消购买"
          });
        });
    },
    jump_to_transaction(row){
      this.$router.push({name:'transaction',params: {TsId: row.TsId}});
    }
  }
};
</script>

<style></style>
