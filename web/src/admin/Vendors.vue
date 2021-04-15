<template>
  <section class="vendors">
    <el-table :data="vendors">
      <el-table-column prop="did" label="DID">
      </el-table-column>
      <el-table-column prop="types" label="Types">
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
          did: 'did:nuts:346j457kljh7k56jh756k7hjlskjfodig',
          types: "type1, type2",
          active: true
        },
        {
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
.vendors {
  margin: 2em auto;
  max-width: 800px;
}
</style>
