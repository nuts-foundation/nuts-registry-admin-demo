<template>
  <section class="narrow-page customers">
    <header>
      <el-button type="primary" @click="createOrganisation()">
        <i class="el-icon-plus"></i>
        Add
      </el-button>
    </header>
    <el-table :data="customers">
      <el-table-column type="expand">
        <template #default="props">
          <el-form label-width="100px" size="mini">
            <el-form-item>
              <el-button type="primary" size="normal" @click="editOrganisation(props.row)">
                <i class="el-icon-edit"></i>
                Edit
              </el-button>
              <el-button type="danger" size="normal">
                <i class="el-icon-delete"></i>
                Delete
              </el-button>
            </el-form-item>
            <el-form-item label="Identifier:">
              {{ props.row.identifier }}
            </el-form-item>
            <el-form-item label="DID:">
              {{ props.row.did }}
            </el-form-item>
            <el-form-item label="Status:">
              <el-tag type="success" v-if="props.row.active"><i class="el-icon-circle-check"></i> Enabled</el-tag>
              <el-tag type="danger" v-if="!props.row.active"><i class="el-icon-circle-close"></i> Not enabled</el-tag>
            </el-form-item>
            <el-form-item label="Services:">
              <ul>
                <li>Service 1</li>
                <li>Service 2</li>
              </ul>
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

  <el-dialog title="Add a care organisation" v-model="editFormVisible">
    <el-form label-width="120px">
      <el-form-item label="Name:">
        <el-input v-model="customer.name"></el-input>
      </el-form-item>
      <el-form-item label="Town:">
        <el-input v-model="customer.town"></el-input>
      </el-form-item>
      <el-form-item label="Identifier:">
        <el-input v-model="customer.identifier"></el-input>
      </el-form-item>
      <el-form-item label="Services:">
        <el-checkbox-group v-model="customer.services">
          <el-checkbox label="Service 1" name="services"></el-checkbox>
          <el-checkbox label="Service 2" name="services"></el-checkbox>
          <el-checkbox label="Service 3" name="services"></el-checkbox>
          <el-checkbox label="Service 4" name="services"></el-checkbox>
        </el-checkbox-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="editFormVisible = false">Cancel</el-button>
        <el-button type="primary" @click="editFormVisible = false">Save</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import { onMounted, reactive } from "vue";

export default {
  data() {
    return {
      customers: [],
      customer: {},
      editFormVisible: false
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
    },
    createOrganisation() {
      this.customer = {
        name: '',
        town: '',
        identifier: '',
        services: []
      };
      this.editFormVisible = true;
    },
    editOrganisation(org) {
      this.customer.name = org.name;
      this.customer.town = org.town;
      this.customer.identifier = org.identifier;
      this.customer.services = org.services || [];
      this.editFormVisible = true;
    }
  },
}
</script>

<style>
.customers header {
  text-align: right;
}

.customers ul {
  margin: 0;
  padding: 0 25px;
}
</style>
