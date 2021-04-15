<template>
  <section class="narrow-page vendors">
    <el-table :data="vendors">
      <el-table-column type="expand">
        <template #default="props">
          <el-form label-width="200px" size="mini">
            <el-form-item label="Address:">
              {{ props.row.address }}
            </el-form-item>
            <el-form-item label="Contact phone number:">
              <a :href="`tel:${ props.row.phone }`">{{ props.row.phone }}</a>
            </el-form-item>
            <el-form-item label="Contact e-mail address:">
              <a :href="`mailto:${ props.row.email }`">{{ props.row.email }}</a>
            </el-form-item>
            <el-form-item label="Types:">
              <ul>
                <li>Type 1</li>
                <li>Type 2</li>
              </ul>
            </el-form-item>
          </el-form>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="Name">
      </el-table-column>
      <el-table-column prop="did" label="DID">
      </el-table-column>
      <el-table-column fixed="right" label="Trusted" width="100">
        <template #default="scope">
          <el-switch
            v-model="scope.row.active"
            active-color="#13ce66"
            inactive-color="#ff4949"
            @click="confirmDialogVisible = scope.row.active">
          </el-switch>
        </template>
      </el-table-column>
    </el-table>
  </section>

  <el-dialog
    title="Are you sure?"
    v-model="confirmDialogVisible"
    width="30%"
    center>
    <span>Are you really sure you want to trust this vendor?</span>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="confirmDialogVisible = false">No, please cancel</el-button>
        <el-button type="primary" @click="confirmDialogVisible = false">Yes, I'm sure</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
export default {
  data() {
    return {
      vendors: [
        {
          name: "Enovation Point",
          phone: "+31 435 346 33",
          email: "nuts-admin@point.nl",
          address: "Rotterdamseweg 523, 1538EM Rotterdam",
          did: 'did:nuts:346j457kljh7k56jh756k7hjlskjfodig',
          types: "type1, type2",
          active: true
        },
        {
          name: "Nedap Healthcare",
          phone: "+31 435 346 33",
          email: "nuts-admin@nedap.nl",
          address: "Groenloseweg 25, 1538EM Groenlo",
          did: 'did:nuts:dfg8hfg87hyw3kjrhwg9867yqk34jhsgd',
          types: "type1",
          active: false
        }
      ],
      confirmDialogVisible: false
    }
  }
}
</script>

<style>
.vendors ul {
  margin: 0;
  padding: 0 25px;
}
</style>
