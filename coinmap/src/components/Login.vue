<template>
  <el-row type="flex" justify="center">
    <el-col
      :xs="24"
      :sm="15"
      :md="10"
      :lg="8"
      :xl="6"
      style="margin-top: 3rem; margin-bottom: 2rem"
    >
      <!-- 头像 -->
      <div class="avater_box">
        <img src="../assets/logo.png" />
      </div>
      <!-- 表单 -->
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginFormRules"
        class="login_form"
        label-width="0px"
      >
        <div
          style="
            font-weight: bolder;
            font-size: 18px;
            margin-bottom: 5px;
            margin-left: 5px;
          "
        >
          {{ mode ? "登录" : "注册" }}
        </div>
        <!-- 用户名 -->
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            prefix-icon="iconfont icon-user"
            placeholder="用户名"
          ></el-input>
        </el-form-item>
        <!-- 密码 -->
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            prefix-icon="iconfont icon-3702mima"
            type="password"
            placeholder="密码"
            show-password
          ></el-input>
        </el-form-item>
        <!-- 按钮 -->
        <el-form-item class="btns">
          <el-button type="primary" @click="mode?login():register()">{{
            mode ? "登录" : "注册"
          }}</el-button>
          <el-button type="info" @click="mode = !mode">{{
            mode ? "注册" : "登录"
          }}</el-button>
        </el-form-item>
      </el-form>
    </el-col>
  </el-row>
</template>

<script>
export default {
  data() {
    return {
      //登录表单对象
      loginForm: {
        username: "",
        password: "",
      },
      //表单验证规则
      loginFormRules: {
        //用户名验证
        username: [
          { required: true, message: "请输入用户名", trigger: "blur" },
          { min: 3, max: 10, message: "长度在3到10个字符", trigger: "blur" },
        ],
        //密码验证
        password: [
          { required: true, message: "请输入密码", trigger: "blur" },
          { min: 6, max: 15, message: "长度在6到15个字符", trigger: "blur" },
        ],
      },
      mode: true,
    };
  },
  methods: {
    async register() {
      let vaild=await this.$refs.loginFormRef.validate();
      if (!vaild) return;
      const { data: res } = await this.$http.post("register", this.loginForm);
      if (res.code !== 200) return this.$message.error("注册失败");
      this.$message.success("注册成功");
      await this.login();

    },
    async login() {
      let vaild=await this.$refs.loginFormRef.validate();
      if (!vaild) return;
      const { data: res } = await this.$http.post("login", this.loginForm);
      if (res.code == 200) this.$message.success("登录成功");
      else return this.$message.error("登录失败");
      // 将登录成功后的token保存到seessionStorage中
      window.sessionStorage.setItem("token", res.data);
      // 通过编程式导航跳转到后台主页，路由地址是/home
      this.$router.push({
        path: "/cryptocurrency",
      });
    },
  },
};
</script>

<style scoped>
.btns {
  display: flex;
  justify-content: flex-end;
}
.login_form {
  width: 100%;
  padding: 20px;
  box-sizing: border-box;
}
.avater_box {
  height: 130px;
  width: 130px;
  border: 1px solid #eee;
  border-radius: 50%;
  padding: 10px;
  box-shadow: #ddd 0 0 10px;
  background-color: #fff;
  margin: auto;
}
img {
  height: 100%;
  width: 100%;
  border-radius: 50%;
  background-color: #eee;
}
</style>
