<template>
  <section class="narrow-page vendor">
    <h1>{{ vendor.name }}</h1>
    <el-form label-width="200px" size="mini">
      <el-form-item label="DID:">
        <el-button v-if="!registered" type="success" @click="registered = true">
          <i class="el-icon-document-add"></i>
          Register me
        </el-button>
        <span v-if="registered">did:nuts:98df7g89dfg7hoq3irj</span>
      </el-form-item>
      <el-form-item label="Address:">
        {{ vendor.address }}
      </el-form-item>
      <el-form-item label="Contact phone number:">
        {{ vendor.phone }}
      </el-form-item>
      <el-form-item label="Contact e-mail address:">
        {{ vendor.email }}
      </el-form-item>
      <el-form-item size="normal">
        <el-button type="primary" @click="editFormVisible = true">
          <i class="el-icon-edit"></i>
          Edit
        </el-button>
      </el-form-item>
    </el-form>

    <el-dialog title="Edit your vendor information" v-model="editFormVisible">
      <el-form ref="form" :model="vendor" label-width="200px">
        <el-form-item label="Vendor name:">
          <el-input v-model="vendor.name"></el-input>
        </el-form-item>
        <el-form-item label="Address:">
          <el-input v-model="vendor.address"></el-input>
        </el-form-item>
        <el-form-item label="Contact phone number:">
          <el-input v-model="vendor.phone"></el-input>
        </el-form-item>
        <el-form-item label="Contact e-mail address:">
          <el-input v-model="vendor.email"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editFormVisible = false">Cancel</el-button>
          <el-button type="primary" @click="editFormVisible = false">Save</el-button>
        </span>
      </template>
    </el-dialog>
  </section>
</template>

<script>
export default {
  data() {
    return {
      vendor: {
        name: "Nuts Zorgservices BV",
        phone: "+31 435 346 33",
        email: "nuts-admin@zorgservices.nl",
        address: "Nootjeslaan 15, 8252OD NieuwCashewdam"
      },
      editFormVisible: false,
      registered: false
    }
  }
}
</script>

<style>
.vendor {
  text-align: left;
}

.vendor .edit-button {
  text-align: right;
}
</style>
