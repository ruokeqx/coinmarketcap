<template>
<el-container>
    <el-header style="height: auto">
      <el-row justify="center" type="flex">
        <el-col :xs="24" :sm="22" :md="20" :lg="18" :xl="16">
          <div id="head-line">
            <el-image :src="src" style="width: 216px; height: 37px;float: left" class="img" @click="homeClick()"></el-image>
            <el-popover
              class="user"
              placement="bottom"
              width="80"
              trigger="hover"
              style="float: right">
              <div>
                <el-button class="pop-item" @click="homeClick()" style="border: none">主页</el-button>
                <el-button class="pop-item" @click="scClick()" style="border: none">我的收藏</el-button>
                <el-button class="pop-item" @click="gmClick()" style="border: none">购买货币</el-button>
                <el-button class="pop-item" @click="assetClick()" style="border: none">个人资产</el-button>
                <el-button class="pop-item" @click="logoutClick()" style="border: none">退出登陆</el-button>
              </div>
              <el-button slot="reference" icon="el-icon-user" size="medium" circle></el-button>
            </el-popover>
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
import COIN from '../assets/COIN.png'
export default {
  data () {
    return {
      src: COIN,
    }
  },
  methods: {
    scClick(){
      this.$router.push({
        path: '/collection'
      })
    },
    assetClick(){
      this.$router.push({
        path: '/asset'
      })
    },
    gmClick(){
      console.log("gm")
    },
    homeClick(){
      this.$router.push({
        path: '/cryptocurrency'
      })
    },
    logoutClick(){
      window.sessionStorage.removeItem("token");
      this.$router.push({
        path: '/login'
      })
    }
  }
}
</script>

<style scoped>
*{
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: "Poppins" , sans-serif;
}
body{
  min-height: 100vh;
  width: 100%;
  background: #EEECEB;
}
el-main{
  overflow: hidden;
}
#head-line{
  width:90%;
  text-align:center;
  margin-left: 5%;
  margin-top: 10px;
}

.pop-item{
  display: block;
  padding: 8px;
  margin: 3px auto;
}
</style>
