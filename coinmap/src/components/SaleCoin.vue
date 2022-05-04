<template>
  <div>
    <h1>出售货币</h1>
    <el-form ref="formRef" :model="form" :rules="formRules" label-width="0px">
      <el-row style="margin-top: 3rem">
        <el-col
          :xs="24"
          :sm="{ span: 20, offset: 2 }"
          :md="{ span: 18, offset: 3 }"
          :lg="{ span: 16, offset: 4 }"
          :xl="{ span: 14, offset: 5 }"
        >
          <el-card shadow="hover" style="text-align: center;">
            <el-row :gutter="20">
              <el-col :span="6" style="line-height: 48px;text-align: right;"
                >币种</el-col
              >
              <el-col :span="18">
                <el-form-item prop="selectedCoin">
                  <el-select
                    :disabled="fixed"
                    v-model="form.selectedCoin"
                    filterable
                    placeholder="请选择币种"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="item in coins"
                      :key="item.id"
                      :label="item.name"
                      :value="item.id"
                    >
                    </el-option>
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="6" style="line-height: 48px;text-align: right;"
                >出售数量</el-col
              >
              <el-col :span="18">
                <el-form-item prop="amount">
                  <el-input
                    placeholder="请输入出售数量"
                    v-model="form.amount"
                    style="width:100%"
                  ></el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="6" style="line-height: 48px;text-align: right;"
                >预期交易时间</el-col
              >
              <el-col :span="18">
                <el-form-item prop="expectedTime">
                  <el-date-picker
                    v-model="form.expectedTime"
                    type="datetime"
                    placeholder="请选择"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            <el-button @click="saleButton()">出售</el-button>
          </el-card>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>

<script>
import coins from "../coins";
export default {
  data() {
    let fixedId = null;
    if (this.$route.params.coin != undefined) {
      let r = coins.find(x => x.name == this.$route.params.coin);
      if (r != undefined) {
        fixedId = r.id;
      }
    }
    return {
      coins,
      fixed: fixedId != null,
      form: {
        amount: null,
        selectedCoin: fixedId,
        expectedTime: new Date().getTime() + 3600 * 1000 * 24
      },
      formRules: {
        amount: [
          {
            required: true,
            pattern: /[1-9][0-9]*(\.[0-9]*)?/,
            message: "请输入数量",
            trigger: "blur"
          }
        ],
        selectedCoin: [
          { required: true, message: "请选择币种", trigger: "blur" }
        ],
        expectedTime: [
          { required: true, message: "请选择时间", trigger: "blur" }
        ]
      }
    };
  },
  methods: {
    async saleButton() {
      let b = await this.$refs.formRef.validate();
      if (!b) {
        return;
      }
      let oldThis = this;

      this.$confirm("确认出售吗", "", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          oldThis.$http
            .post(
              "data-api/v3/cryptocurrency/transaction/create",
              {
                ExpectedTime: this.form.expectedTime / 1000,
                TsCid: this.form.selectedCoin,
                TsNum: parseFloat(this.form.amount)
              },
              {
                headers: {
                  token: window.sessionStorage["token"]
                }
              }
            )
            .then(res => {
              if (res.data.code == 200) {
                oldThis.$message.success("上架成功");
                oldThis.$router.push({ path: "/trade" });
              } else {
                oldThis.$alert(res.data.msg, "提示", {
                  confirmButtonText: "OK"
                });
              }
            })
            .catch(err => {
              oldThis.$alert("错误: " + err.response.data.msg, "提示", {
                confirmButtonText: "OK"
              });
            });
        })
        .catch(() => {
          oldThis.$message({
            type: "info",
            message: "取消出售"
          });
        });
    }
  }
};
</script>

<style scoped></style>
