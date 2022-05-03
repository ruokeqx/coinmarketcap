<template>
  <div>
    <h1>个人资产</h1>
    <div>
      <el-card shadow="hover" class="my-card">
        <div
          style="display:flex;align-items: flex-start;flex-wrap: nowrap;justify-content: space-between;"
        >
          <div class="highlight-title" style="display: inline;">账户余额</div>
          <el-button type="small" @click="charge(-1)">充值</el-button>
        </div>
        <span class="num-large">{{ Math.floor(rmbAmount) }}</span
        ><span class="num-small">.{{ pricePToString(rmbAmount) }}元</span>
      </el-card>
      <el-card shadow="hover" class="my-card">
        <div class="highlight-title">虚拟货币价值</div>
        <span class="num-large">{{ Math.floor(coinAmount) }}</span
        ><span class="num-small">.{{ pricePToString(coinAmount) }}元</span>
      </el-card>
    </div>
    <h2 style="margin-top: 4rem">快速充值</h2>
    <el-form
      ref="chargeFormRef"
      :model="chargeForm"
      :rules="chargeRules"
      label-width="0px"
    >
      <div
        style="display: flex;gap: 10px;flex-wrap: nowrap;align-content: center;justify-content: space-around;align-items: flex-start;flex-direction: row;"
      >
        <el-form-item prop="selectedCoin">
          <el-select
            v-model="chargeForm.selectedCoin"
            filterable
            placeholder="请选择币种"
          >
            <el-option :key="-1" :label="'CYN'" :value="-1"> </el-option>
            <el-option
              v-for="item in coins"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item prop="chargeAmount" style="flex: 1">
          <el-input
            placeholder="请输入充值金额"
            v-model="chargeForm.chargeAmount"
          ></el-input>
        </el-form-item>
        <el-button @click="chargeButton()">充值</el-button>
      </div>
    </el-form>

    <h2 style="margin-top: 4rem">详情</h2>
    <el-card shadow="hover" class="my-card" style="width:100%">
      <el-table
        :data="tableData"
        style="text-algin:center; font-size: 15px"
        :header-cell-style="{ color: '#000000' }"
        :cell-style="{ color: '#000000' }"
        :row-style="{ height: 60 + 'px' }"
      >
        <el-table-column prop="CoinName" label="币种" align="center" />
        <el-table-column prop="Cnum" label="余额" align="center" />
        <el-table-column label="操作" align="center">
          <template slot-scope="scope">
            <el-button size="small" @click="ChargeClick(scope.row.CoinName)"
              >充值</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
import coins from "../coins";
export default {
  data() {
    return {
      coins,
      tableData: [],
      rmbAmount: 0,
      coinAmount: 0,
      chargeForm: {
        chargeAmount: null,
        selectedCoin: null
      },
      chargeRules: {
        //用户名验证
        chargeAmount: [
          { required: true,pattern: /[1-9][0-9*](\.[0-9]*)?/, message: "请输入金额",  trigger: "blur" },
        ],
        selectedCoin:[
          { required: true, message: "请选择币种", trigger: "blur" },
        ],
      }
    };
  },
  created() {
    this.getData();
  },
  mounted() {},
  methods: {
    pricePToString(x) {
      let r = Math.round((x - Math.floor(x)) * 100).toString();
      if (r.length < 2) {
        r = "0" + r;
      }
      return r;
    },
    getData() {
      this.index = (this.page - 1) * this.pageSize;
      let oldThis = this;
      this.$http
        .get("data-api/v3/cryptocurrency/myaccount", {
          headers: {
            token: window.sessionStorage["token"]
          }
        })
        .then(res => {
          oldThis.coinAmount = res.data[0];
          oldThis.rmbAmount = res.data[1];
        })
        .catch(function(error) {
          oldThis.$message.error("获取数据失败");
        });
      this.$http
        .get("data-api/v3/cryptocurrency/account", {
          headers: {
            token: window.sessionStorage["token"]
          }
        })
        .then(res => {
          oldThis.tableData = res.data;
        })
        .catch(function(error) {
          oldThis.$message.error("获取数据失败");
        });
    },
    charge(type) {
      let out = this;
      this.$prompt("请输入充值数量", "", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        inputPattern: /[1-9][0-9*](\.[0-9]*)?/,
        inputErrorMessage: "金额错误"
      }).then(({ value }) => {
        out.$http
          .post(
            "data-api/v3/cryptocurrency/getmoney",
            { num: parseFloat(value), cid: type },
            {
              headers: {
                token: window.sessionStorage["token"]
              }
            }
          )
          .then(() => {
            out.$message.success("充值成功");
            out.getData();
          })
          .catch(function(error) {
            out.$message.error("充值失败");
          });
      });
    },
    ChargeClick(name) {
      let id = coins.find(x => x["name"] == name)["id"];
      if (id == -1) {
        this.$message.error("币种不存在");
      } else {
        this.charge(id);
      }
    },
    async chargeButton() {
      let b = await this.$refs.chargeFormRef.validate();
      if (!b) {
        return;
      }
      this.$http
          .post(
            "data-api/v3/cryptocurrency/getmoney",
            { num: parseFloat(this.chargeForm.chargeAmount), cid: this.chargeForm.selectedCoin },
            {
              headers: {
                token: window.sessionStorage["token"]
              }
            }
          )
          .then(() => {
            this.$message.success("充值成功");
            this.chargeForm.selectedCoin=null;
            this.chargeForm.chargeAmount=null;
            this.getData();
          })
          .catch(function(error) {
            this.$message.error("充值失败");
          });
    }
  }
};
</script>

<style scoped>
.el-table {
  margin-top: 5px;
}
.user {
  margin-left: 83%;
}
.num-large {
  font-size: 3em;
}

.num-small {
  font-size: 2em;
}

.highlight-title {
  font-size: 1.2em;
  font-weight: bold;
  margin-bottom: 1rem;
  /* white-space: nowrap; */
}

.asset-cards {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  align-content: stretch;
  align-items: stretch;
  justify-content: space-around;
  gap: 10px;
}

.my-card {
  display: inline-block;
  margin: 20px 2%;
  width: 45%;
}
@media only screen and (max-width: 700px) {
  .my-card {
    display: block;
    margin: 20px auto;
    width: 90%;
  }
}
</style>
