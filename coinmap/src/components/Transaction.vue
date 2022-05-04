<template>
  <div>
    <el-row style="margin-top: 3rem">
      <el-col
        :xs="24"
        :sm="{ span: 22, offset: 1 }"
        :md="{ span: 20, offset: 2 }"
        :lg="{ span: 18, offset: 3 }"
        :xl="{ span: 16, offset: 4 }"
      >
        <el-card shadow="hover" style="text-align: center;">
          <el-empty
            v-if="transaction == null"
            description="没有内容"
          ></el-empty>
          <el-descriptions
            v-else
            class="margin-top"
            title="详情"
            :column="2"
            border
          >
            <template slot="extra">
              <el-button
                size="small"
                v-if="
                  (transaction.TsStatus == 0 || transaction.TsStatus == 1) &&
                    transaction.SellerId != uid
                "
                @click="buyClick()"
                >购买</el-button
              >
              <el-button
                size="small"
                v-if="
                  (transaction.TsStatus == 0 || transaction.TsStatus == 1) &&
                    transaction.SellerId == uid
                "
                @click="takeDown()"
                >下架</el-button
              >
              <el-button
                size="small"
                v-if="
                  (transaction.TsStatus == 0 || transaction.TsStatus == 1) &&
                    transaction.SellerId == uid
                "
                @click="discount()"
                >打折</el-button
              >
            </template>
            <el-descriptions-item>
              <template slot="label">
                币种
              </template>
              {{ coins.find(x => x.id == transaction.TsCid).name }}
            </el-descriptions-item>
            <el-descriptions-item>
              <template slot="label">
                数量
              </template>
              {{ transaction.TsNum }}
            </el-descriptions-item>
            <el-descriptions-item>
              <template slot="label">
                折扣
              </template>
              {{ transaction.Discount }}
            </el-descriptions-item>
            <el-descriptions-item>
              <template slot="label">
                当前市场价
              </template>
              {{ transaction.Cost.toFixed(2) }}
              <!-- <el-tag size="small">School</el-tag> -->
            </el-descriptions-item>
            <el-descriptions-item>
              <template slot="label">
                交易创建时间
              </template>
              {{ dateToString(transaction.TsCreaTime) }}
            </el-descriptions-item>
            <el-descriptions-item>
              <template slot="label">
                预期交易时间
              </template>
              {{ dateToString(transaction.ExpectedTime) }}
            </el-descriptions-item>
            <el-descriptions-item>
              <template slot="label">
                交易状态
              </template>
              {{ statusToStr(transaction.TsStatus) }}
            </el-descriptions-item>
            <el-descriptions-item v-if="transaction.TsStatus == 3">
              <template slot="label">
                交易完成时间
              </template>
              {{ dateToString(transaction.TsCloseTime) }}
            </el-descriptions-item>
            <el-descriptions-item v-if="transaction.TsStatus == 2">
              <template slot="label">
                交易关闭时间
              </template>
              {{ dateToString(transaction.TsCloseTime) }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import coins from "../coins";
export default {
  data() {
    let uid = JSON.parse(
      atob(window.sessionStorage.getItem("token").split(".")[1])
    )["uid"];
    return {
      uid,
      coins,
      transaction: null
    };
  },
  async mounted() {
    await this.getData();
  },
  methods: {
    async getData() {
      let TsId = parseInt(this.$route.params.TsId);
      try {
        let res = await this.$http.post(
          "/data-api/v3/cryptocurrency/transaction/onetransaction",
          { TsId },
          {
            headers: {
              token: window.sessionStorage["token"]
            }
          }
        );
        if (res.status != 200) {
          this.$message.error("错误: " + res.data.msg);
          return false;
        }
        this.transaction = res.data;
        return true;
      } catch (err) {
        this.$message.error("错误: " + err.response.data.msg);
        return false;
      }
    },
    buyClick() {
      let TsId = parseInt(this.$route.params.TsId);
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
            message: "取消购买"
          });
        });
    },
    takeDown() {
      let TsId = parseInt(this.$route.params.TsId);
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
            message: "取消购买"
          });
        });
    },
    discount() {
      let TsId = parseInt(this.$route.params.TsId);
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
        .catch(() => {});
    },
    dateToString(timeStamp) {
      let date = new Date(timeStamp * 1000);
      const toTwoStr = x => (x < 10 ? "0" + x : "" + x);
      return (
        date.getFullYear() +
        "-" +
        toTwoStr(date.getMonth() + 1) +
        "-" +
        toTwoStr(date.getDate()) +
        " " +
        toTwoStr(date.getHours()) +
        ":" +
        toTwoStr(date.getMinutes()) +
        ":" +
        toTwoStr(date.getSeconds())
      );
    },
    statusToStr(status) {
      switch (status) {
        case 0:
          return "出售中";
        case 1:
          return "交易超时";
        case 2:
          return "交易关闭";
        case 3:
          return "交易完成";
      }
    }
  }
};
</script>

<style></style>
