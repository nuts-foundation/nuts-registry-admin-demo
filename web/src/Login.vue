<template>
  <el-row :align="middle" :gutter="20">
    <el-col :offset="6" :span="10">
      <el-form ref="form" :model="form" label-width="120px">
        <el-form-item label="Enter username">
          <el-input placeholder="Username" v-model="form.username"></el-input>
        </el-form-item>
        <el-form-item label="Enter password">
          <el-input placeholder="Password" v-model="form.password" show-password></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">Login</el-button>
          <el-button>Cancel</el-button>
        </el-form-item>
      </el-form>
    </el-col>
  </el-row>
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
      this.$router.push("/admin/")
    }
  }
}
</script>