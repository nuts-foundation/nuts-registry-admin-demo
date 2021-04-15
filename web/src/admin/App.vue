<template>
  <el-container>
    <el-aside>
      <h1>
        <i class="nuts-logo">Nuts</i><br/>
        registry admin
      </h1>
      <el-menu default-active="1">
        <router-link to="customers">
          <el-menu-item index="1">
            <i class="el-icon-office-building"></i>
            <span>Your care organisations</span>
          </el-menu-item>
        </router-link>
        <router-link to="vendor">
          <el-menu-item index="2">
            <i class="el-icon-school"></i>
            <span>You as a vendor</span>
          </el-menu-item>
        </router-link>
        <router-link to="vendors">
          <el-menu-item index="3">
            <i class="el-icon-lock"></i>
            <span>Trusted vendors</span>
          </el-menu-item>
        </router-link>
      </el-menu>
      <p><router-link to="/">Log out</router-link></p>
    </el-aside>
    <el-main>
      <router-view></router-view>
    </el-main>
  </el-container>
</template>

<script>
module.exports = {
  data: function () {
    return {
      greeting: 'Hello'
    }
  }
}
</script>

<style>
.el-container {
  min-height: 100vh;
}

.el-aside {
  text-align: center;
  border-right: solid 1px #e6e6e6;
}
.el-aside a {
  text-decoration: none;
  color: #2c3e50;
  font-size: 14px;
}
.el-aside a:hover {
  color: navy;
}

.el-menu {
  border-right: 0;
  text-align: left;
  margin: 4em 0;
}
.el-menu a {
  text-decoration: none;
}

.el-container .el-aside h1 {
  font-size: 1.5em;
  margin: 3em 0;
  line-height: 2em;
  font-weight: normal;
}

.narrow-page {
  margin: 2em auto;
  max-width: 800px;
}
</style>
