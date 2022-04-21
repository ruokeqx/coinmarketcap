<template>
    <div class="login_container">
        <div class="login_box">
            <!-- 头像 -->
            <div class="avater_box">
                <img src="../assets/logo.png">
            </div>
            <!-- 表单 -->
            <el-form ref="loginFormRef" :model="loginForm" :rules="loginFormRules" class="login_form" label-width="0px">
               <div style="font-weight: bolder;font-size: 18px; margin-bottom: 5px; margin-left: 5px">登录</div>
                <!-- 用户名 -->
                <el-form-item prop="username">
                    <el-input v-model="loginForm.username" prefix-icon="iconfont icon-user" placeholder="用户名"></el-input>
                </el-form-item>
                <!-- 密码 -->
                <el-form-item prop="password">
                    <el-input v-model="loginForm.password" prefix-icon="iconfont icon-3702mima" type="password" placeholder="密码" show-password></el-input>
                </el-form-item>
                <!-- 按钮 -->
                <el-form-item class="btns">
                    <el-button type="primary" @click="login">登录</el-button>
                    <el-button type="info" @click="register">注册</el-button>
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>

<script>
export default {
    data(){
        return{
            //登录表单对象
            loginForm:{
                username:'',
                password:''
            },
            //表单验证规则
            loginFormRules:{
                //用户名验证
                username:[
                    { required: true, message:"请输入用户名", trigger:"blur"},
                    { min:3, max:10, message:"长度在3到10个字符", trigger:"blur"}
                ],
                //密码验证
                password:[
                     { required: true, message:"请输入密码", trigger:"blur"},
                    { min:6, max:15, message:"长度在6到15个字符", trigger:"blur"}
                ]
            }
        }
    },
    methods:{
        //重置表单
        register(){
            // console.log(this);
           this.$router.push('/register');
        },
       async login(){
             this.$refs.loginFormRef.validate(async vaild =>{
                if(!vaild) return;
                if(this.loginForm.password == '123456'){
                    this.$message.error('登录失败:用户名或密码错误！');
                    return;
                }
                const {data: res}=await this.$http.post('login',this.loginForm);
                console.log(res)
                if(res.code == 200)
                    this.$message.success('登录成功');
                else
                    return this.$message.error('登录失败');
                // 将登录成功后的token保存到seessionStorage中
                window.sessionStorage.setItem("token",res.data);
                // 通过编程式导航跳转到后台主页，路由地址是/home
                this.$router.push({
                    path:'/cryptocurrency',
                    query: {
                        username: this.loginForm.username
                    }
                });
            });
        }
    }
}
</script>

<style scoped>
    .login_container{
        background-color: #2b4b6b;
        height: 100%;
    }
    .login_box{
        width: 450px;
        height: 300px;
        background-color: #fff;
        border-radius: 3px;
        position: absolute;
        left: 50%;
        top: 50%;
        transform: translate(-50%,-50%);
    }
    .btns{
        display: flex;
        justify-content: flex-end;
    }
    .login_form{
        position: absolute;
        bottom: 0;
        width: 100%;
        padding: 20px;
        box-sizing: border-box;
    }
    .avater_box{
            height: 130px;
            width: 130px;
            border: 1px solid #eee;
            border-radius: 50%;
            padding: 10px;
            box-shadow: #ddd 0 0 10px;
            position: absolute;
            left: 50%;
            transform: translate(-50%,-50%);
            background-color: #fff;
            
    }
    img{
        height: 100%;
        width: 100%;
        border-radius: 50%;
        background-color: #eee;
    }
</style>
