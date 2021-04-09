<template>
  <el-table :data="customers">
    <el-table-column prop="name" label="Name">
    </el-table-column>
    <el-table-column prop="did" label="DID">
    </el-table-column>
  </el-table>
</template>

<script>
export default {
  data() {
    return { customers: [
      {name: "Verpleeghuis de nootjes", did: "did:nuts:123"},
      {name: "Verpleeghuis notenboom", did: "did:nuts:456"},
    ]}
  }
}
</script>