<template>
  <el-container>
    <el-header style="height: auto">
      <el-row justify="center" type="flex">
        <el-col :xs="24" :sm="22" :md="20" :lg="18" :xl="16">
          <div id="head-line">
            <el-image
              :src="src"
              style="width: 140px; height: 25px;float: left;margin-top:5px;"
              class="img"
              @click="homeClick()"
            ></el-image>
            <div style="float: right">

              <el-badge
                :value="messages.length"
                :hidden="messages.length == 0"
                class="item"
                style="margin-right: 10px;"
              >
                <el-popover placement="bottom" width="80" trigger="hover">
              
                  <div v-if="messages.length == 0" style="text-align: center">
                    无消息
                  </div>
                  <div v-for="item in messages" :key="item.MsgId">
                    <el-link
                      type="info"
                      :underline="false"
                      @click="clickMessage(item)"
                    >
                      <div v-if="item.MsgType == 0">关注的货币降价了！</div>
                      <div v-if="item.MsgType == 1">交易超过预期时间</div>
                      <div v-if="item.MsgType == 2">出售的货币交易成功！</div>
                    </el-link>
                    <hr />
                  </div>
                  <div>
                    <el-button
                      class="pop-item"
                      v-if="messages.length != 0"
                      @click="clearMessage()"
                      >清空消息</el-button
                    >
                  </div>
                  <el-button
                    slot="reference"
                    :icon="
                      messages.length == 0
                        ? 'el-icon-bell'
                        : 'el-icon-message-solid'
                    "
                    size="medium"
                    circle
                  ></el-button>
                </el-popover>
              </el-badge>
              <el-popover placement="bottom" width="80" trigger="hover" v-model="visible">
                <div>
                  <el-button
                    class="pop-item"
                    @click="$router.push({ path: '/trade' });visible=false;"
                    >交易中心</el-button
                    >
                  <el-button
                    class="pop-item"
                    @click="$router.push({ path: '/collection' });visible=false;"
                    >我的收藏</el-button
                  >
                  <el-button
                    class="pop-item"
                    @click="$router.push({ path: '/my_trading' });visible=false;"
                    >个人交易</el-button
                  >
                  <el-button
                    class="pop-item"
                    @click="$router.push({ path: '/asset' });visible=false;"
                    >个人资产</el-button
                  >
                  <el-button
                    class="pop-item"
                    @click="$router.push({ path: '/sale_coin' });visible=false;"
                    >出售货币</el-button
                  >
                  <el-button class="pop-item" @click="logoutClick()"
                    >退出登陆</el-button
                  >
                </div>
                <el-button
                  slot="reference"
                  icon="el-icon-user"
                  size="medium"
                  circle
                ></el-button>
              </el-popover>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-header>
    <el-row justify="center" type="flex" style="margin-top: 3rem">
      <el-col :xs="24" :sm="22" :md="20" :lg="18" :xl="16">
        <router-view></router-view>
      </el-col>
    </el-row>
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
  </el-container>
</template>

<script>
import COIN from "../assets/COIN.png";
export default {
  data() {
    return {
      visible:false,
      src: COIN,
      messages: [],
      timer: null
    };
  },
  created() {
    this.getMessage();
    this.timer = setInterval(this.getMessage, 60*1000);
  },
  beforeDestroy() {
    clearInterval(this.timer);
  },

  methods: {
    homeClick() {
      this.$router.push({
        path: "/cryptocurrency"
      });
      this.visible=false;
    },
    logoutClick() {
      window.sessionStorage.removeItem("token");
      this.$router.push({
        path: "/login"
      });
      this.visible=false;
    },

    async getMessage() {
      try {
        let res = await this.$http.get(
          "/data-api/v3/cryptocurrency/transaction/msg",
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
        this.messages = res.data;
        return true;
      } catch (err) {
        this.$message.error("错误: " + err.response.data.msg);
        return false;
      }
    },
    async readMessage(msgIds) {
      try {
        let res = await this.$http.post(
          "/data-api/v3/cryptocurrency/transaction/readmsg",
          msgIds,
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
        return true;
      } catch (err) {
        this.$message.error("错误: " + err.response.data.msg);
        return false;
      }
    },
    async clickMessage(msg) {
      let TsId = msg.TsId;
      await this.readMessage([msg.MsgId]);
      await this.getMessage();
      this.$router.push({ name: "transaction", params: { TsId } });
    },
    async clearMessage() {
      let result = await this.readMessage(this.messages.map(x => x.MsgId));
      if (result) {
        this.$message.success("清除成功");
      } else {
        this.$message.error("清除失败");
      }
      await this.getMessage();
    }
  }
};
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: "Poppins", sans-serif;
}
body {
  min-height: 100vh;
  width: 100%;
  background: #eeeceb;
}
el-main {
  overflow: hidden;
}
#head-line {
  width: 90%;
  text-align: center;
  margin-left: 5%;
  margin-top: 10px;
}

.pop-item {
  display: block;
  padding: 8px;
  margin: 3px auto;
  border: none;
}
hr {
  display: block;
  height: 1px;
  border: 0;
  border-top: 1px solid #ccc;
  margin: 7px 0;
  padding: 0;
}
</style>
