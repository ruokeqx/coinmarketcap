<template>
    <el-container>
        <el-header>
          <el-row style="width:90%;text-algin:center; margin-left: 5%;" >
            <el-col :span="8"><div class="grid-content bg-purple"><el-image :src="src" style="width: 298px; height: 50px" class="img"></el-image></div></el-col>
          </el-row>
        </el-header>
        <el-main>
            <router-view></router-view>
        </el-main>
    </el-container>
</template>

<script>
import COIN from '../assets/COIN.png'
export default {
  data () {
    return {
      src: COIN
    }
  },
  methods: {
    handleSelect (key, keyPath) {
      console.log(key, keyPath)
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
.row{
  margin-top: 14px;
}
el-main{
  overflow: hidden;
}
</style>
