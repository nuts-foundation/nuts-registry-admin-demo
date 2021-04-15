<template>
  <section class="login-form">
    <h1>
      <i class="nuts-logo">Nuts</i><br/>
      registry admin
    </h1>
    <el-row>
      <el-col :offset="8" :span="8">
        <el-form ref="form" :model="form" label-position="top">
          <el-form-item label="What's the password?">
            <el-input placeholder="Password" v-model="form.password" show-password></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onSubmit">Go!</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </section>
</template>

<script>
export default {
  data() {
    return {
      form: {
        username: '',
        password: ''

      }
    }
  },
  methods: {
    onSubmit() {
      this.$router.push("/admin/customers")
    }
  }
}
</script>

<style>
.login-form h1 {
  margin: 3em 0;
  line-height: 3em;
  font-weight: normal;
}
</style>
