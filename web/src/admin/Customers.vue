<template>
  <section class="customers">
    <header>
      <el-button type="primary" @click="onSubmit">
        <i class="el-icon-plus"></i>
        Add
      </el-button>
    </header>
    <el-table :data="customers">
      <el-table-column type="expand">
        <template #default="props">
          <el-form label-width="100px" size="mini">
            <el-form-item label="Identifier:">
              {{ props.row.identifier }}
            </el-form-item>
            <el-form-item label="DID:">
              {{ props.row.did }}
            </el-form-item>
            <el-form-item label="Status:">
              <el-tag type="success" v-if="props.row.active"><i class="el-icon-circle-check"></i></el-tag>
              <el-tag type="danger" v-if="!props.row.active"><i class="el-icon-circle-close"></i></el-tag>
              {{ props.row.active ? 'Enabled' : 'Not enabled' }}
            </el-form-item>
            <el-form-item label="Services:">
              <ul>
                <li>Service 1</li>
                <li>Service 2</li>
              </ul>
            </el-form-item>
            <el-form-item>
              <el-button>Edit</el-button>
            </el-form-item>
          </el-form>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="Name">
      </el-table-column>
      <el-table-column prop="town" label="Town">
      </el-table-column>
      <el-table-column width="100" fixed="right">
        <template #default="scope">
          <el-tag type="success" v-if="scope.row.active"><i class="el-icon-circle-check"></i></el-tag>
          <el-tag type="danger" v-if="!scope.row.active"><i class="el-icon-circle-close"></i></el-tag>
        </template>
      </el-table-column>
    </el-table>
  </section>
</template>

<script>
import { onMounted, reactive } from "vue";

export default {
  data() {
    return {
      customers: []
    }
  },
  created() {
    // watch the params of the route to fetch the data again
    this.$watch(
        () => this.$route.params,
        () => {
          this.fetchData()
        },
        // fetch the data when the view is created and the data is
        // already being observed
        { immediate: true }
    )
  },
  methods: {
    fetchData() {
      fetch("api/customers")
          .then(response => response.json())
          .then(data => this.customers = data)
    }
  },
}
</script>

<style>
.customers {
  margin: 2em auto;
  max-width: 800px;
}

.customers header {
  text-align: right;
}

.customers ul {
  margin: 0;
  padding: 0 25px;
}
</style>
